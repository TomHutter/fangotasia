package view

import (
	"fangotasia/setup"
	"fmt"
	"regexp"
	"strings"
)

var Notice struct {
	Message string
	Color   string
	Sleep   int
}

func Highlight(s string, c string) (h string) {
	re := regexp.MustCompile("::(.+)::")

	h = string(re.ReplaceAll([]byte(s), []byte("[white::b]"+"$1"+c)))
	return
}

func Surroundings(area setup.Area) (text []string) {
	desc := strings.Split(area.Properties.Description.Long, "\\n")
	desc0, desc := desc[0], desc[1:]
	text = append(text, fmt.Sprintf("%sIch bin %s", "[yellow]", desc0))
	for _, v := range desc {
		if strings.Contains(v, "++") {
			v = strings.ReplaceAll(v, "++", "")
			text = append(text, fmt.Sprintf("%s%s", "[cyan]", v))
		} else {
			text = append(text, fmt.Sprintf("%s%s", "[yellow]", v))
		}
	}
	text = append(text, "[-:-:-]")
	var items []string
	for _, object := range setup.ObjectsInArea(area) {
		item := Highlight(object.Properties.Description.Long, "[blue:black:-]")
		items = append(items, fmt.Sprintf("%s  - %s", "[blue:black]", item))
	}
	if len(items) > 0 {
		text = append(text, fmt.Sprintf("%sIch sehe:", "[blue:black]"))
		for _, item := range items {
			text = append(text, item)
		}
		text = append(text, "[-:-:-]")
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
	text = append(text, fmt.Sprintf("%sRichtungen: %s", "[white:black:b]", strings.Join(directions, ", ")))
	return
}
