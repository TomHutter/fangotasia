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
	switch knownVerb.Func {
	case "Move", "Climb":
		// Answer
		for i := 0; i < val[0].Field(2).Len(); i++ {
			notice = append(notice, val[0].Field(2).Index(i).String())
		}
		sleep := int(val[0].Field(3).Int())
		view.AddFlashNotice(strings.Join(notice, "\n"), sleep, config.RED)
		view.FlashNotice()
		// OK
		if val[0].Field(0).Bool() == true {
			//fmt.Printf("New area: %d\n", val[0].Field(4).Int())
			// Area
			return config.GetAreaByID(int(val[0].Field(4).Int()))
		}
		// KO
		if val[0].Field(1).Bool() == true {
			GameOver()
		}
	//case "Stab":
	default:
		// OK
		if val[0].Field(0).Bool() == true {
			color = config.GREEN
		} else {
			color = config.RED
		}
		// Sleep
		sleep := val[0].Field(3).Int()
		// Answer
		for i := 0; i < val[0].Field(2).Len(); i++ {
			notice = append(notice, val[0].Field(2).Index(i).String())
		}
		view.AddFlashNotice(strings.Join(notice, "\n"), int(sleep), color)
		view.FlashNotice()
		// KO
		if val[0].Field(1).Bool() == true {
			GameOver()
		}
	}
	/*
		default:
			fmt.Println(val[0])
			sleep := 1
			if val[0].Bool() == false {
				color = config.RED
				sleep = knownVerb.Sleep
			} else {
				color = config.GREEN
			}
			for i := 0; i < val[1].Len(); i++ {
				//fmt.Println(val[1].Index(i).String())
				notice = append(notice, val[1].Index(i).String())
				view.Flash(text, strings.Join(notice, "\n"), sleep, color)
			}
		}*/
	return area

	/*
			279 rem ** kommandoabfrage ************
		280 ifge(13)=-2thenge(13)=-4
		281 ifoa<>olorve=3thenprinte$:ol=oa:gosub311
		282 ifge(13)=-2thenprintf$"die tarnkappe hat sich in luft"
		283 ifge(13)=-2thenprint"aufgeloest !":ge(13)=2:in=in-1
		284 ifge(13)=-4thenge(13)=-2
		285 pokevc,peek(vc)or16
		286 ze=ze+1:ko$="":printc$f$"und nun ";:inputko$:printd$;
		287 iflen(ko$)=0thenprintchr$(145)chr$(145);:goto286
		288 v$="":n$="":ve=0:no=0
		289 fori=1tolen(ko$)
		290 ifmid$(ko$,i,1)<>" "thenv$=v$+mid$(ko$,i,1):next
		291 iflen(v$)+1>=len(ko$)then293
		292 n$=right$(ko$,(len(ko$)-i))
		293 fori=1to31:ifv$=ve$(i)thenve=i:goto297
		294 next
		295 ifn$=""thenn$=v$
		296 goto298
		297 ifv$=ko$thengoto300
		298 fori=1to44:ifn$=no$(i)thenno=i:goto300
		299 next
		300 ifno<9andno<>0and(ve=0orve=1)thenve=1
		301 iflen(n$)>0andno=0andve<>16thenve=0:fl=1
		302 iffl=1thenfl=0:printchr$(34);n$;chr$(34);" kenne ich nicht.":goto280
		303 ifn$=""and(ve>8orve=2)andve<>22andve<>15andve<>30thenfl=1
		304 iffl=1thenfl=0:print"bitte gib ein objekt an.":goto280
		305 ifve=0thenprintchr$(34);v$;chr$(34);"kenne ich nicht.":goto280
		306 onvegoto338,346,280,376,567,586,390,397,403,413,417
		307 onve-11goto438,446,455,461,466,555,471,477,487,496
		308 onve-21goto383,504,510,517,360,529,535,539,547,517
		309 :
	*/
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
	/*
		opponent := config.GetObjectByID(10)
		if opponent.Area == area && !objectInUse(*config.GetObjectByID(13)) {
			answer = append(answer, fmt.Sprintf("%s %s %s",
				strings.Title(opponent.Description.Article),
				opponent.Description.Short,
				config.Answers[3]))
			return false, answer
		}
		return inventory.add(object)
	*/
	/*
			359 rem ** nimm ***********************
		360 f=0:gosub605:iffl=1thenfl=0:goto280
		361 ifge(no)=-1orge(no)=-2thenprint"habe ich dabei.":goto280
		362 ifno=10orno=16orno=18orno=21orno=22orno=27orno=36orno=40orno=42thenfl=1
		363 iffl=1thenfl=0:printa$(0):goto280
		364 ifno=29orno=14or(no=34andge(9)<>-3)thenprint"zu schwer.":goto280
		365 ifno=17andge(16)=oaandge(13)<>-2thenprint"der baer";a$(3):goto281
		366 ifno=19andge(18)=oaandge(13)<>-2thenprint"der zwerg";a$(3):goto281
		367 ifno=35andge(36)=oaandge(13)<>-2thenprint"der gnom";a$(3):goto281
		368 ifge(10)=oaandge(13)<>-2thenprint"die raupe";a$(3):goto281
		369 ifno=44andge(42)=oaandge(13)<>-2thenprint"der drache";a$(3):goto281
		370 ifno=32orno=43thenprint"kann ich nicht erreichen.":goto280
		371 ifin+1>7thenprint"ich habe zuviel zu tragen."
		372 ifin+1>7thenprint"ich muesste etwas weglegen.":goto280
		373 in=in+1:ge(no)=-1:print"gut.":goto281
		374 :
	*/
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
		r.Answer = append(r.Answer, config.Answers["stabDwarf"])
		r.OK = false
		r.KO = true
		r.Sleep = 2
		return
	case 36:
		r.Answer = append(r.Answer, config.Answers["stabGnome"])
		r.OK = false
		r.KO = true
		r.Sleep = 2
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
