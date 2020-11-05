package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"gopkg.in/yaml.v2"
)

const (
	// https://en.wikipedia.org/wiki/ANSI_escape_code
	red     = "\033[01;31m"
	green   = "\033[01;32m"
	yellow  = "\033[01;33m"
	blue    = "\033[01;34m"
	white   = "\033[01;97m"
	neutral = "\033[0m"
)

var areas = [52][4]int{
	{}, {0, 0, 2, 0}, {8, 0, 3, 1}, {9, 0, 0, 2}, {10, 0, 5, 0},
	{11, 0, 0, 4}, {13, 0, 7, 0}, {0, 8, 0, 0}, {7, 2, 9, 0},
	{15, 3, 10, 8}, {0, 4, 11, 9}, {0, 5, 12, 10}, {0, 0, 0, 11},
	{16, 6, 0, 0}, {17, 0, 15, 0}, {0, 9, 0, 14}, {0, 13, 17, 0},
	{0, 14, 18, 16}, {24, 0, 19, 17}, {25, 0, 20, 18}, {0, 0, 21, 19},
	{26, 0, 0, 20}, {27, 0, 23, 0}, {0, 0, 24, 22}, {29, 18, 0, 23},
	{30, 19, 0, 0}, {32, 21, 0, 0}, {0, 22, 28, 0}, {0, 0, 0, 27},
	{0, 0, 30, 0}, {0, 25, 0, 0}, {0, 0, 0, 0}, {0, 26, 0, 33},
	{0, 34, 32, 0}, {33, 36, 37, 35}, {0, 0, 34, 0}, {34, 0, 5, 0},
	{37, 38, 39, 34}, {37, 38, 38, 38}, {40, 39, 39, 37}, {42, 39, 0, 0},
	{43, 41, 42, 41}, {44, 40, 42, 41}, {47, 41, 43, 43}, {48, 42, 45, 44},
	{49, 45, 45, 44}, {50, 51, 0, 0}, {47, 43, 48, 47}, {48, 44, 48, 47},
	{0, 45, 50, 0}, {0, 46, 0, 49}, {46, 0, 0, 0},
}

var areaMap = [12][10]int{
	{0, 0, 0, 0, 0, 0, 47, 48, 49, 50},
	{0, 0, 0, 0, 0, 0, 43, 44, 45, 46},
	{0, 0, 0, 0, 0, 0, 41, 42, 0, 51},
	{0, 0, 0, 0, 33, 32, 0, 40, 0, 0},
	{0, 0, 0, 35, 34, 0, 37, 39, 0, 0},
	{27, 28, 29, 30, 36, 0, 38, 0, 0, 0},
	{22, 23, 24, 25, 0, 26, 0, 0, 0, 0},
	{16, 17, 18, 19, 20, 21, 0, 0, 0},
	{13, 14, 15, 0, 0, 0, 0, 0, 0, 0},
	{6, 7, 0, 31, 0, 0, 0, 0, 0, 0},
	{0, 8, 9, 10, 11, 12, 0, 0, 0, 0},
	{1, 2, 3, 4, 5, 0, 0, 0, 0, 0},
}

var objectsInArea = [45][2]int{
	{-1, 0}, {}, {}, {}, {}, {}, {}, {}, {},
	{28, 0}, {29, 0}, {8, 0}, {24, 0}, {2, 0}, {26, 0}, {20, 10}, {19, 0}, {19, 22},
	{18, 0}, {18, 20}, {16, 0}, {13, 0}, {14, 0}, {0, 0}, {0, 26}, {15, 5}, {8, 5},
	{9, 0}, {11, 0}, {12, 0}, {1, 0}, {1, 0}, {1, 0}, {3, 7}, {4, 0}, {4, 0},
	{4, 0}, {0, 0}, {0, 0}, {0, 18}, {30, 0}, {33, 10}, {51, 0}, {40, 0}, {51, 47},
}

var doorOpen = false

type conf struct {
	Verbs     []string `yaml:"verbs"`
	Nouns     []string `yaml:"nouns"`
	Objects   []string `yaml:"objects"`
	Answers   []string `yaml:"answers"`
	Locations []string `yaml:"locations"`
}

func (c *conf) getConf(filename string) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func appendText(block *[]string, newText string, color ...string) {
	text := *block
	if color == nil {
		*block = append(text, newText)
	}
	*block = append(text, fmt.Sprintf("%s%s%s", color[0], newText, white))
}

func printScreen(text []string) {
	// clear screen
	fmt.Print("\033[H\033[2J")
	block := strings.Join(text, "\n")
	fmt.Println(block)
}

func getBoxLen(locations []string) int {
	boxLen := 0
	for _, v := range locations {
		lineLen := len([]rune(strings.Split(v, "\n")[0]))
		if lineLen > boxLen {
			boxLen = lineLen
		}
	}
	// make boxLen odd for middle connection piece
	if boxLen%2 == 0 {
		boxLen++
	}
	return boxLen + 4
}

