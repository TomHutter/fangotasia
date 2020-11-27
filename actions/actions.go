package actions

import (
	"fantasia/movement"
	"fantasia/setup"
	"fantasia/view"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
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
		answer := setup.Reactions["unknownVerb"]
		notice := fmt.Sprintf(answer.Statement, command)
		view.AddFlashNotice(notice, answer.Sleep, setup.RED)
		return area
	}

	if knownVerb.Single {
		call := reflect.ValueOf(&order).MethodByName(knownVerb.Func)
		if !call.IsValid() {
			view.AddFlashNotice(fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, setup.RED)
			return area
		}
		val := call.Call([]reflect.Value{})
		var notice []string
		for i := 0; i < val[0].Len(); i++ {
			notice = append(notice, val[0].Index(i).String())
		}
		// ToDo: get rid of knownVerb.Sleep
		view.AddFlashNotice(strings.Join(notice, "\n"), knownVerb.Sleep, setup.GREEN)
		return area
	}

	argv := make([]reflect.Value, 0)

	switch knownVerb.Name {
	case "n", "s", "o", "w":
		obj = Object{}
		argv = append(argv, reflect.ValueOf(area))
		argv = append(argv, reflect.ValueOf(knownVerb.Name))
	case "load":
		obj = Object{}
		argv = append(argv, reflect.ValueOf(area))
	case "save":
		obj = Object{}
		argv = append(argv, reflect.ValueOf(area))
	default:
		if len(parts) < 1 {
			answer := setup.Reactions["needObject"]
			notice := fmt.Sprintln(answer.Statement)
			view.AddFlashNotice(notice, answer.Sleep, setup.RED)
			return area
		}
		for _, p := range parts {
			obj = Object(setup.GetObjectByName(p))
			if obj != (Object{}) {
				argv = append(argv, reflect.ValueOf(area))
				break
			}
		}
	}

	if len(argv) < 1 {
		answer := setup.Reactions["unknownNoun"]
		notice := fmt.Sprintf(answer.Statement, strings.Join(parts, " "))
		view.AddFlashNotice(notice, answer.Sleep, answer.Color)
		return area
	}

	// now method and all args should be known
	call := reflect.ValueOf(obj).MethodByName(knownVerb.Func)
	if !call.IsValid() {
		view.AddFlashNotice(fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, setup.RED)
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
		color = setup.GREEN
	default:
		color = setup.RED
	}
	switch knownVerb.Func {
	case "Move", "Climb", "Load":
		color = setup.GREEN
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
		view.FlashNotice()
	}
	// KO
	if val[0].Field(2).Bool() == true {
		GameOver()
	}
	return area
}

func (object Object) inArea(area setup.Area) bool {
	return object.Properties.Area == area.ID
}

func (object Object) inInventory() bool {
	return object.Properties.Area == 1000
}

func (object Object) inUse() bool {
	return object.Properties.Area == 2000
}

func (object Object) available(area setup.Area) bool {
	return object.inArea(area) || object.inInventory() || object.inUse()
}

func (object Object) snatchFrom(opponent Object) (r setup.Reaction) {
	hood := Object(setup.GetObjectByID(13))
	area := setup.GetAreaByID(object.Properties.Area)
	if opponent.inArea(area) && !hood.inUse() {
		r = setup.Reactions["wontLet"]
		r.Statement = fmt.Sprintf("%s %s %s",
			strings.Title(opponent.Properties.Description.Article),
			opponent.Properties.Description.Short,
			r.Statement)
		return
	}
	return object.pick()
}

func (object Object) Open(area int) (ok bool, answer []string) {
	answer = append(answer, "lässt sich öffnen")
	return true, answer
}

