package main

import (
	"fangotasia/actions"
	"fangotasia/grid"
	"fangotasia/intro"
	"fangotasia/movement"
	"fangotasia/setup"
	"fangotasia/view"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
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

func init() {
	setup.PathName, _ = os.Getwd()
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
	//setupCloseHandler()
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
	//exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	//exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// reenable display entered characters on the screen
	//defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	// Setup our Ctrl+C handler
	grid.SetupGrid()
	//intro.Intro()
	//grid.REPL(setup.GetAreaByID(1))
	area := setup.GetAreaByID(1)
	//oldArea := area
	movement.RevealArea(area.ID)
	go actions.REPL(area)
	grid.Surroundings.SetText(strings.Join(view.Surroundings(area), "\n"))
	intro.Prelude(func() { actions.REPL(area) })
	//for {
	//actions.REPL(area)
	//}
	//var dir rune
	//var direction int
	//var text []string
	//text = surroundings(area, locations, objects)
	//text := movement.DrawMap(area)
	//text = append(text, "\n", "\n", "\n")
	//surroundings := view.Surroundings(area)
	//text = append(text, surroundings...)
	//view.Input()
	//view.PrintScreen(text)
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
	/*
		actions.Parse("Inventar", area, text)
		actions.Parse("nimm zauberschuhe", area, text)
		actions.Parse("Inventar", area, text)
		actions.Parse("trage zauberschuhe", area, text)
		area = actions.Parse("o", area, text)
		area = actions.Parse("w", area, text)
		area = actions.Parse("o", area, text)
		area = actions.Parse("n", area, text)
		area = actions.Parse("o", area, text)
		area = actions.Parse("klettere baum", area, text)
		area = actions.Parse("klettere baum", area, text)
		area = actions.Parse("klettere baum", area, text)
	*/
	//area = actions.Parse("spring", area, text)
	//area = actions.Parse("klettere baum", area, text)
	//area = actions.Parse("w", area, text)
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
	/*
		for {
			area = actions.Parse(view.Scanner("prompt: und nun? > "), area, text)
			text = movement.DrawMap(area)
			surroundings = view.Surroundings(area)
			text = append(text, surroundings...)
			view.PrintScreen(text)
		}
	*/
	if err := grid.App.SetRoot(grid.Grid, true).SetFocus(grid.AreaField).Run(); err != nil {
		//if err := grid.App.SetRoot(grid.Pages, true).SetFocus(grid.Pages).Run(); err != nil {
		panic(err)
	}
}
