package actions

import (
	"fantasia/config"
	"fantasia/movement"
	"fantasia/view"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

type verb config.Verb
type Object config.Object
type reaction struct {
	OK     bool
	KO     bool
	Answer []string
	Sleep  int
	AreaID int
}

var object string
var moves int
var visited [51]bool

//type verb []string

//type command struct{}

func Parse(input string, area config.Area, text []string) config.Area {

	var command string
	var order verb
	var knownVerb config.Verb
	var obj Object

	re := regexp.MustCompile("[\\s,\\t]+")

	parts := re.Split(input, -1)
	command = strings.ToLower(parts[0])
	// pop parts
	parts = parts[1:]

	//fmt.Println(parts) // ["Have", "a", "great", "day!"]
	//parts := strings.r  Split(input, )
	for _, v := range config.Verbs {
		if v.Name == command {
			knownVerb = v
			break
		}
	}

	if knownVerb == (config.Verb{}) {
		notice := fmt.Sprintf(config.Answers["unknownVerb"], command)
		view.AddFlashNotice(notice, 4, config.RED)
		return area
	}

	//fmt.Printf("Valid verb '%s' found.\n", knownVerb.Name)
	if knownVerb.Single {
		call := reflect.ValueOf(&order).MethodByName(knownVerb.Func)
		//view.Flash(text, fmt.Sprintf("%s", call.IsNil()), 2, config.RED)
		//fmt.Println(call.String())
		//fmt.Println(call.IsValid())
		if !call.IsValid() {
			view.AddFlashNotice(fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, config.RED)
			return area
		}
		//val := reflect.ValueOf(&order).MethodByName(knownVerb.Func).Call(argv)
		val := call.Call([]reflect.Value{})
		/*
			// call reflec.Call for verb without arguments
			call := reflect.ValueOf(&order).MethodByName(knownVerb.Func)
			fmt.Println(call)
			val := reflect.ValueOf(&order).MethodByName(knownVerb.Func).Call([]reflect.Value{})
		*/
		var notice []string
		for i := 0; i < val[0].Len(); i++ {
			notice = append(notice, val[0].Index(i).String())
		}
		view.AddFlashNotice(strings.Join(notice, "\n"), knownVerb.Sleep, config.GREEN)
		//for i := 0; i < val[0].Len(); i++ {
		//	fmt.Println(val[0].Index(i).String())
		//}
		return area
	}

	argv := make([]reflect.Value, 0)

	switch knownVerb.Name {
	case "n", "s", "o", "w":
		//obj = Object(config.GetObjectByName(knownVerb.Func))
		obj = Object{}
		argv = append(argv, reflect.ValueOf(area))
		argv = append(argv, reflect.ValueOf(knownVerb.Name))
	default:
		if len(parts) < 1 {
			notice := fmt.Sprintln(config.Answers["needObject"])
			view.AddFlashNotice(notice, knownVerb.Sleep, config.RED)
			return area
		}
		// call reflec.Call for verb with arguments
		// - first valid noun
		// - area
		//argv := make([]reflect.Value, 1)
		// loop over all input parts and see if we find a valid object
		for _, p := range parts {
			obj = Object(config.GetObjectByName(p))
			if obj != (Object{}) {
				//fmt.Printf("Valid object '%s' found.\n", obj.Properties.Description.Short)
				//argv = append(argv, reflect.ValueOf(o))
				argv = append(argv, reflect.ValueOf(area))
				break

				/*
					//val := reflect.ValueOf(&order).MethodByName(v.Func).Call(argv)
					//val := reflect.ValueOf(&order).MethodByName(v.Func).Call([]reflect.Value{reflect.ValueOf(order)})
					val := reflect.ValueOf(&order).MethodByName(knownVerb.Func).Call(argv)
					for i := 0; i < val[0].Len(); i++ {
						fmt.Println(val[0].Index(i).String())
					}
					fmt.Println(val[1])
					return
				*/
			}
		}
	}

	if len(argv) < 1 {
		notice := fmt.Sprintf(config.Answers["unknownNoun"], strings.Join(parts, " "))
		view.AddFlashNotice(notice, knownVerb.Sleep, config.RED)
		return area
	}

	// now method and all args should be known
	call := reflect.ValueOf(obj).MethodByName(knownVerb.Func)
	//view.Flash(text, fmt.Sprintf("%s", call.IsNil()), 2, config.RED)
	//fmt.Println(call.String())
	//fmt.Println(call.IsValid())
	if !call.IsValid() {
		view.AddFlashNotice(fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, config.RED)
		return area
	}
	//val := reflect.ValueOf(&order).MethodByName(knownVerb.Func).Call(argv)
	val := call.Call(argv)

	var notice []string
	var color string
	// Answer
	for i := 0; i < val[0].Field(2).Len(); i++ {
		notice = append(notice, val[0].Field(2).Index(i).String())
	}
	// Sleep
	sleep := int(val[0].Field(3).Int())
	switch knownVerb.Func {
	case "Move", "Climb":
		color = config.RED
		// OK
		if val[0].Field(0).Bool() == true {
			// Area
			area = config.GetAreaByID(int(val[0].Field(4).Int()))
		}
		// add notice. Give feedback in the next screen.
		view.AddFlashNotice(strings.Join(notice, "\n"), int(sleep), color)
	default:
		// OK
		if val[0].Field(0).Bool() == true {
			color = config.GREEN
		} else {
			color = config.RED
		}
		// add notice and give feedback in this screen.
		view.AddFlashNotice(strings.Join(notice, "\n"), int(sleep), color)
		view.FlashNotice()
	}
	// KO
	if val[0].Field(1).Bool() == true {
		GameOver()
	}
	return area
}

func (object Object) inArea(area config.Area) bool {
	return object.Properties.Area == area.ID
}

func (object Object) inInventory() bool {
	return object.Properties.Area == 1000
}

func (object Object) inUse() bool {
	return object.Properties.Area == 2000
}

func (object Object) available(area config.Area) bool {
	return object.inArea(area) || object.inInventory() || object.inUse()
}

func (object Object) snatchFrom(opponent Object) (r reaction) {
	hood := Object(config.GetObjectByID(13))
	area := config.GetAreaByID(object.Properties.Area)
	if opponent.inArea(area) && !hood.inUse() {
		r.Answer = append(r.Answer, fmt.Sprintf("%s %s %s",
			strings.Title(opponent.Properties.Description.Article),
			opponent.Properties.Description.Short,
			config.Answers["wontLet"]))
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}
	return object.pick()
}

func (object Object) Open(area int) (ok bool, answer []string) {
	answer = append(answer, "lässt sich öffnen")
	return true, answer
}

func (object Object) Take(area config.Area) (r reaction) {
	if !object.available(area) {
		r.Answer = append(r.Answer, config.Answers["dontSee"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}
	if object.inInventory() || object.inUse() {
		r.Answer = append(r.Answer, config.Answers["haveAlready"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}

	switch object.ID {
	case 10, 16, 18, 21, 22, 27, 36, 40, 42:
		r.Answer = append(r.Answer, config.Answers["silly"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	case 29, 14:
		r.Answer = append(r.Answer, config.Answers["tooHeavy"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	case 34:
		if !Object(config.GetObjectByID(9)).inUse() {
			r.Answer = append(r.Answer, config.Answers["tooHeavy"])
			r.OK = false
			r.KO = false
			r.Sleep = 2
			return
		}
	case 17:
		return object.snatchFrom(Object(config.GetObjectByID(16)))
	case 19:
		return object.snatchFrom(Object(config.GetObjectByID(18)))
	case 35:
		return object.snatchFrom(Object(config.GetObjectByID(36)))
	case 44:
		return object.snatchFrom(Object(config.GetObjectByID(42)))
	case 32, 43:
		r.Answer = append(r.Answer, config.Answers["unreachable"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}
	return object.snatchFrom(Object(config.GetObjectByID(10)))
}

func (object Object) Stab(area config.Area) (r reaction) {
	if !Object(config.GetObjectByID(15)).inInventory() &&
		!Object(config.GetObjectByID(25)).inInventory() &&
		!Object(config.GetObjectByID(33)).inInventory() {
		r.Answer = append(r.Answer, config.Answers["noTool"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}
	switch object.ID {
	case 14:
		r.Answer = append(r.Answer, config.Answers["tryCut"])
		r.OK = true
		r.KO = false
		r.Sleep = 2
		return
	case 10:
		r.Answer = append(r.Answer, config.Answers["stabGrub"])
		r.OK = true
		r.KO = false
		r.Sleep = 2
		return
	case 16:
		r.Answer = append(r.Answer, config.Answers["starbBaer"])
		r.OK = true
		r.KO = false
		r.Sleep = 2
		return
	case 18:
		if Object(config.GetObjectByID(13)).inUse() {
			dwarf := config.GetObjectByID(18)
			dwarf.Properties.Area = -1
			config.GameObjects[dwarf.ID] = dwarf.Properties
			r.Answer = append(r.Answer, config.Answers["stabDwarfHooded"])
			r.OK = true
			r.KO = false
			r.Sleep = 1
		} else {
			r.Answer = append(r.Answer, config.Answers["stabDwarf"])
			r.OK = false
			r.KO = true
			r.Sleep = 2
		}
		return
	case 36:
		if Object(config.GetObjectByID(13)).inUse() {
			gnome := config.GetObjectByID(36)
			gnome.Properties.Area = -1
			config.GameObjects[gnome.ID] = gnome.Properties
			r.Answer = append(r.Answer, config.Answers["stabGnomeHooded"])
			r.OK = true
			r.KO = false
			r.Sleep = 1
		} else {
			r.Answer = append(r.Answer, config.Answers["stabGnome"])
			r.OK = false
			r.KO = true
			r.Sleep = 2
		}
		return
	case 42:
		r.Answer = append(r.Answer, config.Answers["stabDragon"])
		r.OK = false
		r.KO = true
		r.Sleep = 2
		return
	}
	r.Answer = append(r.Answer, config.Answers["dontKnow"])
	r.OK = true
	r.KO = false
	r.Sleep = 2
	return
}

// As Move is called in context of object handling, Move reflects on Object even obj is not used.
func (obj Object) Move(area config.Area, dir string) (r reaction) {
	//func Move(area int, direction int, text []string) int {
	//if direction == 0 {
	//	return 0, "Ich brauche eine Richtung."
	//}
	// barefoot on unknown terrain?
	if !Object(config.GetObjectByID(31)).inUse() {
		r.Answer = append(r.Answer, config.Answers["noShoes"])
		r.OK = false
		r.KO = true
		r.Sleep = 3
		return
	}

	// wearing the hood?
	hood := config.GetObjectByID(13)
	if Object(hood).inUse() {
		hood.Properties.Area = -1
		config.GameObjects[hood.ID] = hood.Properties
		r.Answer = append(r.Answer, config.Answers["hoodInUse"])
		r.Sleep = 2
	}

	var direction = map[string]int{"n": 0, "s": 1, "o": 2, "w": 3}

	newArea := area.Properties.Directions[direction[dir]]
	if newArea == 0 {
		r.Answer = append(r.Answer, config.Answers["noWay"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		r.AreaID = area.ID
		return
		//view.Flash(text, "In diese Richtung führt kein Weg.")
		//return area
	}
	// Area 30 and 25 are connected by a door. Is it open?
	if (area.ID == 30 || area.ID == 25 && direction[dir] == 0) && !config.DoorOpen {
		r.Answer = append(r.Answer, config.Answers["locked"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		r.AreaID = area.ID
		return
		//view.Flash(text, "Die Tür ist versperrt.")
		//return area
	}
	movement.RevealArea(newArea)
	moves += 1
	visited[area.ID] = true
	r.OK = true
	r.KO = false
	r.AreaID = newArea
	// Direction Moor?
	if newArea == 5 {
		for _, o := range config.ObjectsInArea(config.GetAreaByID(1000)) {
			o.Properties.Area = 29
			config.GameObjects[o.ID] = o.Properties
		}
		r.Answer = append(r.Answer, config.Answers["inTheMoor"])
		r.Sleep = 6
		r.OK = true
	}
	return
}

func useDoor() {
	/*
		495 rem ** sperre *********************
		496 f=0:gosub605:iffl=1thenfl=0:goto280
		497 ifno<>40andno<>35thenprinta$(2):goto280
		498 ifno=35thenprint"versuche 'oeffne'.":goto280
		499 iftu=1thenprint"ist schon offen !":goto280
		500 ifge(26)<>-1thenprint"ich habe keinen schluessel.":goto280
		501 print"gut.":tu=1:goto281
		502 :
	*/
}

func (obj Object) Use(area config.Area) (r reaction) {
	r.KO = false
	r.Sleep = 2
	switch obj.ID {

	case 13:
		if !obj.inInventory() {
			r.Answer = append(r.Answer, config.Answers["dontHave"])
			r.OK = false
		} else {
			obj.Properties.Area = 2000
			config.GameObjects[obj.ID] = obj.Properties
			r.Answer = append(r.Answer, config.Answers["hood"])
			r.OK = true
		}
	case 31:
		if !obj.inInventory() {
			r.Answer = append(r.Answer, config.Answers["dontHave"])
			r.OK = false
		} else {
			obj.Properties.Area = 2000
			config.GameObjects[obj.ID] = obj.Properties
			r.Answer = append(r.Answer, config.Answers["shoes"])
			r.OK = true
		}
	default:
		r.Answer = append(r.Answer, config.Answers["dontKnow"])
		r.OK = false
	}
	return
}

func (object Object) Climb(area config.Area) (r reaction) {
	if area.ID == 31 {
		moves += 1
		r.OK = true
		r.AreaID = 9
		return
	}
	if area.ID == 9 && object.ID == 27 {
		movement.RevealArea(31)
		moves += 1
		visited[31] = true
		r.OK = true
		r.AreaID = 31
		return
	}
	if !object.available(area) {
		r.Answer = append(r.Answer, config.Answers["dontSee"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		r.AreaID = area.ID
		return
	}
	if object.ID != 27 {
		r.Answer = append(r.Answer, config.Answers["silly"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		r.AreaID = area.ID
		return
	}
	return
}

func (v *verb) Verbs() (verbs []string) {
	//func (c *command) Verben() {
	verbs = append(verbs, "Verben, die ich kenne: ")
	var line []string
	for i, val := range config.Verbs {
		if i > 1 && i%10 == 0 {
			verbs = append(verbs, strings.Join(line, ", "))
			line = make([]string, 0)
		}
		line = append(line, string(val.Name))
	}
	verbs = append(verbs, strings.Join(line, ", "))
	return
}

func (c *verb) Inventory() (inv []string) {
	objects := config.ObjectsInArea(config.GetAreaByID(1000))
	if len(objects) == 0 {
		//fmt.Println("Ich habe nichts dabei.")
		inv = append(inv, config.Answers["invEmpty"])
		return
	}

	inv = append(inv, config.Answers["inv"])
	for _, o := range objects {
		inv = append(inv, fmt.Sprintf("- %s", o.Properties.Description.Long))
	}
	return
}

func GameOver() {
	var board []string
	sum := 0
	board = append(board, fmt.Sprintln("G A M E    O V E R"))
	//board = append(board, fmt.Sprint(""))
	inv := config.ObjectsInArea(config.GetAreaByID(1000))
	if len(inv) > 0 {
		board = append(board, fmt.Sprintln("Du besitzt:"))
		//board = append(board, fmt.Sprint(""))
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
	if visited[5] {
		board = append(board, fmt.Sprintf("- Du bist im Moor gewesen: %d Punkte", 2))
		sum += 2
	}
	if visited[29] {
		board = append(board, fmt.Sprintf("- Du hast die verlassene Burg besucht: %d Punkte", 3))
		sum += 3
	}
	if visited[31] {
		board = append(board, fmt.Sprintf("- Du bist auf einen Baum geklettert: %d Punkte", 4))
		sum += 4
	}
	board = append(board, fmt.Sprint(""))
	board = append(board, fmt.Sprintf("Du hast %d von 200 Punkten!", sum))
	board = append(board, fmt.Sprint("Noch ein Spiel (j/n)?"))
	view.PrintScreen(board)
	res := view.Scanner("once: true")
	if strings.ToLower(res) == "j" {
		config.Init()
	}
	/*
		621 poke214,9:poke211,13:sysvd:printb$"-rang ";
		622 ifpu=0thenprint"10 -":goto632
		623 ifpu<25thenprint"9 -":goto632
		624 ifpu<50thenprint"8 -":goto632
		625 ifpu<75thenprint"7 -":goto632
		626 ifpu<100thenprint"6 -":goto632
		627 ifpu<125thenprint"5 -":goto632
		628 ifpu<150thenprint"4 -":goto632
		629 ifpu<175thenprint"3 -":goto632
		630 ifpu<200thenprint"2 -"goto632
		631 ifpu=200thenprint"1 -"
	*/
	//return true
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	os.Exit(0)
}

//func (obj object) hit(inv []string) {
//}

func (obj Object) pick() (r reaction) {
	inv := config.GetAreaByID(1000)
	inventory := config.ObjectsInArea(inv)
	if len(inventory) > 6 {
		r.Answer = append(r.Answer, config.Answers["invFull"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}

	obj.Properties.Area = 1000
	config.GameObjects[obj.ID] = obj.Properties
	r.Answer = append(r.Answer, config.Answers["ok"])
	r.OK = true
	r.KO = false
	r.Sleep = 1
	return
}

func (obj Object) drop(area config.Area) (r reaction) {
	if obj.Properties.Area != 1000 {
		r.Answer = append(r.Answer, config.Answers["dontHave"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}
	obj.Properties.Area = area.ID
	config.GameObjects[obj.ID] = obj.Properties
	r.Answer = append(r.Answer, config.Answers["ok"])
	r.OK = true
	r.KO = false
	r.Sleep = 1
	return
}
