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
	desc := strings.Split(area.Properties.Description[setup.Language].Long, "\\n")
	desc0, desc := desc[0], desc[1:]
	text = append(text, fmt.Sprintf("%s%s %s", "[yellow]", setup.TextElements["iAm"][setup.Language], desc0))
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
		// skip hood, if vanished
		if object.ID == 13 && setup.Flags["HoodVanished"] {
			continue
		}
		item := Highlight(object.Properties.Description[setup.Language].Long, "[blue:black:-]")
		items = append(items, fmt.Sprintf("%s  - %s", "[blue:black]", item))
	}
	if len(items) > 0 {
		text = append(text, fmt.Sprintf("%s%s:", "[blue:black]", setup.TextElements["iSee"][setup.Language]))
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
				directions = append(directions, setup.TextElements["north"][setup.Language])
			case 1: // S
				directions = append(directions, setup.TextElements["south"][setup.Language])
			case 2: // O
				directions = append(directions, setup.TextElements["east"][setup.Language])
			case 3: // W
				directions = append(directions, setup.TextElements["west"][setup.Language])
			}
		}
	}
	text = append(text, fmt.Sprintf("%s%s: %s", "[white:black:b]",
		setup.TextElements["directions"][setup.Language],
		strings.Join(directions, ", ")))
	return
}
