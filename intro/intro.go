package intro

import (
	"fangotasia/grid"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func Intro() {
	var text = []string{
		"[blue:black:-]Mach Dich auf den gefahrenreichen Weg in",
		"das zauberhafte Land Fangotasia und suche",
		"nach märchenhaften Schätzen.",
		"Führe mich mit einfachen Kommandos in",
		"einem oder zwei Worten, z.B.:",
		"",
		"[yellow:black:b]SPRING      BENUTZE TARNKAPPE      ENDE",
		"",
		"LEGE RUBIN     FÜTTERE DRACHE     INVENTAR",
		"",
		"[blue:black:-]Mit  [white:black:b]SAVE  [blue:black:-]kannst Du den aktuellen Stand",
		"des Spieles abspeichern,",
		"mit  [white:black:b]LOAD  [blue:black:-]wieder einlesen.",
	}
	//grid.Grid.Clear()
	//grid.Grid.AddItem(grid.AreaGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AreaMap.SetText(strings.Join(text, "\n"))
	//grid.AreaField.SetText("")
	//grid.App.SetFocus(grid.AreaField)
	grid.AreaField.SetDoneFunc(func(key tcell.Key) {
		grid.Grid.Clear()
		grid.Grid.AddItem(grid.InputGrid, 0, 0, 1, 1, 0, 0, false)
		grid.App.SetFocus(grid.InputField)
	})
	//view.PrintScreen(text)
	//view.Scanner("once: true")
}

func Prelude() {
	var text = []string{
		"[red:black:b]F A N G O T A S I A",
		"",
		"[blue:black:b]- Ein Adventure von Klaus Hartmuth -",
		"",
		"[yellow:black:b]- GO Version von Tom Hutter -",
	}
	grid.Grid.Clear()
	grid.Grid.AddItem(grid.AreaGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AreaMap.SetText(strings.Join(text, "\n"))
	grid.AreaField.SetText("")
	grid.App.SetFocus(grid.AreaField)
	grid.AreaField.SetDoneFunc(func(key tcell.Key) { Intro() })
	//grid.Pages.SwitchToPage("map")

	//view.PrintScreen(text)
	//view.Scanner("once: true")
}
