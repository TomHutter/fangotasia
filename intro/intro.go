package intro

import (
	"fangotasia/grid"
	"fangotasia/setup"
	"fangotasia/view"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func Intro() {
	var text = setup.TextElements["intro"][setup.Language]
	grid.AreaMap.SetText(text)
	grid.AreaField.SetDoneFunc(func(key tcell.Key) {
		grid.Grid.Clear()
		grid.Grid.AddItem(grid.InputGrid, 0, 0, 1, 1, 0, 0, false)
		grid.App.SetFocus(grid.InputField)
	})
}

func Prelude() {
	var text = setup.TextElements["prelude"][setup.Language]
	grid.Grid.Clear()
	grid.Grid.AddItem(grid.AreaGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AreaMap.SetText(text)
	grid.AreaField.SetText("")
	grid.App.SetFocus(grid.AreaField)
	grid.AreaField.SetDoneFunc(func(key tcell.Key) {
		//Intro()
		if _, err := os.Stat(setup.PathName + "/save/fangotasia.lang"); os.IsNotExist(err) {
			SetupLanguage()
		} else {
			lang, _ := ioutil.ReadFile(setup.PathName + "/save/fangotasia.lang")
			setup.Language = string(lang)
			Intro()
		}
	})
}

func SetupLanguage() {
	var keys []string
	keys = make([]string, 0)

	filename := setup.PathName + "/save/fangotasia.lang"

	for k := range setup.TextElements["selectLanguage"] {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	grid.Grid.Clear()
	grid.Grid.AddItem(grid.LanguageGrid, 0, 0, 1, 1, 0, 0, false)
	grid.LanguageSelect.
		SetLabel(setup.TextElements["selectLanguage"][setup.Language]).
		SetOptions(keys, nil).
		SetSelectedFunc(func(o string, i int) {
			//i, o := grid.LanguageSelect.GetCurrentOption()
			if i >= 0 {
				ioutil.WriteFile(filename, []byte(o), 0777)
			}
			setup.Language = o
			grid.Grid.Clear()
			grid.Grid.AddItem(grid.AreaGrid, 0, 0, 1, 1, 0, 0, false)
			grid.App.SetFocus(grid.AreaField)
			//grid.InputField.SetLabel(fmt.Sprintf("%s? > ", setup.TextElements["andNow"][setup.Language]))
			//grid.Surroundings.SetText(strings.Join(view.Surroundings(area), "\n"))
			area := setup.GetAreaByID(1)
			grid.Surroundings.SetText(strings.Join(view.Surroundings(area), "\n"))
			Intro()
			/*
				grid.Grid.Clear()
				grid.Grid.AddItem(grid.InputGrid, 0, 0, 1, 1, 0, 0, false)
				grid.App.SetFocus(grid.InputField)
				grid.Response.SetText("")
			*/
		})
	grid.App.SetFocus(grid.LanguageSelect)
}
