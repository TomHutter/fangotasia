package main

import (
	"fantasia/actions"
	"fantasia/movement"
	"fantasia/setup"
	"fantasia/view"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unicode/utf8"
)

/*
ToDos:
- func Jump
- func Help
- func Chase
*/

// Setup keyboard scanning
func scanner() (r rune) {
	var b []byte = make([]byte, 4)
	os.Stdin.Read(b)
	r, _ = utf8.DecodeRune(b)
	//fmt.Println(b)
	//fmt.Println(r)
	return
}

func intro() {
	var text = []string{
		"\033[01;34mMach Dich auf den gefahrenreichen Weg in",
		"das zauberhafte Land Fantasia und suche",
		"nach märchenhaften Schätzen.",
		"Führe mich mit einfachen Kommandos in",
		"einem oder zwei Worten, z.B.:",
		"",
		"\033[01;33mNORDEN      BENUTZE TARNKAPPE      ENDE",
		"",
		"LEGE RUBIN     TÖTE DRACHE     INVENTAR",
		"",
		"\033[01;34mMit  \033[01;97mSAVE  \033[01;34mkannst Du den aktuellen Stand",
		"des Spieles abspeichern,",
		"mit  \033[01;97mLOAD  \033[01;34mwieder einlesen.",
		"\033[0mWeiter \u23CE",
	}
	view.PrintScreen(text)
	view.Scanner("once: true")
}

func prelude() {
	var text = []string{
		"\033[01;31mF A N T A S I A",
		"",
		"\033[01;34m- Ein Adventure von Klaus Hartmuth -",
		"",
		"\033[01;33m- GO Version von Tom Hutter -",
		"",
		"\033[0mWeiter \u23CE",
	}
	view.PrintScreen(text)
	view.Scanner("once: true")
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
func init() {
	setup.PathName, _ = os.Getwd()
	//_, filename, _, _ := runtime.Caller(0)
	//setup.PathName = path.Dir(filename)
}

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
	setup.Setup()
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
	intro()
	area := setup.GetAreaByID(1)
	//oldArea := area
	movement.RevealArea(area.ID)
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
	//area = actions.Parse("klettere baum", area, text)
	//text = movement.DrawMap(area)
	//surroundings = movement.Surroundings(area)
	//text = append(text, surroundings...)
	//view.PrintScreen(text)
	//area = actions.Parse("klettere baum", area, text)
	//text = movement.DrawMap(area)
	//surroundings = movement.Surroundings(area)
	//text = append(text, surroundings...)
	//view.PrintScreen(text)
	//actions.Parse("inventar", area)
	//actions.Parse("load", area, text)
	//actions.Parse("nimm Zwergendolch", area, text)
	//actions.Parse("Inventar", area, text)
	//actions.Parse("n", area, text)
	//actions.Parse("o", area, text)
	//actions.Parse("s", area, text)
	//actions.Parse("stich Gnom", area, text)
	//actions.Parse("Inventar", area, text)
	//surroundings = movement.Surroundings(area)
	actions.Parse("nimm zauberschuhe", area, text)
	actions.Parse("trage zauberschuhe", area, text)
	area = actions.Parse("o", area, text)
	area = actions.Parse("w", area, text)
	area = actions.Parse("o", area, text)
	area = actions.Parse("n", area, text)
	area = actions.Parse("o", area, text)
	area = actions.Parse("klettere baum", area, text)
	area = actions.Parse("klettere baum", area, text)
	area = actions.Parse("klettere baum", area, text)
	area = actions.Parse("klettere baum", area, text)
	area = actions.Parse("w", area, text)
	/*
		actions.Parse("Inventar", area, text)
		area = actions.Parse("w", area, text)
		area = actions.Parse("w", area, text)
		text = movement.DrawMap(area)
		surroundings = movement.Surroundings(area)
		text = append(text, surroundings...)
		view.PrintScreen(text)
		area = actions.Parse("sag simsalabim", area, text)
		area = actions.Parse("benutze karte", area, text)
		area = actions.Parse("füttere tafel", area, text)
	*/
	/*
		area = actions.Parse("load", area, text)
		area = actions.Parse("gieße strauch", area, text)
		area = actions.Parse("gieße strauch", area, text)
	*/

	/*
		actions.Parse("nimm zauberschuhe", area, text)
		actions.Parse("trage zauberschuhe", area, text)
		area = actions.Parse("o", area, text)
		text = movement.DrawMap(area)
		surroundings = movement.Surroundings(area)
		text = append(text, surroundings...)
		view.PrintScreen(text)
		area = actions.Parse("o", area, text)
		text = movement.DrawMap(area)
		surroundings = movement.Surroundings(area)
		text = append(text, surroundings...)
		view.PrintScreen(text)
		actions.Parse("nimm zwergendolch", area, text)
		area = actions.Parse("n", area, text)
		text = movement.DrawMap(area)
		surroundings = movement.Surroundings(area)
		text = append(text, surroundings...)
		view.PrintScreen(text)
	*/
	/*
		actions.Parse("klettere baum", area, text)
		actions.Parse("klettere baum", area, text)
		area = actions.Parse("o", area, text)
		text = movement.DrawMap(area)
		surroundings = movement.Surroundings(area)
		text = append(text, surroundings...)
		view.PrintScreen(text)
		area = actions.Parse("s", area, text)
		text = movement.DrawMap(area)
		surroundings = movement.Surroundings(area)
		text = append(text, surroundings...)
		view.PrintScreen(text)
		area = actions.Parse("o", area, text)
		text = movement.DrawMap(area)
		surroundings = movement.Surroundings(area)
		text = append(text, surroundings...)
		view.PrintScreen(text)
	*/
	//actions.Parse("stich gnom", area, text)
	//actions.Parse("n", area, text)
	//actions.Parse("öffne	Tür", area)
	//actions.Parse(view.Scanner("prompt: und nun? > "), area)
	for {
		area = actions.Parse(view.Scanner("prompt: und nun? > "), area, text)
		text = movement.DrawMap(area)
		surroundings = movement.Surroundings(area)
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
