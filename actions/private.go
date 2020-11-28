package actions

import (
	"fantasia/setup"
	"fmt"
	"strings"
)

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

func (object Object) snatchFrom(opponent Object) (r setup.Reaction) {
	hood := Object(setup.GetObjectByID(13))
	area := setup.GetAreaByID(object.Properties.Area)
	if opponent.inArea(area) && !hood.inUse() {
		r = setup.Reactions["wontLet"]
		r.Statement = fmt.Sprintf("%s %s %s",
			strings.Title(opponent.Properties.Description.Article),
			opponent.Properties.Description.Short,
			r.Statement)
		return
	}
	return object.pick()
}

func (obj Object) pick() (r setup.Reaction) {
	inv := setup.GetAreaByID(1000)
	inventory := setup.ObjectsInArea(inv)
	if len(inventory) > 6 {
		r = setup.Reactions["invFull"]
		return
	}

	obj.Properties.Area = 1000
	setup.GameObjects[obj.ID] = obj.Properties
	r = setup.Reactions["ok"]
	return
}

func (obj Object) drop(area setup.Area) (r setup.Reaction) {
	if obj.Properties.Area != 1000 {
		r = setup.Reactions["dontHave"]
		return
	}
	obj.Properties.Area = area.ID
	setup.GameObjects[obj.ID] = obj.Properties
	r = setup.Reactions["ok"]
	return
}
