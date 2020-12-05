package actions

import (
	"fantasia/movement"
	"fantasia/setup"
	"fantasia/view"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type verb setup.Verb
type Object setup.Object
type reaction struct {
	OK       bool
	KO       bool
	Reaction []string
	Sleep    int
	AreaID   int
	Color    string
}

var object string
var moves int

func Parse(input string, area setup.Area, text []string) setup.Area {

	var command string
	var order verb
	var knownVerb setup.Verb
	var obj Object

	re := regexp.MustCompile("[\\s,\\t]+")

	parts := re.Split(input, -1)
	command = strings.ToLower(parts[0])
	// pop parts
	parts = parts[1:]

	for _, v := range setup.Verbs {
		if v.Name == command {
			knownVerb = v
			break
		}
	}

	if knownVerb == (setup.Verb{}) {
		// do we have input?
		if len(command) > 0 {
			answer := setup.Reactions["unknownVerb"]
			notice := fmt.Sprintf(answer.Statement, command)
			view.AddFlashNotice(notice, answer.Sleep, "[red]")
		}
		return area
	}

	if knownVerb.Single {
		call := reflect.ValueOf(&order).MethodByName(knownVerb.Func)
		if !call.IsValid() {
			view.AddFlashNotice(fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, "[red]")
			return area
		}
		val := call.Call([]reflect.Value{})
		var notice []string
		for i := 0; i < val[0].Len(); i++ {
			notice = append(notice, val[0].Index(i).String())
		}
		sleep := int(val[1].Int())
		// ToDo: get rid of knownVerb.Sleep
		view.AddFlashNotice(strings.Join(notice, "\n"), sleep, "[green]")
		return area
	}

	argv := make([]reflect.Value, 0)

	switch knownVerb.Func {
	case "Move":
		obj = Object{}
		argv = append(argv, reflect.ValueOf(area))
		argv = append(argv, reflect.ValueOf(knownVerb.Name))
	case "Load", "Save", "Jump":
		obj = Object{}
		argv = append(argv, reflect.ValueOf(area))
	case "Say":
		obj = Object{}
		argv = append(argv, reflect.ValueOf(area))
		if len(parts) > 0 {
			argv = append(argv, reflect.ValueOf(strings.Join(parts, " ")))
		} else {
			argv = append(argv, reflect.ValueOf(""))
		}
	case "Drink":
		obj = Object(setup.GetObjectByID(30))
		argv = append(argv, reflect.ValueOf(area))
	default:
		if len(parts) < 1 {
			answer := setup.Reactions["needObject"]
			notice := fmt.Sprintln(answer.Statement)
			view.AddFlashNotice(notice, answer.Sleep, "[red]")
			return area
		}
		for _, p := range parts {
			obj = Object(getObjectByName(p, area))
			if obj != (Object{}) {
				argv = append(argv, reflect.ValueOf(area))
				break
			}
		}
	}

	if len(argv) < 1 {
		return area
	}

	// now method and all args should be known
	call := reflect.ValueOf(obj).MethodByName(knownVerb.Func)
	if !call.IsValid() {
		view.AddFlashNotice(fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, "[red]")
		return area
	}
	val := call.Call(argv)

	var color string
	// Reaction
	notice := val[0].Field(0).String()
	// Sleep
	sleep := int(val[0].Field(4).Int())
	switch val[0].Field(3).String() {
	case "GREEN":
		color = "[green]"
	default:
		color = "[red]"
	}
	switch knownVerb.Func {
	case "Move", "Climb", "Load", "Jump":
		//color = setup.GREEN
		// OK
		if val[0].Field(1).Bool() == true {
			// Area
			area = setup.GetAreaByID(int(val[1].Int()))
			area.Properties.Visited = true
			setup.GameAreas[area.ID] = area.Properties
		}
		// add notice. Give feedback in the next screen.
		view.AddFlashNotice(notice, sleep, color)
	default:
		// OK
		/*
			if val[0].Field(1).Bool() == true {
				color = setup.GREEN
			} else {
				color = setup.RED
			}
		*/
		// add notice and give feedback in this screen.
		view.AddFlashNotice(notice, sleep, color)
		//view.FlashNotice()
	}
	// KO
	if val[0].Field(2).Bool() == true {
		view.FlashNotice()
		GameOver(true)
		area = setup.GetAreaByID(1)
		movement.RevealArea(area.ID)
		text := movement.DrawMap(area)
		surroundings := view.Surroundings(area)
		text = append(text, surroundings...)
		view.PrintScreen(text)
	}
	return area
}
