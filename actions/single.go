package actions

import (
	"fantasia/setup"
	"fmt"
	"strings"
)

func (v *verb) Verbs() (verbs []string, sleep int) {
	//func (c *command) Verben() {
	verbs = append(verbs, "Verben, die ich kenne: ")
	var line []string
	for i, val := range setup.Verbs {
		if i > 1 && i%10 == 0 {
			verbs = append(verbs, strings.Join(line, ", "))
			line = make([]string, 0)
		}
		line = append(line, string(val.Name))
	}
	verbs = append(verbs, strings.Join(line, ", "))
	sleep = -1
	return
}

func (v *verb) Inventory() (inv []string, sleep int) {
	objects := setup.ObjectsInArea(setup.GetAreaByID(1000))
	if len(objects) == 0 {
		//fmt.Println("Ich habe nichts dabei.")
		inv = append(inv, setup.Reactions["invEmpty"].Statement)
		return
	}

	inv = append(inv, setup.Reactions["inv"].Statement)
	for _, o := range objects {
		inv = append(inv, fmt.Sprintf("- %s", o.Properties.Description.Long))
	}
	sleep = -1
	return
}
