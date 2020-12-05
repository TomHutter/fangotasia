package main

import (
	"fantasia/actions"
	"fantasia/intro"
	"fantasia/movement"
	"fantasia/setup"
	"fantasia/view"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetupGrid(area setup.Area) {
	app := tview.NewApplication()
	//SetAcceptanceFunc(tview.InputFieldInteger).
	//if err := app.SetRoot(inputField, true).SetFocus(inputField).Run(); err != nil {
	//	panic(err)
	//}
	/*
		newPrimitive := func(text string) *tview.TextView {
			return tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText(text)
		}
	*/
	//menu := newPrimitive("Menu")
	//main := newPrimitive("Main content")
	//sideBar := newPrimitive("Side Bar")
	inputField := tview.NewInputField()
	areaField := tview.NewInputField()

	surroundings := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetText(strings.Join(view.Surroundings(area), "\n"))

	areaMap := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
		//SetText(strings.Join(movement.DrawMap(area), "\n"))

	response := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	inputGrid := tview.NewGrid().
		SetRows(0, 1, 0).
		SetColumns(0).
		SetBorders(false).
		AddItem(surroundings, 0, 0, 15, 1, 0, 0, false).
		AddItem(inputField, 15, 0, 1, 1, 0, 0, false).
		AddItem(response, 16, 0, 15, 1, 0, 0, false)
		//AddItem(areaMap, 0, 0, 1, 3, 0, 0, false).

	areaGrid := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0).
		SetBorders(false).
		AddItem(areaMap, 0, 0, 1, 1, 0, 0, false).
		AddItem(areaField, 1, 0, 1, 1, 0, 0, false)

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(0).
		SetBorders(true).
		AddItem(inputGrid, 0, 0, 1, 1, 0, 0, false)

	inputField.
		SetLabel("und nun? > ").
		SetLabelColor(tcell.ColorDarkCyan).
		SetFieldWidth(80).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorDarkCyan).
		SetDoneFunc(func(key tcell.Key) {
			//app.SetFocus(sourroundings)
			//fmt.Println(inputField.GetText())
			area = actions.Parse(inputField.GetText(), area, []string{})
			inputField.SetText("")
			//response.Clear()
			//surroundings.Clear()
			if area.ID == 0 {
				areaMap.SetText(strings.Join(movement.DrawMap(area), "\n"))
				grid.Clear()
				grid.AddItem(areaGrid, 0, 0, 1, 1, 0, 0, false)
				areaField.SetText("")
				app.SetFocus(areaField)
			} else {
				surroundings.SetText(strings.Join(view.Surroundings(area), "\n"))
				if len(view.Notice.Message) != 0 {
					//var notice string
					response.SetText(
						fmt.Sprintf("\n%s%s%s\n",
							view.Notice.Color,
							view.Notice.Message, "[-:black:-]"))
					//response.SetText("")
					view.Notice.Message = ""
					view.Notice.Color = ""
					view.Notice.Sleep = 0
				}
			}
			//areaMap.Clear()
			//text := DrawMap(area)
			//surroundings = view.Surroundings(area)
			//areaMap.SetText(strings.Join(movement.DrawMap(area), "\n"))
			//response.SetText(strings.Join(view.Surroundings(area), "\n"))
			//app.SetFocus(inputField)
		})

		/*
			areaField.
				SetLabel("und nun? > ").
				SetLabelColor(tcell.ColorDarkCyan).
				SetFieldWidth(80).
				SetFieldBackgroundColor(tcell.ColorBlack).
				SetFieldTextColor(tcell.ColorDarkCyan).
				SetDoneFunc(func(key tcell.Key) {
					grid.RemoveItem(inputGrid)
					grid.AddItem(areaGrid, 0, 0, 1, 1, 0, 0, false)
					//app.SetFocus(sourroundings)
					//fmt.Println(inputField.GetText())
					area = actions.Parse(inputField.GetText(), area, []string{})
					//text := DrawMap(area)
					//surroundings = view.Surroundings(area)
					areaMap.SetText(strings.Join(movement.DrawMap(area), "\n"))
					response.SetText(strings.Join(view.Surroundings(area), "\n"))
					inputField.SetText("")
					//app.SetFocus(inputField)
				})
		*/

	areaField.
		SetLabel("Weiter \u23CE ").
		SetLabelColor(tcell.ColorDarkCyan).
		SetFieldWidth(20).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorDarkCyan).
		SetAcceptanceFunc(tview.InputFieldMaxLength(0)).
		SetDoneFunc(func(key tcell.Key) {
			//fmt.Println("blubb")
			//areaMap.SetText("blubb")
			areaMap.Clear()
			grid.Clear()
			//grid.RemoveItem(areaGrid)
			grid.AddItem(inputGrid, 0, 0, 1, 1, 0, 0, false)
			app.SetFocus(inputField)
			//surroundings.Clear()
			//response.Clear()
			//areaMap.SetText(strings.Join(movement.DrawMap(area), "\n"))
			surroundings.SetText(strings.Join(view.Surroundings(area), "\n"))
			//response.SetText(strings.Join(view.Surroundings(area), "\n"))
			inputField.SetText("")
			//grid.AddItem(inputGrid, 0, 0, 1, 1, 0, 0, false)
			//app.SetFocus(sourroundings)
			//fmt.Println(inputField.GetText())
			//area = actions.Parse(areaField.GetText(), area, []string{})
			//text := DrawMap(area)
			//surroundings = view.Surroundings(area)
			//response.SetText(strings.Join(view.Surroundings(area), "\n"))
			//inputField.SetText("")
			//app.SetFocus(inputField)
		})

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	/*
		grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
			AddItem(main, 1, 0, 1, 3, 0, 0, false).
			AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

		// Layout for screens wider than 100 cells.
		grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
			AddItem(main, 1, 1, 1, 1, 0, 100, false).
			AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)
	*/
	if err := app.SetRoot(grid, true).SetFocus(inputField).Run(); err != nil {
		panic(err)
	}
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
	setupCloseHandler()
	SetupGrid(setup.GetAreaByID(1))
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
	intro.Prelude()
	intro.Intro()
	area := setup.GetAreaByID(1)
	//oldArea := area
	movement.RevealArea(area.ID)
	//var dir rune
	//var direction int
	//var text []string
	//text = surroundings(area, locations, objects)
	text := movement.DrawMap(area)
	//text = append(text, "\n", "\n", "\n")
	surroundings := view.Surroundings(area)
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
	for {
		area = actions.Parse(view.Scanner("prompt: und nun? > "), area, text)
		text = movement.DrawMap(area)
		surroundings = view.Surroundings(area)
		text = append(text, surroundings...)
		view.PrintScreen(text)
	}
}
