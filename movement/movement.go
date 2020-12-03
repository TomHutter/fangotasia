package movement

import (
	"fantasia/setup"
	"fantasia/view"
	"fmt"
	"strings"
)

func Surroundings(area setup.Area) (text []string) {
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
	}
	text = append(text, setup.NEUTRAL)
	var items []string
	for _, object := range setup.ObjectsInArea(area) {
		item := view.Highlight(object.Properties.Description.Long, setup.BLUE)
		items = append(items, fmt.Sprintf("%s  - %s", setup.BLUE, item))
	}
	if len(items) > 0 {
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
	text = append(text, fmt.Sprintf("%sGebiet: %d, Richtungen: %s", setup.WHITE, area.ID, strings.Join(directions, ", ")))
	return
}
