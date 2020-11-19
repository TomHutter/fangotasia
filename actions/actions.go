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
type storage []config.Object

var inventory storage
var inUse storage
var object string
var moves int
var visited [51]bool

//type verb []string

//type command struct{}

func Parse(input string, area int, text []string) int {

	var command string
	var order verb
	var knownVerb config.Verb

	re := regexp.MustCompile("[\\s,\\t]+")

	parts := re.Split(input, -1)
	command = strings.ToLower(parts[0])
	// pop parts
	parts = parts[1:]

	fmt.Println(parts) // ["Have", "a", "great", "day!"]
	//parts := strings.r  Split(input, )
	for _, v := range config.Verbs {
		if v.Name == command {
			knownVerb = v
			break
		}
	}

	if knownVerb == (config.Verb{}) {
		notice := fmt.Sprintf("'%s' kenne ich nicht.\nDas Kommando 'Verben' gibt eine Liste aller verfügbaren Verben aus.\n", command)
		view.Flash(text, notice, knownVerb.Sleep, config.RED)
		return area
	}

	fmt.Printf("Valid verb '%s' found.\n", knownVerb.Name)
	if knownVerb.Single {
		call := reflect.ValueOf(&order).MethodByName(knownVerb.Func)
		//view.Flash(text, fmt.Sprintf("%s", call.IsNil()), 2, config.RED)
		//fmt.Println(call.String())
		//fmt.Println(call.IsValid())
		if !call.IsValid() {
			view.Flash(text, fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, config.RED)
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
		view.Flash(text, strings.Join(notice, "\n"), knownVerb.Sleep, config.GREEN)
		//for i := 0; i < val[0].Len(); i++ {
		//	fmt.Println(val[0].Index(i).String())
		//}
		return area
	}

	argv := make([]reflect.Value, 0)

	switch knownVerb.Name {
	case "n", "s", "o", "w":
		argv = append(argv, reflect.ValueOf(area))
		argv = append(argv, reflect.ValueOf(knownVerb.Name))
	default:
		// call reflec.Call for verb with arguments
		// - first valid noun
		// - area
		//argv := make([]reflect.Value, 1)
		// loop over all input parts and see if we find a valid object
		for _, p := range parts {
			word := strings.ToLower(p)
			for _, o := range config.GameObjects {
				if strings.ToLower(o.Description.Short) == word {
					fmt.Printf("Valid object '%s' found.\n", o.Description.Short)
					argv = append(argv, reflect.ValueOf(o))
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
	}

	if len(argv) < 2 {
		notice := fmt.Sprintf("'%s' kenne ich nicht.\n", strings.Join(parts, " "))
		view.Flash(text, notice, knownVerb.Sleep, config.RED)
		return area
	}

	// now method and all args should be known
	call := reflect.ValueOf(&order).MethodByName(knownVerb.Func)
	//view.Flash(text, fmt.Sprintf("%s", call.IsNil()), 2, config.RED)
	//fmt.Println(call.String())
	//fmt.Println(call.IsValid())
	if !call.IsValid() {
		view.Flash(text, fmt.Sprintf("Func '%s' not yet implemented\n", knownVerb.Func), 2, config.RED)
		return area
	}
	//val := reflect.ValueOf(&order).MethodByName(knownVerb.Func).Call(argv)
	val := call.Call(argv)

	var notice []string
	var color string
	switch knownVerb.Func {
	case "Move":
		if val[0].Bool() == true {
			fmt.Printf("New area: %d\n", val[1].Int())
			return int(val[1].Int())
		}
		for i := 0; i < val[2].Len(); i++ {
			notice = append(notice, val[2].Index(i).String())
		}
		view.Flash(text, strings.Join(notice, "\n"), knownVerb.Sleep, config.RED)
	case "Stab":
		for i := 0; i < val[1].Len(); i++ {
			notice = append(notice, val[1].Index(i).String())
		}
		if val[0].Bool() == false {
			view.Flash(text, strings.Join(notice, "\n"), knownVerb.Sleep, config.RED)
			order.GameOver()
		} else {
			view.Flash(text, strings.Join(notice, "\n"), knownVerb.Sleep, config.GREEN)
		}
	default:
		fmt.Println(val[0])
		if val[0].Bool() == false {
			color = config.RED
		} else {
			color = config.GREEN
		}
		for i := 0; i < val[1].Len(); i++ {
			//fmt.Println(val[1].Index(i).String())
			notice = append(notice, val[1].Index(i).String())
			view.Flash(text, strings.Join(notice, "\n"), knownVerb.Sleep, color)
		}
	}
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

func objectInArea(object config.Object, area int) bool {
	return object.Area == area
}

func objectInInventory(object config.Object) bool {
	return object.Area == -1
}

func objectInUse(object config.Object) bool {
	return object.Area == -2
}

func objectAvailable(object config.Object, area int) bool {
	return objectInArea(object, area) || objectInInventory(object) || objectInUse(object)
}

func (inv *storage) add(object config.Object) (ok bool, answer []string) {
	if len(*inv) > 6 {
		answer = append(answer, config.Answers[6])
		return false, answer
	}
	*inv = append(*inv, object)
	o := config.GetObjectByID(object.ID)
	o.Area = -1
	answer = append(answer, config.Answers[7])
	return true, answer
}

func (inv storage) drop(object config.Object, area int) (ok bool, answer []string) {
	var newInv storage
	for _, o := range inv {
		if o.ID == object.ID {
			obj := config.GetObjectByID(object.ID)
			obj.Area = area
			continue
		}
		newInv = append(newInv, o)
	}
	copy(inv, newInv)
	answer = append(answer, config.Answers[7])
	return true, answer
}

func invisible(area int, opponent config.Object, object config.Object) (ok bool, answer []string) {
	if opponent.Area == area && !objectInUse(*config.GetObjectByID(13)) {
		answer = append(answer, fmt.Sprintf("%s %s %s",
			strings.Title(opponent.Description.Article),
			opponent.Description.Short,
			config.Answers[3]))
		return false, answer
	}
	return inventory.add(object)
}

func (v *verb) Open(object config.Object, area int) (ok bool, answer []string) {
	answer = append(answer, "lässt sich öffnen")
	return true, answer
}

func (v *verb) Take(object config.Object, area int) (ok bool, answer []string) {
	if !objectAvailable(object, area) {
		answer = append(answer, "sehe ich hier nicht.")
		return false, answer
	}
	if objectInInventory(object) || objectInUse(object) {
		answer = append(answer, "habe ich schon.")
		return false, answer
	}

	switch object.ID {
	case 10, 16, 18, 21, 22, 27, 36, 40, 42:
		answer = append(answer, config.Answers[0])
		return false, answer
	case 29, 14:
		answer = append(answer, config.Answers[4])
		return false, answer
	case 34:
		if !objectInUse(*config.GetObjectByID(9)) {
			answer = append(answer, config.Answers[4])
			return false, answer
		}
	case 17:
		return invisible(area, *config.GetObjectByID(16), object)
	case 19:
		return invisible(area, *config.GetObjectByID(18), object)
	case 35:
		return invisible(area, *config.GetObjectByID(36), object)
	case 44:
		return invisible(area, *config.GetObjectByID(42), object)
	case 32, 43:
		answer = append(answer, config.Answers[5])
		return false, answer
	}
	return invisible(area, *config.GetObjectByID(10), object)
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

func (v *verb) Stab(object config.Object, area int) (ok bool, answer []string) {
	if !objectInInventory(*config.GetObjectByID(15)) &&
		!objectInInventory(*config.GetObjectByID(25)) &&
		!objectInInventory(*config.GetObjectByID(33)) {
		answer = append(answer, config.Answers[0])
		return false, answer
	}
	switch object.ID {
	case 14:
		answer = append(answer, config.Answers[10])
		return true, answer
	case 10:
		answer = append(answer, config.Answers[11])
		return true, answer
	case 16:
		answer = append(answer, config.Answers[12])
		return true, answer
	case 18:
		answer = append(answer, config.Answers[13])
		return false, answer
	case 36:
		answer = append(answer, config.Answers[14])
		return false, answer
	case 42:
		answer = append(answer, config.Answers[15])
		return false, answer
	}
	answer = append(answer, config.Answers[1])
	return true, answer
	/*
		345 rem ** stich **********************
		346 f=0:gosub605:iffl=1thenfl=0:goto280
		347 ifno<>10andno<>14andno<>16andno<>18andno<>36andno<>42thenfl=1
		348 iffl=1thenfl=0:printa$(2):goto280
		349 ifno=14thenprint"versuche 'schneide'.":goto280
		350 ifge(15)<>-1andge(25)<>-1andge(33)<>-1thenfl=1
		351 iffl=1thenfl=0:printa$(1):goto280
		352 ifno=10thenprint"die raupe ist kitzelig und lacht laut !":goto281
		353 ifno=16thenprint"der baer brummt unwillig.":goto281
		354 ifno=18thenprint"der zwerg wird boes und toetet mich !":goto357
		355 ifno=36thenprint"der gnom verzaubert mich !":goto357
		356 ifno=42thenprint"der drache verbrennt mich !"
		357 fori=1to2000:next:goto611
	*/
}

func (v *verb) Move(area int, dir string) (ok bool, newArea int, answer []string) {
	//func Move(area int, direction int, text []string) int {
	//if direction == 0 {
	//	return 0, "Ich brauche eine Richtung."
	//}
	var direction = map[string]int{"n": 0, "s": 1, "o": 2, "w": 3}

	newArea = config.GetAreaByID(area).Directions[direction[dir]]
	if newArea == 0 {
		answer = append(answer, config.Answers[8])
		return false, area, answer
		//view.Flash(text, "In diese Richtung führt kein Weg.")
		//return area
	}
	// Area 30 and 25 are connected by a door. Is it open?
	if (area == 30 || area == 25 && direction[dir] == 0) && !config.DoorOpen {
		answer = append(answer, config.Answers[9])
		return false, area, answer
		//view.Flash(text, "Die Tür ist versperrt.")
		//return area
	}
	movement.RevealArea(newArea)
	answer = append(answer, config.Answers[7])
	moves += 1
	visited[area] = true
	return true, newArea, answer
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

func use(object int, area int) {
	/*
		604 rem ** unterprogramm **************
		605 ifge(no)<>oaandge(no)<>-1andge(no)<>-2thenfl=1
		606 iffl=1thenprint"sehe ich hier nicht.":return
		607 iffl=1andge(no)<>-1andge(no)<>-2thenprint"habe ich nicht dabei.":fl=1
		608 return
		if objectsInArea[object][0] != area
	*/
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
	if len(inventory) == 0 {
		//fmt.Println("Ich habe nichts dabei.")
		inv = append(inv, "Ich habe nichts dabei.")
		return
	}

	inv = append(inv, "Ich habe:")
	for _, i := range inventory {
		inv = append(inv, fmt.Sprintf("- %s", i.Description.Long))
	}
	return
}

func (v *verb) GameOver() {
	var board []string
	sum := 0
	for _, o := range inventory {
		sum += o.Value
	}
	// all valuable objects found
	if sum == 170 {
		switch {
		case moves < 500:
			sum += 7
			fallthrough
		case moves < 400:
			sum += 7
			fallthrough
		case moves < 300:
			sum += 7
		}
	}
	switch {
	case visited[5]:
		sum += 2
	case visited[29]:
		sum += 3
	case visited[31]:
		sum += 4
	}
	board = append(board, fmt.Sprintf("Du hast %d von 200 Punkten!\n", sum))
	board = append(board, "Noch ein Spiel (j/n)?\n")
	view.Flash(board, "", -1, "")
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
