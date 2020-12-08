package actions

import (
	"fangotasia/setup"
	"fangotasia/view"
	"fmt"
	"strings"
)

func (v *verb) Verbs() (verbs []string) {
	verbs = append(verbs, fmt.Sprintf("%s: ", setup.TextElements["verbsKnown"]))
	var line []string
	for i, val := range setup.Verbs {
		if i > 1 && i%10 == 0 {
			verbs = append(verbs, strings.Join(line, ", "))
			line = make([]string, 0)
		}
		line = append(line, string(val.Name))
	}
	verbs = append(verbs, strings.Join(line, ", "))
	return
}

func (v *verb) Inventory() (inv []string) {
	objects := setup.ObjectsInArea(setup.GetAreaByID(setup.INVENTORY))
	if len(objects) == 0 {
		inv = append(inv, setup.Reactions["invEmpty"].Statement[0])
		return
	}

	inv = append(inv, setup.Reactions["inv"].Statement[0])
	for _, o := range objects {
		obj := view.Highlight(o.Properties.Description.Long, "[green:black:-]")
		inv = append(inv, fmt.Sprintf("- %s", obj))
	}
	return
}

func (v *verb) End() (r []string) {
	scoreBoard(true, false)
	return
}
