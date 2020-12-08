package actions

import (
	"fangotasia/grid"
	"fangotasia/setup"
	"fangotasia/view"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"
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

func Parse(input string, area setup.Area) setup.Area {

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
			answer := setup.GetReactionByName("unknownVerb")
			notice := fmt.Sprintf(answer.Statement[0], command)
			grid.InputField.SetText("")
			grid.Response.SetText(
				fmt.Sprintf("\n%s%s%s\n",
					"[red]",
					notice, "[-:black:-]"))
		}
		return area
	}

	if knownVerb.Single {
		call := reflect.ValueOf(&order).MethodByName(knownVerb.Func)
		if !call.IsValid() {
			grid.InputField.SetText("")
			grid.Response.SetText(
				fmt.Sprintf(fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, "[red]"))
			return area
		}
		val := call.Call([]reflect.Value{})
		var notice []string
		for i := 0; i < val[0].Len(); i++ {
			notice = append(notice, val[0].Index(i).String())
		}
		grid.InputField.SetText("")
		grid.Response.SetText(
			fmt.Sprintf("\n%s%s%s\n",
				"[green:black:-]",
				strings.Join(notice, "\n"),
				"[-:black:-]"))
		return area
	}

	argv := make([]reflect.Value, 0)

	switch knownVerb.Func {
	case "Move":
		obj = Object{}
		argv = append(argv, reflect.ValueOf(area))
		argv = append(argv, reflect.ValueOf(knownVerb.Name))
	case "Load", "Save", "Jump", "Map", "Help":
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
			answer := setup.GetReactionByName("needObject")
			notice := fmt.Sprintln(answer.Statement[0])
			grid.InputField.SetText("")
			grid.Response.SetText(
				fmt.Sprintf("\n%s%s%s\n",
					"[red]",
					notice,
					"[-:black:-]"))
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
		grid.InputField.SetText("")
		grid.Response.SetText(
			fmt.Sprintf(fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, "[red]"))
		return area
	}
	val := call.Call(argv)

	var color string
	var notice string
	// Reaction
	resp := val[0].Field(0)
	respLen := resp.Len()
	if respLen > 0 {
		rand.Seed(time.Now().UnixNano())
		notice = resp.Index(rand.Intn(respLen)).String()
		for notice == grid.Response.GetText(true) {
			// get random choice of possible reactions
			notice = resp.Index(rand.Intn(respLen)).String()
		}
	}
	switch val[0].Field(3).String() {
	case "GREEN":
		color = "[green]"
	default:
		color = "[red]"
	}
	switch knownVerb.Func {
	case "Move", "Climb", "Load", "Jump":
		// OK
		if val[0].Field(1).Bool() == true {
			// Area
			area = setup.GetAreaByID(int(val[1].Int()))
			area.Properties.Visited = true
			setup.GameAreas[area.ID] = area.Properties
		}
		fallthrough
	default:
		grid.InputField.SetText("")
		grid.Response.SetText(
			fmt.Sprintf("\n%s%s%s\n",
				color,
				notice, "[-:black:-]"))
	}
	// KO ?
	if val[0].Field(2).Bool() == true {
		time.Sleep(time.Duration(6) * time.Second)
		area = setup.GetAreaByID(1)
		grid.InputField.SetText("")
		grid.Response.SetText("")
		grid.Surroundings.SetText("")
		GameOver(true)
	}
	return area
}

func REPL(area setup.Area) {
	for {
		command := <-grid.Input
		area = Parse(command, area)
		grid.Surroundings.SetText(strings.Join(view.Surroundings(area), "\n"))
	}
}
