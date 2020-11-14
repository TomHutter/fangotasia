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
	Content [3]string
}

type places struct {
	Long  string
	Short string
}

type conf struct {
	Verbs      []string        `yaml:"verbs"`
	Nouns      []string        `yaml:"nouns"`
	Objects    []string        `yaml:"objects"`
	Answers    []string        `yaml:"answers"`
	Locations  []places        `yaml:"locations"`
	Overwrites []mapOverwrites `yaml:"overwrites"`
}

var visibleMap [12][10]int

var objects []string
var locations []places
var overwrites [][3]string

func getMapOverwrites() (overwrites [][3]string) {
	var c conf
	c.getConf("map_overwrites.yaml")
	for _, v := range c.Overwrites {
		// overwrites is already large enough to address v.Area
		if v.Area < len(overwrites) {
			var dummy [3]string
			for i, line := range v.Content {
				dummy[i] = line
			}
			overwrites[v.Area] = dummy
		}
		// overwrites needs expansion to address v.Area
		if v.Area > len(overwrites) {
			var dummy = make([][3]string, v.Area)
			copy(dummy, overwrites)
			overwrites = dummy
		}
		if v.Area == len(overwrites) {
			var dummy [3]string
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
	case 5:
		if areaVisible(36) {
			visibleMap[11][4] = 57
			visibleMap[11][5] = 58
			visibleMap[5][5] = 59
		}
	case 6:
		if areaVisible(7) {
			visibleMap[9][1] = 52
		}
	case 7:
		if areaVisible(6) {
			visibleMap[9][1] = 52
		}
	case 9:
		if areaVisible(31) {
			visibleMap[10][2] = 54
		}
	case 15:
		if areaVisible(31) {
			visibleMap[9][2] = 55
		} else {
			visibleMap[9][2] = 53
		}
	case 31:
		visibleMap[10][2] = 54
		if areaVisible(15) {
			visibleMap[9][2] = 55
		} else {
			visibleMap[9][2] = 56
		}
	case 32:
		if areaVisible(37) {
			visibleMap[4][5] = 60
		} else {
			visibleMap[4][5] = 53
		}
		if visibleMap[11][5] != 0 {
			visibleMap[5][5] = 59
		} else {
			visibleMap[5][5] = 53
		}
	case 37:
		visibleMap[4][5] = 60
		if areaVisible(40) {
			visibleMap[3][6] = 61
			visibleMap[4][6] = 62
		}
	case 38:
		if areaVisible(40) {
			visibleMap[5][6] = 63
			visibleMap[6][6] = 64
		} else {
			visibleMap[5][6] = 0
		}
	case 39:
		if areaVisible(40) {
			visibleMap[4][7] = 65
			visibleMap[5][7] = 64
		} else {
			visibleMap[4][7] = 0
		}
	case 40:
		visibleMap[3][6] = 61
		visibleMap[4][6] = 62
		visibleMap[5][6] = 63
		visibleMap[6][6] = 64
		visibleMap[4][7] = 65
		visibleMap[5][7] = 64
	case 41:
		if areaVisible(51) {
			visibleMap[2][6] = 66
			if areaVisible(40) {
				visibleMap[3][6] = 64
			} else {
				visibleMap[3][6] = 67
			}
		}
	case 42:
		if areaVisible(51) {
			visibleMap[2][7] = 68
		} else {
			visibleMap[2][7] = 42
		}
	case 43:
		if areaVisible(51) {
			visibleMap[1][6] = 69
		} else {
			visibleMap[1][6] = 0
		}
	case 44:
		if areaVisible(51) {
			visibleMap[1][7] = 70
		} else {
			visibleMap[1][7] = 0
		}
	case 45:
		if areaVisible(51) {
			visibleMap[1][8] = 71
			visibleMap[2][8] = 64
		} else {
			visibleMap[1][8] = 0
			visibleMap[2][8] = 0
		}
	case 46:
		if !areaVisible(51) {
			visibleMap[1][9] = 0
		}
	case 47:
		if areaVisible(51) {
			visibleMap[0][6] = 72
		} else {
			visibleMap[0][6] = 0
		}
	case 48:
		if areaVisible(51) {
			visibleMap[0][7] = 73
		} else {
			visibleMap[0][7] = 0
		}
	case 49:
		if !areaVisible(51) {
			visibleMap[0][8] = 0
		}
	case 50:
		if !areaVisible(51) {
			visibleMap[0][9] = 0
		}
	case 51:
		visibleMap[2][6] = 66
		visibleMap[3][6] = 67
		visibleMap[2][7] = 68
		visibleMap[2][7] = 68
		visibleMap[1][6] = 69
		visibleMap[1][7] = 70
		visibleMap[1][8] = 71
		visibleMap[2][8] = 64
		visibleMap[0][6] = 72
		visibleMap[0][7] = 73
		visibleMap[1][9] = 46
		visibleMap[0][8] = 49
		visibleMap[0][9] = 50
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

func getBoxLen(locations []places) int {
	boxLen := 0
	for _, v := range locations {
		lineLen := len([]rune(strings.Split(v.Short, "\n")[0]))
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

func drawBox(area int, boxLen int) (box [3]string) {
	// draw emty field, if area == 0
	if area == 0 {
		// boxlen + left an right connection
		spacer := strings.Repeat(" ", boxLen+2)
		for l := 0; l < 3; l++ {
			box[l] = fmt.Sprintf("%s", spacer)
		}
		return
	}
	// we have an overwrite for this box?
	if len(overwrites) >= area && len(overwrites[area][0]) > 0 {
		var dummy [3]string
		for i, v := range overwrites[area] {
			dummy[i] = v
		}
		box = dummy
		return
	}
	var leftCon, rightCon, topCon, bottomCon string
	// get first line of area from locations
	text := strings.Split(locations[area-1].Short, "\n")[0]
	textLen := len([]rune(text)) + 2 // two space left and right
	leftSpacer := strings.Repeat(" ", (boxLen-textLen)/2)
	rightSpacer := strings.Repeat(" ", boxLen-len(leftSpacer)-textLen)
	// horizontal line - left/right corner and middle connection element
	horLine := strings.Repeat(config.HL, (boxLen-3)/2)
	// can we walk to the north?
	if config.Areas[area][0] == 0 {
		// no => draw a hoizontal line
		topCon = config.HL
	} else {
		// yes => draw a connection to north
		topCon = config.TC
	}
	// can we walk to the south?
	if config.Areas[area][1] == 0 {
		// no => draw a hoizontal line
		bottomCon = config.HL
	} else {
		// yes => draw a connection to south
		bottomCon = config.BC
	}
	// can we walk to the east?
	if config.Areas[area][2] == 0 {
		// no => draw a vertical line
		rightCon = fmt.Sprintf("%s ", config.VL)
	} else {
		// yes => draw a connection to west
		rightCon = fmt.Sprintf("%s%s", config.RC, config.HL)
	}
	// can we walk to the west?
	if config.Areas[area][3] == 0 {
		// no => draw a vertical line
		leftCon = fmt.Sprintf(" %s", config.VL)
	} else {
		// yes => draw a connection to west
		leftCon = fmt.Sprintf("%s%s", config.HL, config.LC)
	}
	box[0] = fmt.Sprintf(" %s%s%s%s%s ", config.BTL, horLine, topCon, horLine, config.BTR)
	box[1] = fmt.Sprintf("%s%s%s%s%s", leftCon, leftSpacer, text, rightSpacer, rightCon)
	box[2] = fmt.Sprintf(" %s%s%s%s%s ", config.BBL, horLine, bottomCon, horLine, config.BBR)
	return
}

func drawMap(area int) (text []string) {
	coordinates := config.AreaCoordinates[area]
	x := coordinates.X
	y := coordinates.Y
	// max x = 9, don't go further east than 8
	/*
		if x > 8 {
			x = 8
		}
	*/
	boxLen := getBoxLen(locations)
	var boxes [5][3]string
	for i := 0; i < 6; i++ {
		iy := y + i - 2
		// outside y range => draw empty boxes
		if iy < 0 || iy > 11 {
			for j := 0; j < 5; j++ {
				boxes[j] = drawBox(0, boxLen)
			}
		} else {
			for j := 0; j < 5; j++ {
				ix := x + j - 2
				if ix < 0 || ix > 9 {
					boxes[j] = drawBox(0, boxLen)
				} else {
					v := visibleMap[iy][ix]
					boxes[j] = drawBox(v, boxLen)
				}
			}
		}
		for l := 0; l < 3; l++ {
			if iy == y {
				text = append(text, fmt.Sprintf("%s%s%s%s%s%s%s%s", config.NEUTRAL, boxes[0][l], boxes[1][l],
					config.YELLOW, boxes[2][l],
					config.NEUTRAL, boxes[3][l], boxes[4][l]))
			} else {
				text = append(text, fmt.Sprintf("%s%s%s%s%s%s", config.NEUTRAL, boxes[0][l],
					boxes[1][l], boxes[2][l],
					boxes[3][l], boxes[4][l]))
			}
		}
	}
	//printScreen(text)
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

func surroundings(area int, locations []places, objects []string) (text []string) {
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
	text = append(text, fmt.Sprintf("%sIch bin %s", config.YELLOW, locations[area-1].Long))
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
	//printScreen(text)
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
	/*
		visibleMap[11][1] = 2
		visibleMap[11][2] = 3
		visibleMap[10][1] = 8
		visibleMap[10][2] = 9
		visibleMap[9][0] = 6
		visibleMap[9][1] = 7
		visibleMap[9][3] = 52
	*/

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
	overwrites = getMapOverwrites()

	/*
		for y := 0; y < 12; y++ {
			for x := 0; x < 10; x++ {
				visibleMap[y][x] = config.AreaMap[y][x]
			}
		}
	*/
	//visibleMap[11][0] = 1
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
	area := 40
	oldArea := area
	revealArea(area)
	var dir rune
	var direction int
	var text []string
	//text = surroundings(area, locations, objects)
	text = drawMap(area)
	text = append(text, "\n", "\n", "\n")
	text = append(text, surroundings(area, locations, objects)...)
	printScreen(text)
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
		// are we lost? (show old area)
		if !areaVisible(area) {
			text = drawMap(oldArea)
			text = append(text, "\n", "\n", "\n")
			text = append(text, surroundings(oldArea, locations, objects)...)
			printScreen(text)
		} else {
			//text = drawMap(area)
			//text = surroundings(area, locations, objects)
			text = drawMap(area)
			text = append(text, "\n", "\n", "\n")
			text = append(text, surroundings(area, locations, objects)...)
			oldArea = area
			printScreen(text)
		}
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
