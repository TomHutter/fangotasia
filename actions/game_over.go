package actions

import (
	"fangotasia/grid"
	"fangotasia/setup"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GameOver(KO bool) {
	var board []string
	sum := 0
	if KO {
		board = append(board, fmt.Sprintln("G A M E    O V E R"))
	}
	inv := setup.ObjectsInArea(setup.GetAreaByID(setup.INVENTORY))
	if len(inv) > 0 {
		board = append(board, fmt.Sprintln("Du besitzt:"))
	}
	for _, o := range inv {
		val := o.Properties.Value
		if val > 0 {
			sum += val
			desc := strings.Replace(o.Properties.Description.Long, "::", "", -1)
			board = append(board, fmt.Sprintf("- %s: %d Punkte", desc, val))
		}
	}
	board = append(board, fmt.Sprint(""))
	// all valuable objects found
	if sum == 170 {
		switch {
		case moves < 500:
			board = append(board, fmt.Sprintf("- Du hast weniger als %d Züge gebraucht: 7 Punkte", 500))
			sum += 7
			fallthrough
		case moves < 400:
			board = append(board, fmt.Sprintf("- Du hast weniger als %d Züge gebraucht: 7 Punkte", 400))
			sum += 7
			fallthrough
		case moves < 300:
			board = append(board, fmt.Sprintf("- Du hast weniger als %d Züge gebraucht: 7 Punkte", 300))
			sum += 7
		}
	}
	board = append(board, fmt.Sprint(""))
	if setup.GetAreaByID(5).Properties.Visited {
		board = append(board, fmt.Sprintf("- Du bist im Moor gewesen: %d Punkte", 2))
		sum += 2
	}
	if setup.GetAreaByID(29).Properties.Visited {
		board = append(board, fmt.Sprintf("- Du hast die verlassene Burg besucht: %d Punkte", 3))
		sum += 3
	}
	if setup.GetAreaByID(31).Properties.Visited {
		board = append(board, fmt.Sprintf("- Du bist auf einen Baum geklettert: %d Punkte", 4))
		sum += 4
	}
	board = append(board, fmt.Sprint(""))
	board = append(board, fmt.Sprintf("Du hast %d von 200 Punkten!", sum))
	if KO {
		board = append(board, fmt.Sprint("Noch ein Spiel? (j/n)"))
	} else {
		board = append(board, fmt.Sprint("Ach komm, noch 5 Minuten? (j/n)"))
	}
	grid.Grid.Clear()
	grid.Grid.AddItem(grid.AreaGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AreaMap.SetTextAlign(tview.AlignCenter).SetText(strings.Join(board, "\n"))
	grid.AreaField.SetText("")
	grid.App.SetFocus(grid.AreaField)
	grid.AreaField.SetLabel("Weiter (j/n) \u23CE ").
		SetAcceptanceFunc(tview.InputFieldMaxLength(1)).
		SetDoneFunc(func(key tcell.Key) {
			if strings.ToLower(grid.AreaField.GetText()) != "j" {
				grid.App.Stop()
			}
			grid.Grid.Clear()
			grid.Grid.AddItem(grid.InputGrid, 0, 0, 1, 1, 0, 0, false)
			grid.App.SetFocus(grid.InputField)
			grid.Response.SetText("")
		})
	if KO {
		setup.Setup()
	}
}