func (object Object) Take(area setup.Area) (r setup.Reaction) {
	if !object.available(area) {
		r = setup.Reactions["dontSee"]
		return
	}
	if object.inInventory() || object.inUse() {
		r = setup.Reactions["haveAlready"]
		return
	}

	switch object.ID {
	case 10, 16, 18, 21, 22, 27, 36, 40, 42:
		r = setup.Reactions["silly"]
		return
	case 29, 14:
		r = setup.Reactions["tooHeavy"]
		return
	case 34:
		if !Object(setup.GetObjectByID(9)).inUse() {
			r = setup.Reactions["tooHeavy"]
			return
		}
	case 17:
		return object.snatchFrom(Object(setup.GetObjectByID(16)))
	case 19:
		return object.snatchFrom(Object(setup.GetObjectByID(18)))
	case 35:
		return object.snatchFrom(Object(setup.GetObjectByID(36)))
	case 44:
		return object.snatchFrom(Object(setup.GetObjectByID(42)))
	case 32, 43:
		r = setup.Reactions["unreachable"]
		return
	}
	return object.snatchFrom(Object(setup.GetObjectByID(10)))
}

func (object Object) Stab(area setup.Area) (r setup.Reaction) {
	if !Object(setup.GetObjectByID(15)).inInventory() &&
		!Object(setup.GetObjectByID(25)).inInventory() &&
		!Object(setup.GetObjectByID(33)).inInventory() {
		r = setup.Reactions["noTool"]
		return
	}
	switch object.ID {
	case 14:
		r = setup.Reactions["tryCut"]
		return
	case 10:
		r = setup.Reactions["stabGrub"]
		return
	case 16:
		r = setup.Reactions["stabBaer"]
		return
	case 18:
		if Object(setup.GetObjectByID(13)).inUse() {
			dwarf := setup.GetObjectByID(18)
			dwarf.Properties.Area = -1
			setup.GameObjects[dwarf.ID] = dwarf.Properties
			r = setup.Reactions["stabDwarfHooded"]
		} else {
			r = setup.Reactions["stabDwarf"]
		}
		return
	case 36:
		if Object(setup.GetObjectByID(13)).inUse() {
			gnome := setup.GetObjectByID(36)
			gnome.Properties.Area = -1
			setup.GameObjects[gnome.ID] = gnome.Properties
			r = setup.Reactions["stabGnomeHooded"]
		} else {
			r = setup.Reactions["stabGnome"]
		}
		return
	case 42:
		r = setup.Reactions["stabDragon"]
		return
	}
	r = setup.Reactions["dontKnowHow"]
	return
}

// As Move is called in context of object handling, Move reflects on Object even obj is not used.
func (obj Object) Move(area setup.Area, dir string) (r setup.Reaction, areaID int) {
	//func Move(area int, direction int, text []string) int {
	//if direction == 0 {
	//	return 0, "Ich brauche eine Richtung."
	//}
	// wearing the hood?
	hood := setup.GetObjectByID(13)
	if Object(hood).inUse() {
		hood.Properties.Area = -1
		setup.GameObjects[hood.ID] = hood.Properties
		r = setup.Reactions["hoodInUse"]
	}

	var direction = map[string]int{"n": 0, "s": 1, "o": 2, "w": 3}

	newArea := area.Properties.Directions[direction[dir]]
	if newArea == 0 {
		r = setup.Reactions["noWay"]
		areaID = area.ID
		return
		//view.Flash(text, "In diese Richtung führt kein Weg.")
		//return area
	}

	// barefoot on unknown terrain?
	if !Object(setup.GetObjectByID(31)).inUse() {
		r = setup.Reactions["noShoes"]
		return
	}

	// Area 30 and 25 are connected by a door. Is it open?
	if (area.ID == 30 || area.ID == 25 && direction[dir] == 0) && !setup.DoorOpen {
		r = setup.Reactions["locked"]
		areaID = area.ID
		return
		//view.Flash(text, "Die Tür ist versperrt.")
		//return area
	}
	movement.RevealArea(newArea)
	moves += 1
	r.OK = true
	r.KO = false
	areaID = newArea
	// Direction Moor?
	if newArea == 5 {
		for _, o := range setup.ObjectsInArea(setup.GetAreaByID(1000)) {
			o.Properties.Area = 29
			setup.GameObjects[o.ID] = o.Properties
		}
		r = setup.Reactions["inTheMoor"]
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

func (obj Object) Use(area setup.Area) (r setup.Reaction) {
	r.KO = false
	r.Sleep = 2
	switch obj.ID {

	case 13:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
		} else {
			obj.Properties.Area = 2000
			setup.GameObjects[obj.ID] = obj.Properties
			r = setup.Reactions["useHood"]
		}
	case 31:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
		} else {
			obj.Properties.Area = 2000
			setup.GameObjects[obj.ID] = obj.Properties
			r = setup.Reactions["useShoes"]
		}
	default:
		r = setup.Reactions["dontKnowHow"]
	}
	return
}

