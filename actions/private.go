package actions

import (
	"fangotasia/grid"
	"fangotasia/setup"
	"fmt"
	"strings"
)

func setResponse(response string) {
	grid.App.QueueUpdate(func() {
		grid.InputField.SetText("")
		grid.Response.SetText(response)
	})
}

func (object Object) inArea(area setup.Area) bool {
	return object.Properties.Area == area.ID
}

func (object Object) inInventory() bool {
	return object.Properties.Area == setup.INVENTORY
}

func (object Object) inUse() bool {
	return object.Properties.Area == setup.INUSE
}

func (object Object) available(area setup.Area) bool {
	return object.inArea(area) || object.inInventory() || object.inUse()
}

func (object Object) snatchFrom(opponent Object) (r setup.Reaction) {
	hood := Object(setup.GetObjectByID(13))
	area := setup.GetAreaByID(object.Properties.Area)
	if opponent.inArea(area) {
		if !hood.inUse() {
			r = setup.GetReactionByName("wontLet")
			r.Statement[setup.Language][0] = fmt.Sprintf("%s %s %s",
				strings.Title(opponent.Properties.Description[setup.Language].Article),
				opponent.Properties.Description[setup.Language].Short,
				r.Statement[setup.Language][0])
			return
		} else {
			// hood in use and object picked successfully?
			r = object.pick()
			if r.OK {
				r = setup.Reactions["hoodInUse"]
				setup.Flags["HoodVanished"] = true
				hood.NewAreaID(object.Properties.Area)
				return
			}
		}
	}
	return object.pick()
}

func (obj Object) pick() (r setup.Reaction) {
	inv := setup.GetAreaByID(setup.INVENTORY)
	inventory := setup.ObjectsInArea(inv)
	if len(inventory) > 6 {
		r = setup.Reactions["invFull"]
		return
	}

	r = setup.Reactions["ok"]
	obj.NewAreaID(setup.INVENTORY)
	return
}

func (obj Object) drop(area setup.Area) (r setup.Reaction) {
	if obj.Properties.Area != setup.INVENTORY {
		r = setup.Reactions["dontHave"]
		return
	}
	obj.NewAreaID(area.ID)
	r = setup.Reactions["ok"]
	return
}

func getObjectByName(name string, area setup.Area) (ok bool, object Object) {
	found := false
	var obj Object
	for id, prop := range setup.GameObjects {
		if strings.ToLower(prop.Description[setup.Language].Short) == strings.ToLower(name) ||
			strings.ToLower(prop.Description[setup.Language].Alt) == strings.ToLower(name) {
			found = true
			obj = Object{id, prop}
			if obj.available(area) {
				return found, obj
			}
		}
	}
	// on treetop climbing down?
	if area.ID == 31 && obj.ID == 27 {
		return found, obj
	}
	// found an object but not in the current area?
	if found {
		r := setup.Reactions["dontSee"]
		grid.InputField.SetText("")
		grid.Response.SetText(
			fmt.Sprintf("\n%s%s%s\n",
				"[red]",
				r.Statement[setup.Language][0],
				"[-:black:-]"))
		return
	}
	// don't know what you are talking about
	r := setup.GetReactionByName("unknownNoun")
	statement := fmt.Sprintf(r.Statement[setup.Language][0], name)
	grid.InputField.SetText("")
	grid.Response.SetText(
		fmt.Sprintf("\n%s%s%s\n",
			"[red]",
			statement,
			"[-:black:-]"))
	return
}
