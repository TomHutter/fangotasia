package intro

import (
	"fangotasia/grid"
	"fangotasia/setup"

	"github.com/gdamore/tcell/v2"
)

func Intro() {
	var text = setup.TextElements["intro"]
	grid.AreaMap.SetText(text)
	grid.AreaField.SetDoneFunc(func(key tcell.Key) {
		grid.Grid.Clear()
		grid.Grid.AddItem(grid.InputGrid, 0, 0, 1, 1, 0, 0, false)
		grid.App.SetFocus(grid.InputField)
	})
}

func Prelude() {
	var text = setup.TextElements["prelude"]
	grid.Grid.Clear()
	grid.Grid.AddItem(grid.AreaGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AreaMap.SetText(text)
	grid.AreaField.SetText("")
	grid.App.SetFocus(grid.AreaField)
	grid.AreaField.SetDoneFunc(func(key tcell.Key) { Intro() })
}
