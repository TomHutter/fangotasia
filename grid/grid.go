package grid

import (
	"fangotasia/setup"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	App          *tview.Application
	Grid         *tview.Grid
	InputGrid    *tview.Grid
	InputField   *tview.InputField
	Surroundings *tview.TextView
	Response     *tview.TextView
	AreaGrid     *tview.Grid
	AreaField    *tview.InputField
	AreaMap      *tview.TextView
)

var Input = make(chan string, 1)

func SetupGrid() {
	App = tview.NewApplication()

	InputField = tview.NewInputField().
		SetLabel(fmt.Sprintf("%s? > ", setup.TextElements["andNow"])).
		SetLabelColor(tcell.ColorDarkCyan).
		SetFieldWidth(80).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorDarkCyan).
		SetDoneFunc(func(key tcell.Key) {
			Input <- InputField.GetText()
		})

	AreaField = tview.NewInputField().
		SetLabel(fmt.Sprintf("%s \u23CE ", setup.TextElements["next"])).
		SetLabelColor(tcell.ColorDarkCyan).
		SetFieldWidth(20).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorDarkCyan).
		SetAcceptanceFunc(tview.InputFieldMaxLength(0))

	Surroundings = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	AreaMap = tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetWrap(false).
		SetChangedFunc(func() {
			App.Draw()
		})

	Response = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			App.Draw()
		})

	InputGrid = tview.NewGrid().
		SetRows(0, 1, 0).
		SetColumns(0).
		SetBorders(false).
		AddItem(Surroundings, 0, 0, 15, 1, 0, 0, false).
		AddItem(InputField, 15, 0, 1, 1, 0, 0, false).
		AddItem(Response, 16, 0, 15, 1, 0, 0, false)

	AreaGrid = tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0).
		SetBorders(false).
		AddItem(AreaMap, 0, 0, 1, 1, 0, 0, false).
		AddItem(AreaField, 1, 0, 1, 1, 0, 0, false)

	Grid = tview.NewGrid().
		SetRows(0).
		SetColumns(0).
		SetBorders(true)

}
