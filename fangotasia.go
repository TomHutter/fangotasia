package main

import (
	"fangotasia/actions"
	"fangotasia/grid"
	"fangotasia/intro"
	"fangotasia/movement"
	"fangotasia/setup"
	"fangotasia/view"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
ToDos:
- forest mapping buggy
*/

// init() will be called before main() by go convention
func init() {
	setup.PathName, _ = os.Getwd()
}

func SetupLanguage() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
	grid.LanguageSelect.
		SetLabel(setup.TextElements["selectLanguage"]).
		SetOptions([]string{"First", "Second", "Third", "Fourth", "Fifth"}, nil)
}

func main() {
	setup.Setup()
	grid.SetupGrid()
	area := setup.GetAreaByID(1)
	movement.RevealArea(area.ID)
	go actions.REPL(area)
	grid.Surroundings.SetText(strings.Join(view.Surroundings(area), "\n"))
	intro.Prelude()
	if err := grid.App.SetRoot(grid.Grid, true).SetFocus(grid.AreaField).Run(); err != nil {
		panic(err)
	}
}
