package main

import (
	"fantasia/actions"
	"fantasia/config"
	"fantasia/movement"
	"fantasia/view"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unicode/utf8"
)

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

	view.AppendText(&text, "fantasia", config.RED)
	view.AppendText(&text, "- Ein Adventure von Klaus Hartmuth -", config.BLUE)
	view.AppendText(&text, "- GO Version von Tom Hutter -", config.YELLOW)
	view.PrintScreen(text)
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

func main() {
	//var c conf
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
	config.Init()
	/*
		c.getConf("config/objects.yaml")
		objects = c.Objects
		c.getConf("config/locations.yaml")
		locations = c.Locations
		overwrites = getMapOverwrites()
	*/

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
	view.Scanner("once: true")
	area := 1
	//oldArea := area
	movement.RevealArea(area)
	//var dir rune
	//var direction int
	//var text []string
	//text = surroundings(area, locations, objects)
	text := movement.DrawMap(area)
	//text = append(text, "\n", "\n", "\n")
	surroundings := movement.Surroundings(area)
	text = append(text, surroundings...)
	//view.Input()
	view.PrintScreen(text)
	//actions.Parse()
	//actions.Parse("verben", area, text)
	//actions.Parse("inventar", area)
	//actions.Parse("nimm Zwergendolch", area, text)
	//actions.Parse("Inventar", area, text)
	//actions.Parse("n", area, text)
	//actions.Parse("o", area, text)
	//actions.Parse("s", area, text)
	//actions.Parse("stich Gnom", area, text)
	//actions.Parse("Inventar", area, text)
	//surroundings = movement.Surroundings(area)
	//actions.Parse("nimm Gnom", area, text)
	//actions.Parse("n", area, text)
	//actions.Parse("öffne	Tür", area)
	//actions.Parse(view.Scanner("prompt: und nun? > "), area)
	for {
		area = actions.Parse(view.Scanner("prompt: und nun? > "), area, text)
		text = movement.DrawMap(area)
		surroundings := movement.Surroundings(area)
		text = append(text, surroundings...)
		view.PrintScreen(text)
		//actions.Parse("verben")
		//fmt.Printf("\ninput: %s\n", string(input))
		/*
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
				area = movement.Move(area, direction, text)
				// are we lost? (show old area)
				if !movement.AreaVisible(area) {
					text = movement.DrawMap(oldArea)
					//text = append(text, "\n", "\n", "\n")
					text = append(text, movement.Surroundings(oldArea)...)
					view.PrintScreen(text)
				} else {
					//text = drawMap(area)
					//text = surroundings(area, locations, objects)
					text = movement.DrawMap(area)
					//text = append(text, "\n", "\n", "\n")
					text = append(text, movement.Surroundings(area)...)
					oldArea = area
					view.PrintScreen(text)
				}
			//text = surroundings(area, locations, objects)
		*/
	}
	//scanner()

	//surroundings(8, locations, objects)
	//scanner()

	//fmt.Println(verbs)
	//fmt.Println(nouns)
	//fmt.Println(objects[0])
	//fmt.Println(locations[0])
}
