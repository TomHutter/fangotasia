package main

import (
	"fantasia/config"
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

var doorOpen = false

type mapOverwrites struct {
	Area    int
	Content [7]string
}

type conf struct {
	Verbs      []string        `yaml:"verbs"`
	Nouns      []string        `yaml:"nouns"`
	Objects    []string        `yaml:"objects"`
	Answers    []string        `yaml:"answers"`
	Locations  []string        `yaml:"locations"`
	Overwrites []mapOverwrites `yaml:"overwrites"`
}

var visibleMap [12][10]int

var objects []string
var locations []string
var overwrites [][7]string

func getMapOverwrites() (overwrites [][7]string) {
	var c conf
	c.getConf("map_overwrites.yaml")
	for _, v := range c.Overwrites {
		// overwrites is already large enough to address v.Area
		if v.Area < len(overwrites) {
			var dummy [7]string
			for i, line := range v.Content {
				dummy[i] = line
			}
			overwrites[v.Area] = dummy
		}
		// overwrites needs expansion to address v.Area
		if v.Area > len(overwrites) {
			var dummy = make([][7]string, v.Area)
			copy(dummy, overwrites)
			overwrites = dummy
		}
		if v.Area == len(overwrites) {
			var dummy [7]string
			for i, line := range v.Content {
				dummy[i] = line
			}
			overwrites = append(overwrites, dummy)
		}
	}
	return
}

func initVisibleAreas() {
	// set all areas to invisible
	for y := 0; y < 12; y++ {
		for x := 0; x < 10; x++ {
			visibleMap[y][x] = 0
		}
	}
	// show first area
	visibleMap[11][0] = 1
}

func areaVisible(area int) bool {
	coordinates := config.AreaCoordinates[area]
	return visibleMap[coordinates.Y][coordinates.X] != 0
}

func revealArea(area int) {
	coordinates := config.AreaCoordinates[area]
	visibleMap[coordinates.Y][coordinates.X] = area
	switch area {
	//case 15:
	//	visibleMap[9][2] = 52
	case 31:
		if areaVisible(15) {
			visibleMap[9][2] = 54
		} else {
			visibleMap[9][2] = 53
		}
		visibleMap[10][2] = 55
	}
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
	*block = append(text, fmt.Sprintf("%s%s%s", color[0], newText, config.WHITE))
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
	return boxLen + 2 // one blank and border left and right
}

func drawBox(area int, boxLen int) (box [7]string) {
	// draw emty field, if area == 0
	if area == 0 {
		// boxlen + left an right connection
		spacer := strings.Repeat(" ", boxLen+4)
		for l := 0; l < 7; l++ {
			box[l] = fmt.Sprintf("%s", spacer)
		}
		return
	}
	// we have an overwrite for this box?
	if len(overwrites) >= area && len(overwrites[area]) > 0 {
		var dummy [7]string
		for i, v := range overwrites[area] {
			dummy[i] = v
		}
		box = dummy
		return
	}
	var leftCon, rightCon, upperCon, topCon, lowerCon, bottomCon string
	// get first line of area from locations
	text := strings.Split(locations[area-1], "\n")[0]
	textLen := len([]rune(text)) + 2 // two space left and right
	leftSpacer := strings.Repeat(" ", (boxLen-textLen)/2)
	conSpacer := strings.Repeat(" ", (boxLen-1)/2)
	rightSpacer := strings.Repeat(" ", boxLen-len(leftSpacer)-textLen)
	empty := strings.Repeat(" ", textLen-2)
	// horizontal line - left/right corner and middle connection element
	horLine := strings.Repeat(config.HL, (boxLen-3)/2)
	// can we walk to the north?
	if config.Areas[area][0] == 0 {
		// no => draw a hoizontal line
		upperCon = " "
		topCon = config.HL
	} else {
		// yes => draw a connection to north
		upperCon = config.VL
		topCon = config.TC
	}
	// can we walk to the south?
	if config.Areas[area][1] == 0 {
		// no => draw a hoizontal line
		bottomCon = config.HL
		lowerCon = " "
	} else {
		// yes => draw a connection to south
		bottomCon = config.BC
		lowerCon = config.VL
	}
	// can we walk to the east?
	if config.Areas[area][2] == 0 {
		// no => draw a vertical line
		rightCon = fmt.Sprintf("%s  ", config.VL)
	} else {
		// yes => draw a connection to west
		rightCon = fmt.Sprintf("%s%s%s", config.RC, config.HL, config.HL)
	}
	// can we walk to the west?
	if config.Areas[area][3] == 0 {
		// no => draw a vertical line
		leftCon = fmt.Sprintf("  %s", config.VL)
	} else {
		// yes => draw a connection to west
		leftCon = fmt.Sprintf("%s%s%s", config.HL, config.HL, config.LC)
	}
	box[0] = fmt.Sprintf("  %s%s%s  ", conSpacer, upperCon, conSpacer)
	box[1] = fmt.Sprintf("  %s%s%s%s%s  ", config.BTL, horLine, topCon, horLine, config.BTR)
	box[2] = fmt.Sprintf("  %s%s%s%s%s  ", config.VL, leftSpacer, empty, rightSpacer, config.VL)
	box[3] = fmt.Sprintf("%s%s%s%s%s", leftCon, leftSpacer, text, rightSpacer, rightCon)
	box[4] = fmt.Sprintf("  %s%s%s%s%s  ", config.VL, leftSpacer, empty, rightSpacer, config.VL)
	box[5] = fmt.Sprintf("  %s%s%s%s%s  ", config.BBL, horLine, bottomCon, horLine, config.BBR)
	box[6] = fmt.Sprintf("  %s%s%s  ", conSpacer, lowerCon, conSpacer)
	return
}

func drawMap(area int) (text []string) {
	coordinates := config.AreaCoordinates[area]
	x := coordinates.X
	y := coordinates.Y
	// max x = 9, don't go further east than 8
	if x > 8 {
		x = 8
	}
	// max y = 11, don't go further south than 8
	//if y > 8 {
	//	y = 8
	//}
	boxLen := getBoxLen(locations)
	//spacer := strings.Repeat(" ", boxLen/2)
	var box1 [7]string
	var box2 [7]string
	var box3 [7]string
	for i := y - 2; i < y+2; i++ {
		// outside y range => draw empty boxes
		if i < 0 || i > 11 {
			box1 = drawBox(0, boxLen)
			box2 = drawBox(0, boxLen)
			box3 = drawBox(0, boxLen)
		} else {
			if x == 0 {
				box1 = drawBox(0, boxLen)
			} else {
				v := visibleMap[i][x-1]
				//if len(overwrites[v]) > 0 {
				//	box1 = overwrites[v]
				//} else {
				box1 = drawBox(v, boxLen)
				//}
			}
			box2 = drawBox(visibleMap[i][x], boxLen)
			if x == 9 {
				box1 = drawBox(0, boxLen)
			} else {
				box3 = drawBox(visibleMap[i][x+1], boxLen)
			}
		}
		/*
			var box2 [7]string
			if len(overwrites[config.AreaMap[i][x+1]]) > 0 {
				box2 = overwrites[config.AreaMap[i][x+1]]
			}
			} else if config.AreaMap[i][x+1] == 8 {
				box2 = overwrites[1].Content
			} else {
				box2 = drawBox(config.AreaMap[i][x+1], boxLen, locations)
			}
			//box3 := drawBox(config.AreaMap[i][x+2], boxLen, locations)
		*/
		for l := 0; l < 7; l++ {
			if i == y {
				text = append(text, fmt.Sprintf("%s%s%s%s%s%s", config.NEUTRAL, box1[l],
					config.YELLOW, box2[l],
					config.NEUTRAL, box3[l]))
			} else {
				text = append(text, fmt.Sprintf("%s%s%s%s", config.NEUTRAL, box1[l], box2[l], box3[l]))
			}
			//fmt.Printf("%s%s\n", box1[l], box2[l])
		}
		/*
			fmt.Printf("%s\u2503%s   %s\u2503%s   %s\u2503%s\n", spacer, spacer, spacer, spacer, spacer, spacer)
			fmt.Printf("%s   %s   %s\n", box1[0], box2[0], box3[0])
			fmt.Printf("%s   %s   %s\n", box1[1], box2[1], box3[1])
			fmt.Printf("%s\u2501\u2501\u25b6%s\u2501\u2513 %s\n", box1[2], box2[2], box3[2])
			fmt.Printf("%s   %s\u25c4\u251B %s\n", box1[3], box2[3], box3[3])
			fmt.Printf("%s   %s   %s\n", box1[4], box2[4], box3[4])
			fmt.Printf("%s\u25B2%s   %s\u2503%s   %s\u2503%s\n", spacer, spacer, spacer, spacer, spacer, spacer)
		*/
	}
	printScreen(text)
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

	appendText(&text, "fantasia", config.RED)
	appendText(&text, "- Ein Adventure von Klaus Hartmuth -", config.BLUE)
	appendText(&text, "- Überarbeitet von Tom Hutter -", config.YELLOW)
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
	flashText = append(text, fmt.Sprintf("%s%s%s", config.RED, err, config.NEUTRAL))
	printScreen(flashText)
	time.Sleep(2 * time.Second)
	printScreen(text)
}

func surroundings(area int, locations []string, objects []string) (text []string) {
	if area == 25 {
		config.ObjectsInArea[40][0] = 25
		objects[40-9] = "eine Tür im Norden"
	}
	if area == 30 {
		config.ObjectsInArea[40][0] = 30
		objects[40-9] = "eine Tür im Süden"
	}

	//	thenge(40)=25:ge$(40)="eine tuer im norden"
	//	ifoa=30thenge(40)=30:ge$(40)="eine tuer im sueden"
	//fmt.Printf("Ich bin %s\n", locations[area-1])
	//var text []string
	text = append(text, fmt.Sprintf("%sIch bin %s", config.YELLOW, locations[area-1]))
	//appendText(&text, fmt.Sprintf("Ich bin %s", locations[area-1]), yellow)
	var items []string
	for i, v := range config.ObjectsInArea {
		if v[0] == area {
			//items = append(items, objects[i-9])
			item := objects[i-9]
			if strings.Contains(item, "::") {
				item = strings.ReplaceAll(item, "::", "")
				items = append(items, fmt.Sprintf("%s  - %s", config.WHITE, item))
				//appendText(&text, fmt.Sprintf("  - %s", item), white)
			} else {
				items = append(items, fmt.Sprintf("%s  - %s", config.BLUE, item))
				//appendText(&text, fmt.Sprintf("  - %s", item), blue)
			}
		}
	}
	if len(items) > 0 {
		text = append(text, "")
		text = append(text, fmt.Sprintf("%sIch sehe:", config.BLUE))
		for _, item := range items {
			text = append(text, item)
		}
		text = append(text, config.NEUTRAL)
	}
	var directions []string
	for d := 0; d < 4; d++ {
		if config.Areas[area][d] != 0 {
			switch d {
			case 0: // N
				directions = append(directions, "Norden")
			case 1: // S
				directions = append(directions, "Süden")
			case 2: // O
				directions = append(directions, "Osten")
			case 3: // W
				directions = append(directions, "Westen")
			}
		}
	}
	text = append(text, fmt.Sprintf("%sRaum: %d, Richtungen: %s", config.WHITE, area, strings.Join(directions, ", ")))
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
	newArea := config.Areas[area][direction]
	if newArea == 0 {
		flash(text, "In diese Richtung führt kein Weg.")
		return area
	}
	// Area 30 and 25 are connected by a door. Is it open?
	if (area == 30 || area == 25 && direction == 0) && !doorOpen {
		flash(text, "Die Tür ist versperrt.")
		return area
	}
	revealArea(newArea)
	return newArea
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
	visibleMap[11][0] = 1
	//var visibleAreas [51]int
	//initVisibleAreas()
	//c.getConf("verbs.yaml")
	//verbs := c.Verbs
	//c.getConf("nouns.yaml")
	//nouns := c.Nouns
	c.getConf("objects.yaml")
	objects = c.Objects
	c.getConf("locations.yaml")
	locations = c.Locations
	//overwrites = getMapOverwrites()

	/*
		for y := 0; y < 12; y++ {
			for x := 0; x < 10; x++ {
				visibleMap[y][x] = config.AreaMap[y][x]
			}
		}
	*/
	//visibleMap[9][2] = 53
	//visibleMap[9][3] = 31
	//visibleMap[10][2] = 54
	//drawMap(0, 7, locations)
	//visibleMap[9][2] = 53
	//drawMap(10, 0, locations, overwrites)
	//return

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
	//text = surroundings(area, locations, objects)
	text = drawMap(1)
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
		drawMap(area)
		//text = surroundings(area, locations, objects)
	}
	//scanner()

	//surroundings(8, locations, objects)
	//scanner()

	//fmt.Println(verbs)
	//fmt.Println(nouns)
	//fmt.Println(objects[0])
	//fmt.Println(locations[0])
}
