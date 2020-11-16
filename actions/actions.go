package actions

import (
	"fantasia/config"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var inventory []config.Object
var inUse []config.Object
var object string

type verb config.Verb

//type verb []string

//type command struct{}

func Parse(input string, area int) {

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
		fmt.Printf("'%s' kenne ich nicht.\nDas Kommando 'Verben' gibt eine Liste aller verfügbaren Verben aus.\n", command)
		return
	}

	fmt.Printf("Valid verb '%s' found.\n", knownVerb.Name)
	if knownVerb.Single {
		// call reflec.Call for verb without arguments
		val := reflect.ValueOf(&order).MethodByName(knownVerb.Func).Call([]reflect.Value{})
		for i := 0; i < val[0].Len(); i++ {
			fmt.Println(val[0].Index(i).String())
		}
		return
	}
	// call reflec.Call for verb with arguments
	// - first valid noun
	// - area
	argv := make([]reflect.Value, 2)
	//argv := make([]reflect.Value, 1)
	for _, p := range parts {
		noun := strings.ToLower(p)
		for _, n := range config.GameObjects {
			if strings.ToLower(n.Description.Short) == noun {
				fmt.Printf("Valid noun '%s' found.\n", n.Description.Short)
				argv[0] = reflect.ValueOf(n)
				argv[1] = reflect.ValueOf(area)

				//val := reflect.ValueOf(&order).MethodByName(v.Func).Call(argv)
				//val := reflect.ValueOf(&order).MethodByName(v.Func).Call([]reflect.Value{reflect.ValueOf(order)})
				val := reflect.ValueOf(&order).MethodByName(knownVerb.Func).Call(argv)
				for i := 0; i < val[0].Len(); i++ {
					fmt.Println(val[0].Index(i).String())
				}
				fmt.Println(val[1])
				return
			}
		}
	}

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
	for _, o := range config.ObjectsInArea(area) {
		if o.Description.Short == object.Description.Short {
			return true
		}
	}
	return false
}

func objectInInventory(object config.Object) bool {
	for _, o := range inventory {
		if o.Description.Short == object.Description.Short {
			return true
		}
	}
	return false
}

func objectInUse(object config.Object) bool {
	for _, o := range inUse {
		if o.Description.Short == object.Description.Short {
			return true
		}
	}
	return false
}

func objectAvailable(object config.Object, area int) bool {
	return objectInArea(object, area) || objectInInventory(object) || objectInUse(object)
}

func (v *verb) Open(object config.Object, area int) (answer []string, ok bool) {
	answer = append(answer, "lässt sich öffnen")
	return answer, true
}

func (v *verb) Take(object config.Object, area int) (answer []string, ok bool) {
	if !objectAvailable(object, area) {
		answer = append(answer, "sehe ich hier nicht.")
		return answer, false
	}
	if objectInInventory(object) || objectInUse(object) {
		answer = append(answer, "habe ich schon.")
		return answer, false
	}

	switch object.ID {
	case 10, 16, 18, 21, 22, 27, 36, 40, 42:
		answer = append(answer, config.Answers[0])
		return answer, false
	case 29, 14:
		answer = append(answer, config.Answers[4])
		return answer, false
	case 34:
		if !objectInUse(config.GetObjectByID(9)) {
			answer = append(answer, config.Answers[4])
			return answer, false
		}
	case 17:
		opponent := config.GetObjectByID(16)
		if opponent.Area == area && !objectInUse(config.GetObjectByID(13)) {
			answer = append(answer, fmt.Sprintf("%s %s %s",
				strings.Title(opponent.Description.Article),
				opponent.Description.Short,
				config.Answers[3]))
			return answer, false
		}
	case 19:
		opponent := config.GetObjectByID(18)
		if opponent.Area == area && !objectInUse(config.GetObjectByID(13)) {
			answer = append(answer, fmt.Sprintf("%s %s %s",
				strings.Title(opponent.Description.Article),
				opponent.Description.Short,
				config.Answers[3]))
			return answer, false
		}
	case 35:
		opponent := config.GetObjectByID(36)
		if opponent.Area == area && !objectInUse(config.GetObjectByID(13)) {
			answer = append(answer, fmt.Sprintf("%s %s %s",
				strings.Title(opponent.Description.Article),
				opponent.Description.Short,
				config.Answers[3]))
			return answer, false
		}
	case 44:
		opponent := config.GetObjectByID(42)
		if opponent.Area == area && !objectInUse(config.GetObjectByID(13)) {
			answer = append(answer, fmt.Sprintf("%s %s %s",
				strings.Title(opponent.Description.Article),
				opponent.Description.Short,
				config.Answers[3]))
			return answer, false
		}
	case 32, 43:
		answer = append(answer, config.Answers[5])
		return answer, false
	}
	opponent := config.GetObjectByID(10)
	if opponent.Area == area && !objectInUse(config.GetObjectByID(13)) {
		answer = append(answer, fmt.Sprintf("%s %s %s",
			strings.Title(opponent.Description.Article),
			opponent.Description.Short,
			config.Answers[3]))
		return answer, false
	}
	if len(inventory) > 6 {
		answer = append(answer, config.Answers[6])
		return answer, false
	}
	inventory = append(inventory, object)
	answer = append(answer, config.Answers[7])
	return answer, true

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

	for _, i := range inventory {
		inv = append(inv, i.Description.Short)
	}
	return
}

//func (obj object) hit(inv []string) {
//}