func (object Object) Climb(area setup.Area) (r setup.Reaction, areaID int) {
	if area.ID == 31 {
		moves += 1
		r.OK = true
		areaID = 9
		return
	}
	if area.ID == 9 && object.ID == 27 {
		movement.RevealArea(31)
		moves += 1
		tree := setup.GetAreaByID(31)
		tree.Properties.Visited = true
		setup.GameAreas[tree.ID] = tree.Properties
		r.OK = true
		areaID = 31
		return
	}
	if !object.available(area) {
		r = setup.Reactions["dontSee"]
		areaID = area.ID
		return
	}
	if object.ID != 27 {
		r = setup.Reactions["silly"]
		areaID = area.ID
		return
	}
	return
}

func (v *verb) Verbs() (verbs []string) {
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
	return
}

func (c *verb) Inventory() (inv []string) {
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
	return
}

func GameOver() {
	var board []string
	sum := 0
	board = append(board, fmt.Sprintln("G A M E    O V E R"))
	//board = append(board, fmt.Sprint(""))
	inv := setup.ObjectsInArea(setup.GetAreaByID(1000))
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
	board = append(board, fmt.Sprint("Noch ein Spiel (j/n)?"))
	view.PrintScreen(board)
	res := view.Scanner("once: true")
	if strings.ToLower(res) == "j" {
		setup.Init()
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

func (obj Object) pick() (r setup.Reaction) {
	inv := setup.GetAreaByID(1000)
	inventory := setup.ObjectsInArea(inv)
	if len(inventory) > 6 {
		r = setup.Reactions["invFull"]
		return
	}

	obj.Properties.Area = 1000
	setup.GameObjects[obj.ID] = obj.Properties
	r = setup.Reactions["ok"]
	return
}

func (obj Object) drop(area setup.Area) (r setup.Reaction) {
	if obj.Properties.Area != 1000 {
		r = setup.Reactions["dontHave"]
		return
	}
	obj.Properties.Area = area.ID
	setup.GameObjects[obj.ID] = obj.Properties
	r = setup.Reactions["ok"]
	return
}

func folderListing() (filename string) {
	_, caller, _, _ := runtime.Caller(0)
	pathname := path.Dir(caller) + "/../save/"
	files, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Files:")

	for _, f := range files {
		fmt.Println(f.Name())
	}
	//filename = pathname + view.Scanner("prompt: filename > ")
	filename = pathname + "mrgl"
	return
}

func (obj Object) Save(area setup.Area) (ok bool, err error) {
	m := make(map[interface{}]interface{})
	m["area"] = area.ID
	m["map"] = setup.Map
	m["objects"] = setup.GameObjects

	filename := folderListing()

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	defer file.Close()

	if err = yaml.NewEncoder(file).Encode(m); err != nil {
		fmt.Println(err)
		return false, err
	}

	fmt.Printf("File %s written successfully\n", filename)
	/*
		err = file.Close()
		if err != nil {
			fmt.Println(err)
			return false, err
		}
	*/
	return true, nil
}

func (obj Object) Load(area setup.Area) (r setup.Reaction, areaID int) {
	var content struct {
		AreaID  int                            `yaml:"area"`
		AreaMap [12][10]int                    `yaml:"map"`
		Objects map[int]setup.ObjectProperties `yaml:"objects"`
	}

	filename := folderListing()

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		r.Statement = err.Error()
		r.OK = false
		return
	}
	err = yaml.Unmarshal(yamlFile, &content)
	if err != nil {
		r.Statement = err.Error()
		r.OK = false
		return
	}

	setup.GameObjects = content.Objects
	setup.Map = content.AreaMap

	r = setup.Reactions["loaded"]
	areaID = content.AreaID
	return
}
