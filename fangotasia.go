package main

import (
	"fangotasia/actions"
	"fangotasia/grid"
	"fangotasia/intro"
	"fangotasia/movement"
	"fangotasia/setup"
	"fangotasia/view"
	"os"
	"strings"
)

/*
ToDos:
- forest mapping buggy
*/

/*
func recoverFromTviewPanic() {
	if r := recover(); r != nil {
		fmt.Println("recovered from ", r)
        (actions.Object{}).Save(area setup.Area) (r setup.Reaction) {
	}
	setup.PathName, _ = os.Getwd()
}
*/

func init() {
	setup.PathName, _ = os.Getwd()
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
