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
	OK     bool
	KO     bool
	Answer []string
	Sleep  int
	AreaID int
}

var object string
var moves int

//type verb []string

//type command struct{}

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

	//fmt.Println(parts) // ["Have", "a", "great", "day!"]
	//parts := strings.r  Split(input, )
	for _, v := range setup.Verbs {
		if v.Name == command {
			knownVerb = v
			break
		}
	}

	if knownVerb == (setup.Verb{}) {
		notice := fmt.Sprintf(setup.Answers["unknownVerb"], command)
		view.AddFlashNotice(notice, 4, setup.RED)
		return area
	}

	//fmt.Printf("Valid verb '%s' found.\n", knownVerb.Name)
	if knownVerb.Single {
		call := reflect.ValueOf(&order).MethodByName(knownVerb.Func)
		//view.Flash(text, fmt.Sprintf("%s", call.IsNil()), 2, setup.RED)
		//fmt.Println(call.String())
		//fmt.Println(call.IsValid())
		if !call.IsValid() {
			view.AddFlashNotice(fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, setup.RED)
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
		view.AddFlashNotice(strings.Join(notice, "\n"), knownVerb.Sleep, setup.GREEN)
		//for i := 0; i < val[0].Len(); i++ {
		//	fmt.Println(val[0].Index(i).String())
		//}
		return area
	}

	argv := make([]reflect.Value, 0)

	switch knownVerb.Name {
	case "n", "s", "o", "w":
		//obj = Object(setup.GetObjectByName(knownVerb.Func))
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
			notice := fmt.Sprintln(setup.Answers["needObject"])
			view.AddFlashNotice(notice, knownVerb.Sleep, setup.RED)
			return area
		}
		// call reflec.Call for verb with arguments
		// - first valid noun
		// - area
		//argv := make([]reflect.Value, 1)
		// loop over all input parts and see if we find a valid object
		for _, p := range parts {
			obj = Object(setup.GetObjectByName(p))
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
		notice := fmt.Sprintf(setup.Answers["unknownNoun"], strings.Join(parts, " "))
		view.AddFlashNotice(notice, knownVerb.Sleep, setup.RED)
		return area
	}

	// now method and all args should be known
	call := reflect.ValueOf(obj).MethodByName(knownVerb.Func)
	//view.Flash(text, fmt.Sprintf("%s", call.IsNil()), 2, setup.RED)
	//fmt.Println(call.String())
	//fmt.Println(call.IsValid())
	if !call.IsValid() {
		view.AddFlashNotice(fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, setup.RED)
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
	case "Move", "Climb", "Load":
		color = setup.RED
		// OK
		if val[0].Field(0).Bool() == true {
			// Area
			area = setup.GetAreaByID(int(val[0].Field(4).Int()))
			area.Properties.Visited = true
			setup.GameAreas[area.ID] = area.Properties
		}
		// add notice. Give feedback in the next screen.
		view.AddFlashNotice(strings.Join(notice, "\n"), int(sleep), color)
	default:
		// OK
		if val[0].Field(0).Bool() == true {
			color = setup.GREEN
		} else {
			color = setup.RED
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

func (object Object) snatchFrom(opponent Object) (r reaction) {
	hood := Object(setup.GetObjectByID(13))
	area := setup.GetAreaByID(object.Properties.Area)
	if opponent.inArea(area) && !hood.inUse() {
		r.Answer = append(r.Answer, fmt.Sprintf("%s %s %s",
			strings.Title(opponent.Properties.Description.Article),
			opponent.Properties.Description.Short,
			setup.Answers["wontLet"]))
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

func (object Object) Take(area setup.Area) (r reaction) {
	if !object.available(area) {
		r.Answer = append(r.Answer, setup.Answers["dontSee"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}
	if object.inInventory() || object.inUse() {
		r.Answer = append(r.Answer, setup.Answers["haveAlready"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}

	switch object.ID {
	case 10, 16, 18, 21, 22, 27, 36, 40, 42:
		r.Answer = append(r.Answer, setup.Answers["silly"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	case 29, 14:
		r.Answer = append(r.Answer, setup.Answers["tooHeavy"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	case 34:
		if !Object(setup.GetObjectByID(9)).inUse() {
			r.Answer = append(r.Answer, setup.Answers["tooHeavy"])
			r.OK = false
			r.KO = false
			r.Sleep = 2
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
		r.Answer = append(r.Answer, setup.Answers["unreachable"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}
	return object.snatchFrom(Object(setup.GetObjectByID(10)))
}

func (object Object) Stab(area setup.Area) (r reaction) {
	if !Object(setup.GetObjectByID(15)).inInventory() &&
		!Object(setup.GetObjectByID(25)).inInventory() &&
		!Object(setup.GetObjectByID(33)).inInventory() {
		r.Answer = append(r.Answer, setup.Answers["noTool"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}
	switch object.ID {
	case 14:
		r.Answer = append(r.Answer, setup.Answers["tryCut"])
		r.OK = true
		r.KO = false
		r.Sleep = 2
		return
	case 10:
		r.Answer = append(r.Answer, setup.Answers["stabGrub"])
		r.OK = true
		r.KO = false
		r.Sleep = 2
		return
	case 16:
		r.Answer = append(r.Answer, setup.Answers["starbBaer"])
		r.OK = true
		r.KO = false
		r.Sleep = 2
		return
	case 18:
		if Object(setup.GetObjectByID(13)).inUse() {
			dwarf := setup.GetObjectByID(18)
			dwarf.Properties.Area = -1
			setup.GameObjects[dwarf.ID] = dwarf.Properties
			r.Answer = append(r.Answer, setup.Answers["stabDwarfHooded"])
			r.OK = true
			r.KO = false
			r.Sleep = 1
		} else {
			r.Answer = append(r.Answer, setup.Answers["stabDwarf"])
			r.OK = false
			r.KO = true
			r.Sleep = 2
		}
		return
	case 36:
		if Object(setup.GetObjectByID(13)).inUse() {
			gnome := setup.GetObjectByID(36)
			gnome.Properties.Area = -1
			setup.GameObjects[gnome.ID] = gnome.Properties
			r.Answer = append(r.Answer, setup.Answers["stabGnomeHooded"])
			r.OK = true
			r.KO = false
			r.Sleep = 1
		} else {
			r.Answer = append(r.Answer, setup.Answers["stabGnome"])
			r.OK = false
			r.KO = true
			r.Sleep = 2
		}
		return
	case 42:
		r.Answer = append(r.Answer, setup.Answers["stabDragon"])
		r.OK = false
		r.KO = true
		r.Sleep = 2
		return
	}
	r.Answer = append(r.Answer, setup.Answers["dontKnowHow"])
	r.OK = true
	r.KO = false
	r.Sleep = 2
	return
}

// As Move is called in context of object handling, Move reflects on Object even obj is not used.
func (obj Object) Move(area setup.Area, dir string) (r reaction) {
	//func Move(area int, direction int, text []string) int {
	//if direction == 0 {
	//	return 0, "Ich brauche eine Richtung."
	//}
	// wearing the hood?
	hood := setup.GetObjectByID(13)
	if Object(hood).inUse() {
		hood.Properties.Area = -1
		setup.GameObjects[hood.ID] = hood.Properties
		r.Answer = append(r.Answer, setup.Answers["hoodInUse"])
		r.Sleep = 2
	}

	var direction = map[string]int{"n": 0, "s": 1, "o": 2, "w": 3}

	newArea := area.Properties.Directions[direction[dir]]
	if newArea == 0 {
		r.Answer = append(r.Answer, setup.Answers["noWay"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		r.AreaID = area.ID
		return
		//view.Flash(text, "In diese Richtung führt kein Weg.")
		//return area
	}

	// barefoot on unknown terrain?
	if !Object(setup.GetObjectByID(31)).inUse() {
		r.Answer = append(r.Answer, setup.Answers["noShoes"])
		r.OK = false
		r.KO = true
		r.Sleep = 3
		return
	}

	// Area 30 and 25 are connected by a door. Is it open?
	if (area.ID == 30 || area.ID == 25 && direction[dir] == 0) && !setup.DoorOpen {
		r.Answer = append(r.Answer, setup.Answers["locked"])
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
	r.OK = true
	r.KO = false
	r.AreaID = newArea
	// Direction Moor?
	if newArea == 5 {
		for _, o := range setup.ObjectsInArea(setup.GetAreaByID(1000)) {
			o.Properties.Area = 29
			setup.GameObjects[o.ID] = o.Properties
		}
		r.Answer = append(r.Answer, setup.Answers["inTheMoor"])
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

func (obj Object) Use(area setup.Area) (r reaction) {
	r.KO = false
	r.Sleep = 2
	switch obj.ID {

	case 13:
		if !obj.inInventory() {
			r.Answer = append(r.Answer, setup.Answers["dontHave"])
			r.OK = false
		} else {
			obj.Properties.Area = 2000
			setup.GameObjects[obj.ID] = obj.Properties
			r.Answer = append(r.Answer, setup.Answers["hood"])
			r.OK = true
		}
	case 31:
		if !obj.inInventory() {
			r.Answer = append(r.Answer, setup.Answers["dontHave"])
			r.OK = false
		} else {
			obj.Properties.Area = 2000
			setup.GameObjects[obj.ID] = obj.Properties
			r.Answer = append(r.Answer, setup.Answers["shoes"])
			r.OK = true
		}
	default:
		r.Answer = append(r.Answer, setup.Answers["dontKnowHow"])
		r.OK = false
	}
	return
}

func (object Object) Climb(area setup.Area) (r reaction) {
	if area.ID == 31 {
		moves += 1
		r.OK = true
		r.AreaID = 9
		return
	}
	if area.ID == 9 && object.ID == 27 {
		movement.RevealArea(31)
		moves += 1
		tree := setup.GetAreaByID(31)
		tree.Properties.Visited = true
		setup.GameAreas[tree.ID] = tree.Properties
		r.OK = true
		r.AreaID = 31
		return
	}
	if !object.available(area) {
		r.Answer = append(r.Answer, setup.Answers["dontSee"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		r.AreaID = area.ID
		return
	}
	if object.ID != 27 {
		r.Answer = append(r.Answer, setup.Answers["silly"])
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
		inv = append(inv, setup.Answers["invEmpty"])
		return
	}

	inv = append(inv, setup.Answers["inv"])
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

func (obj Object) pick() (r reaction) {
	inv := setup.GetAreaByID(1000)
	inventory := setup.ObjectsInArea(inv)
	if len(inventory) > 6 {
		r.Answer = append(r.Answer, setup.Answers["invFull"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}

	obj.Properties.Area = 1000
	setup.GameObjects[obj.ID] = obj.Properties
	r.Answer = append(r.Answer, setup.Answers["ok"])
	r.OK = true
	r.KO = false
	r.Sleep = 1
	return
}

func (obj Object) drop(area setup.Area) (r reaction) {
	if obj.Properties.Area != 1000 {
		r.Answer = append(r.Answer, setup.Answers["dontHave"])
		r.OK = false
		r.KO = false
		r.Sleep = 2
		return
	}
	obj.Properties.Area = area.ID
	setup.GameObjects[obj.ID] = obj.Properties
	r.Answer = append(r.Answer, setup.Answers["ok"])
	r.OK = true
	r.KO = false
	r.Sleep = 1
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

func (obj Object) Load(area setup.Area) (r reaction) {
	var content struct {
		AreaID  int                            `yaml:"area"`
		AreaMap [12][10]int                    `yaml:"map"`
		Objects map[int]setup.ObjectProperties `yaml:"objects"`
	}

	filename := folderListing()

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		r.Answer = append(r.Answer, err.Error())
		r.OK = false
		return
	}
	err = yaml.Unmarshal(yamlFile, &content)
	if err != nil {
		r.Answer = append(r.Answer, err.Error())
		r.OK = false
		return
	}

	setup.GameObjects = content.Objects
	setup.Map = content.AreaMap

	r.OK = true
	r.KO = false
	r.AreaID = content.AreaID
	r.Answer = append(r.Answer, "Spiel erfolgreich geladen.")

	return
}
