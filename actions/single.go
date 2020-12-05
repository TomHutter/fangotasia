package actions

import (
	"fantasia/intro"
	"fantasia/setup"
	"fantasia/view"
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
	objects := setup.ObjectsInArea(setup.GetAreaByID(setup.INVENTORY))
	if len(objects) == 0 {
		//fmt.Println("Ich habe nichts dabei.")
		inv = append(inv, setup.Reactions["invEmpty"].Statement)
		sleep = 2
		return
	}

	inv = append(inv, setup.Reactions["inv"].Statement)
	for _, o := range objects {
		obj := view.Highlight(o.Properties.Description.Long, "[green:black:-]")
		inv = append(inv, fmt.Sprintf("- %s", obj))
	}
	sleep = -1
	return
}

func (v *verb) End() (r []string, sleep int) {
	r = append(r, "Yippeeee....")
	sleep = 3
	GameOver(false)
	return
}

func (v *verb) Help() (r []string, sleep int) {
	fmt.Print("\n\n")
	res := view.Scanner("once: true", "prompt: Ich kann nur die Anleitung wiederholen. (j/n)")
	if strings.ToLower(res) == "j" {
		intro.Intro()
	}
	return
}
