package actions

import (
	"fangotasia/grid"
	"fangotasia/setup"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func scoreBoard(yesNo bool, KO bool) {
	var next string
	var board []string
	sum := 0
	if KO {
		board = append(board, fmt.Sprintf("--  %s  --", setup.TextElements["gameOver"][setup.Language]))
	}
	treasure := setup.ObjectsInArea(setup.GetAreaByID(1))
	if len(treasure) > 0 {
		board = append(board, fmt.Sprintf("%s:", setup.TextElements["haveFound"][setup.Language]))
	}
	for _, o := range treasure {
		val := o.Properties.Value
		if val > 0 {
			sum += val
			desc := strings.Replace(o.Properties.Description[setup.Language].Long, "::", "", -1)
			board = append(board, fmt.Sprintf("- %s: %d %s", desc, val, setup.TextElements["points"][setup.Language]))
		}
	}
	board = append(board, fmt.Sprint(""))
	// all valuable objects found
	if sum == 181 {
		switch {
		case setup.Moves < 500:
			board = append(board, fmt.Sprintf(setup.TextElements["movesNeeded"][setup.Language], 500))
			sum += 10
			fallthrough
		case setup.Moves < 400:
			board = append(board, fmt.Sprintf(setup.TextElements["movesNeeded"][setup.Language], 400))
			sum += 10
			fallthrough
		case setup.Moves < 300:
			board = append(board, fmt.Sprintf(setup.TextElements["movesNeeded"][setup.Language], 300))
			sum += 10
		}
	}
	board = append(board, fmt.Sprint(""))
	if setup.Flags["Swamp"] {
		board = append(board, fmt.Sprintf(setup.TextElements["visitedSwamp"][setup.Language], 2))
		sum += 2
	}
	if setup.Flags["Castle"] {
		board = append(board, fmt.Sprintf(setup.TextElements["visitedCastle"][setup.Language], 3))
		sum += 3
	}
	if setup.Flags["Tree"] {
		board = append(board, fmt.Sprintf(setup.TextElements["climbedTree"][setup.Language], 4))
		sum += 4
	}
	board = append(board, fmt.Sprint(""))
	board = append(board, fmt.Sprintf(setup.TextElements["pointsOutOf"][setup.Language], sum))

	if yesNo {
		if KO {
			board = append(board, fmt.Sprint(setup.TextElements["oneMoreGame"][setup.Language]))
		} else {
			board = append(board, fmt.Sprint(setup.TextElements["comeOn"][setup.Language]))
		}
	}
	if KO {
		setup.Setup()
	}
	grid.Grid.Clear()
	grid.Grid.AddItem(grid.AreaGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AreaMap.SetTextAlign(tview.AlignCenter).SetText(strings.Join(board, "\n"))
	grid.AreaField.SetText("")
	grid.App.SetFocus(grid.AreaField)
	if yesNo {
		next = fmt.Sprintf("%s \u23CE ", setup.TextElements["nextYesNo"][setup.Language])
	} else {
		next = fmt.Sprintf("%s \u23CE ", setup.TextElements["next"][setup.Language])
	}
	grid.AreaField.SetLabel(next).
		SetAcceptanceFunc(tview.InputFieldMaxLength(1)).
		SetDoneFunc(func(key tcell.Key) {
			if yesNo && strings.ToLower(grid.AreaField.GetText()) != setup.TextElements["yesShort"][setup.Language] {
				grid.App.Stop()
			}
			grid.Grid.Clear()
			grid.Grid.AddItem(grid.InputGrid, 0, 0, 1, 1, 0, 0, false)
			grid.App.SetFocus(grid.InputField)
			grid.Response.SetText("")
		})
}
