package movement

import (
	"fantasia/setup"
	"fantasia/view"
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
		item := view.Highlight(object.Properties.Description.Long, setup.BLUE)
		/*
			if strings.Contains(item, "::") {
				item = strings.ReplaceAll(item, "::", "")
				items = append(items, fmt.Sprintf("%s  - %s", setup.WHITE, item))
				//appendText(&text, fmt.Sprintf("  - %s", item), white)
			} else {
				items = append(items, fmt.Sprintf("%s  - %s", setup.BLUE, item))
				//appendText(&text, fmt.Sprintf("  - %s", item), blue)
			}
		*/
		items = append(items, fmt.Sprintf("%s  - %s", setup.BLUE, item))
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
}
