package movement

import (
	"fantasia/setup"
	"fmt"
	"strings"
)

func Surroundings(area setup.Area) (text []string) {
	/*
		if area == 25 {
			setup.ObjectsInArea[40][0] = 25
			setup.Objects[40-9] = "eine Tür im Norden"
		}
		if area == 30 {
			setup.ObjectsInArea[40][0] = 30
			setup.Objects[40-9] = "eine Tür im Süden"
		}
	*/

	//	thenge(40)=25:ge$(40)="eine tuer im norden"
	//	ifoa=30thenge(40)=30:ge$(40)="eine tuer im sueden"
	//fmt.Printf("Ich bin %s\n", locations[area-1])
	//var text []string
	desc := strings.Split(area.Properties.Description.Long, "\\n")
	desc0, desc := desc[0], desc[1:]
	text = append(text, fmt.Sprintf("%sIch bin %s", setup.YELLOW, desc0))
	for _, v := range desc {
		if strings.Contains(v, "++") {
			v = strings.ReplaceAll(v, "++", "")
			text = append(text, fmt.Sprintf("%s%s", setup.CYAN, v))
		} else {
			text = append(text, fmt.Sprintf("%s%s", setup.YELLOW, v))
		}
		//text = append(text, v)
	}
	text = append(text, setup.NEUTRAL)
	//appendText(&text, fmt.Sprintf("Ich bin %s", locations[area-1]), yellow)
	var items []string
	for _, object := range setup.ObjectsInArea(area) {
		item := object.Properties.Description.Long
		if strings.Contains(item, "::") {
			item = strings.ReplaceAll(item, "::", "")
			items = append(items, fmt.Sprintf("%s  - %s", setup.WHITE, item))
			//appendText(&text, fmt.Sprintf("  - %s", item), white)
		} else {
			items = append(items, fmt.Sprintf("%s  - %s", setup.BLUE, item))
			//appendText(&text, fmt.Sprintf("  - %s", item), blue)
		}
	}
	if len(items) > 0 {
		//text = append(text, "")
		text = append(text, fmt.Sprintf("%sIch sehe:", setup.BLUE))
		for _, item := range items {
			text = append(text, item)
		}
		text = append(text, setup.NEUTRAL)
	}
	var directions []string
	for d := 0; d < 4; d++ {
		if area.Properties.Directions[d] != 0 {
			switch d {
			case 0: // N
				directions = append(directions, "Norden")
			case 1: // S
				directions = append(directions, "Süden")
			case 2: // O
				directions = append(directions, "Osten")
			case 3: // W
				directions = append(directions, "Westen")
			}
		}
	}
	//text = append(text, "")
	text = append(text, fmt.Sprintf("%sRaum: %d, Richtungen: %s", setup.WHITE, area.ID, strings.Join(directions, ", ")))
	//printScreen(text)
	return
	//appendText(&text, "Ich sehe:", blue)
	/*if v[1] == area {
		//items = append(items, objects[i-9])
		appendText(&text, fmt.Sprintf("  - %s", objects[i-9]), red)
	}*/
	//printScreen(text)
	/*if len(items) > 0 {
		fmt.Println("Ich sehe:")
		for _, i := range items {
			fmt.Printf("  - %s\n", i)
		}
	}*/

	/*
		ifoa=25thenge(40)=25:ge$(40)="eine tuer im norden"
		ifoa=30thenge(40)=30:ge$(40)="eine tuer im sueden"
		ifoa=6thenp1=1
		ifoa=31thenp2=1
		ifoa=29thenp3=1
		ifoa<>1andge(31)<>-2thenprinte$:poke214,5:poke211,3:sysvd:fl=1
		iffl=1thenprint"hilfe !   ich versinke im boden."
		iffl=1thenfl=0:pokevc,peek(vc)or16:fori=1to2000:next:goto611
		printc$"ich bin "o$(oa)d$:fl=0:fori=9to44:ifge(i)<>oathen323
		iffl=0thenprintf$"ich sehe:"
		printge$(i):fl=1
		next:fl=0
		ifoa=31then335
		ifin>1andoa=5then327
		goto331
		fori=9to44:if(ge(i)=-1orge(i)=-2)andi<>31thenge(i)=29
		next:in=1
		print"im moor ist alles verschwunden,"
		print"was ich bei mir hatte !"
		fl=0:printf$"richtungen:":fori=0to3:ifr(oa,i)=0then334
		iffl=1thenprint", ";
		printno$(i+5);:fl=1
		next:fl=0
		printtc$f$:fori=1to40:printchr$(175);:next:printd$;:return
	*/
}

/*
func Move(area int, direction int, text []string) int {
	//if direction == 0 {
	//	return 0, "Ich brauche eine Richtung."
	//}
	newArea := setup.Areas[area][direction]
	if newArea == 0 {
		view.Flash(text, "In diese Richtung führt kein Weg.")
		return area
	}
	// Area 30 and 25 are connected by a door. Is it open?
	if (area == 30 || area == 25 && direction == 0) && !doorOpen {
		view.Flash(text, "Die Tür ist versperrt.")
		return area
	}
	RevealArea(newArea)
	return newArea
}
*/
/*
	area = movement.Move(area, direction, text)
	// are we lost? (show old area)
	if !movement.AreaVisible(area) {
		text = movement.DrawMap(oldArea)
		//text = append(text, "\n", "\n", "\n")
		text = append(text, movement.Surroundings(oldArea)...)
		view.PrintScreen(text)
	} else {
		//text = drawMap(area)
		//text = surroundings(area, locations, objects)
		text = movement.DrawMap(area)
		//text = append(text, "\n", "\n", "\n")
		text = append(text, movement.Surroundings(area)...)
		oldArea = area
		view.PrintScreen(text)
	}
*/