func drawBox(area int, boxLen int, locations []string) (box [3]string) {
	if area == 0 {
		horLine := strings.Repeat(" ", boxLen)
		box[0] = fmt.Sprintf("%s", horLine)
		box[1] = fmt.Sprintf("%s", horLine)
		box[2] = fmt.Sprintf("%s", horLine)
		return
	}
	var leftCon, rightCon, topCon, bottomCon string
	text := strings.Split(locations[area-1], "\n")[0]
	textLen := len([]rune(text)) + 2 // one space left and right
	leftBuffer := strings.Repeat(" ", (boxLen-textLen)/2)
	rightBuffer := strings.Repeat(" ", boxLen-len(leftBuffer)-textLen)
	horLine := strings.Repeat("\u2501", (boxLen-3)/2)
	if areas[area][3] == 0 {
		leftCon = "\u2503"
	} else {
		leftCon = "\u252B"
	}
	if areas[area][2] == 0 {
		rightCon = "\u2503"
	} else {
		rightCon = "\u2523"
	}
	if areas[area][0] == 0 {
		topCon = "\u2501"
	} else {
		topCon = "\u253B"
	}
	if areas[area][1] == 0 {
		bottomCon = "\u2501"
	} else {
		bottomCon = "\u2533"
	}
	box[0] = fmt.Sprintf("\u250F%s%s%s\u2513", horLine, topCon, horLine)
	box[1] = fmt.Sprintf("%s%s%s%s%s", leftCon, leftBuffer, text, rightBuffer, rightCon)
	box[2] = fmt.Sprintf("\u2517%s%s%s\u251B", horLine, bottomCon, horLine)
	return
}

// Setup keyboard scanning
func scanner() (r rune) {
	var b []byte = make([]byte, 4)
	os.Stdin.Read(b)
	r, _ = utf8.DecodeRune(b)
	//fmt.Println(b)
	//fmt.Println(r)
	return
}

func prelude() {
	var text []string

	appendText(&text, "fantasia", red)
	appendText(&text, "- Ein Adventure von Klaus Hartmuth -", blue)
	appendText(&text, "- Überarbeitet von Tom Hutter -", yellow)
	printScreen(text)
	//scanner()
}

func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
		os.Exit(0)
	}()
}

//func directions() {
/*
	        175 rem ** richtungen *****************
	176 fori=1to51:forj=0to3:readr(i,j):next:next
*/
//}

func flash(text []string, err string) {
	flashText := make([]string, len(text))
	copy(flashText, text)
	flashText = append(text, "")
	flashText = append(text, fmt.Sprintf("%s%s", red, err))
	printScreen(flashText)
	time.Sleep(2 * time.Second)
	printScreen(text)
}

func surroundings(area int, locations []string, objects []string) (text []string) {
	if area == 25 {
		objectsInArea[40][0] = 25
		objects[40-9] = "eine Tür im Norden"
	}
	if area == 30 {
		objectsInArea[40][0] = 30
		objects[40-9] = "eine Tür im Süden"
	}

	//	thenge(40)=25:ge$(40)="eine tuer im norden"
	//	ifoa=30thenge(40)=30:ge$(40)="eine tuer im sueden"
	//fmt.Printf("Ich bin %s\n", locations[area-1])
	//var text []string
	text = append(text, fmt.Sprintf("%sArea: %d [N:%d,S:%d,O:%d,W:%d]", white, area, areas[area][0], areas[area][1], areas[area][2], areas[area][3]))
	text = append(text, fmt.Sprintf("%sIch bin %s", yellow, locations[area-1]))
	//appendText(&text, fmt.Sprintf("Ich bin %s", locations[area-1]), yellow)
	var items []string
	for i, v := range objectsInArea {
		if v[0] == area {
			//items = append(items, objects[i-9])
			item := objects[i-9]
			if strings.Contains(item, "::") {
				item = strings.ReplaceAll(item, "::", "")
				items = append(items, fmt.Sprintf("%s  - %s", white, item))
				//appendText(&text, fmt.Sprintf("  - %s", item), white)
			} else {
				items = append(items, fmt.Sprintf("%s  - %s", blue, item))
				//appendText(&text, fmt.Sprintf("  - %s", item), blue)
			}
		}
	}
	if len(items) > 0 {
		text = append(text, "")
		text = append(text, fmt.Sprintf("%sIch sehe:", blue))
		for _, item := range items {
			text = append(text, item)
		}
		text = append(text, neutral)
	}
	printScreen(text)
	return
	//appendText(&text, "Ich sehe:", blue)
	/*if v[1] == area {
		//items = append(items, objects[i-9])
		appendText(&text, fmt.Sprintf("  - %s", objects[i-9]), red)
	}*/
	//printScreen(text)
	/*if len(items) > 0 {
		fmt.Println("Ich sehe:")
		for _, i := range items {
			fmt.Printf("  - %s\n", i)
		}
	}*/

	/*
		ifoa=25thenge(40)=25:ge$(40)="eine tuer im norden"
		ifoa=30thenge(40)=30:ge$(40)="eine tuer im sueden"
		ifoa=6thenp1=1
		ifoa=31thenp2=1
		ifoa=29thenp3=1
		ifoa<>1andge(31)<>-2thenprinte$:poke214,5:poke211,3:sysvd:fl=1
		iffl=1thenprint"hilfe !   ich versinke im boden."
		iffl=1thenfl=0:pokevc,peek(vc)or16:fori=1to2000:next:goto611
		printc$"ich bin "o$(oa)d$:fl=0:fori=9to44:ifge(i)<>oathen323
		iffl=0thenprintf$"ich sehe:"
		printge$(i):fl=1
		next:fl=0
		ifoa=31then335
		ifin>1andoa=5then327
		goto331
		fori=9to44:if(ge(i)=-1orge(i)=-2)andi<>31thenge(i)=29
		next:in=1
		print"im moor ist alles verschwunden,"
		print"was ich bei mir hatte !"
		fl=0:printf$"richtungen:":fori=0to3:ifr(oa,i)=0then334
		iffl=1thenprint", ";
		printno$(i+5);:fl=1
		next:fl=0
		printtc$f$:fori=1to40:printchr$(175);:next:printd$;:return
	*/
}

func move(area int, direction int, text []string) int {
	//if direction == 0 {
	//	return 0, "Ich brauche eine Richtung."
	//}
	if areas[area][direction] == 0 {
		flash(text, "In diese Richtung führt kein Weg.")
		return area
	}
	// Area 30 and 25 are connected by a door. Is it open?
	if (area == 30 || area == 25 && direction == 0) && !doorOpen {
		flash(text, "Die Tür ist versperrt.")
		return area
	}
	return areas[area][direction]
}

func useDoor() {
	/*
		495 rem ** sperre *********************
		496 f=0:gosub605:iffl=1thenfl=0:goto280
		497 ifno<>40andno<>35thenprinta$(2):goto280
		498 ifno=35thenprint"versuche 'oeffne'.":goto280
		499 iftu=1thenprint"ist schon offen !":goto280
		500 ifge(26)<>-1thenprint"ich habe keinen schluessel.":goto280
		501 print"gut.":tu=1:goto281
		502 :
	*/
}

func use(object int, area int) {
	/*
		604 rem ** unterprogramm **************
		605 ifge(no)<>oaandge(no)<>-1andge(no)<>-2thenfl=1
		606 iffl=1thenprint"sehe ich hier nicht.":return
		607 iffl=1andge(no)<>-1andge(no)<>-2thenprint"habe ich nicht dabei.":fl=1
		608 return
		if objectsInArea[object][0] != area
	*/
}

func main() {
	var c conf
	//c.getConf("verbs.yaml")
	//verbs := c.Verbs
	//c.getConf("nouns.yaml")
	//nouns := c.Nouns
	c.getConf("objects.yaml")
	objects := c.Objects
	c.getConf("locations.yaml")
	locations := c.Locations

	boxLen := getBoxLen(locations)
	for i := 0; i < 12; i++ {
		box1 := drawBox(areaMap[i][1], boxLen, locations)
		box2 := drawBox(areaMap[i][2], boxLen, locations)
		box3 := drawBox(areaMap[i][3], boxLen, locations)
		fmt.Printf("%s%s%s\n", box1[0], box2[0], box3[0])
		fmt.Printf("%s%s%s\n", box1[1], box2[1], box3[1])
		fmt.Printf("%s%s%s\n", box1[2], box2[2], box3[2])
	}
	return
	/*
		box1 := drawBox(1, boxLen, locations)
		box2 := drawBox(2, boxLen, locations)
		box3 := drawBox(3, boxLen, locations)
		fmt.Printf("%s%s%s\n", box1[0], box2[0], box3[0])
		fmt.Printf("%s%s%s\n", box1[1], box2[1], box3[1])
		fmt.Printf("%s%s%s\n", box1[2], box2[2], box3[2])
		return
	*/
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// reenable display entered characters on the screen
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	// Setup our Ctrl+C handler
	setupCloseHandler()
	prelude()
	scanner()
	area := 1
	var dir rune
	var direction int
	var text []string
	text = surroundings(area, locations, objects)
	for {
		dir = scanner()
		switch int(dir) {
		case 110: // N
			direction = 0
		case 115: // S
			direction = 1
		case 111: // O
			direction = 2
		case 119: // W
			direction = 3
		}

		area = move(area, direction, text)
		text = surroundings(area, locations, objects)
	}
	//scanner()

	//surroundings(8, locations, objects)
	//scanner()

	//fmt.Println(verbs)
	//fmt.Println(nouns)
	//fmt.Println(objects[0])
	//fmt.Println(locations[0])
}
