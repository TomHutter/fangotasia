package actions

import (
	"fantasia/movement"
	"fantasia/setup"
)

// As Move is called in context of object handling, Move reflects on Object even obj is not used.
func (obj Object) Move(area setup.Area, dir string) (r setup.Reaction, areaID int) {
	//func Move(area int, direction int, text []string) int {
	//if direction == 0 {
	//	return 0, "Ich brauche eine Richtung."
	//}
	// wearing the hood?
	hood := setup.GetObjectByID(13)
	if Object(hood).inUse() {
		hood.Properties.Area = -1
		setup.GameObjects[hood.ID] = hood.Properties
		r = setup.Reactions["hoodInUse"]
	}

	var direction = map[string]int{"n": 0, "s": 1, "o": 2, "w": 3}

	newArea := area.Properties.Directions[direction[dir]]
	if newArea == 0 {
		r = setup.Reactions["noWay"]
		areaID = area.ID
		return
		//view.Flash(text, "In diese Richtung führt kein Weg.")
		//return area
	}

	// barefoot on unknown terrain?
	if !Object(setup.GetObjectByID(31)).inUse() {
		r = setup.Reactions["noShoes"]
		areaID = area.ID
		return
	}

	// Area 30 and 25 are connected by a door. Is it open?
	if (area.ID == 30 || area.ID == 25 && direction[dir] == 0) && !setup.Flags["DoorOpen"] {
		r = setup.Reactions["locked"]
		areaID = area.ID
		return
		//view.Flash(text, "Die Tür ist versperrt.")
		//return area
	}
	movement.RevealArea(newArea)
	moves += 1
	r.OK = true
	r.KO = false
	areaID = newArea
	// Direction Moor?
	if newArea == 5 {
		for _, o := range setup.ObjectsInArea(setup.GetAreaByID(setup.INVENTORY)) {
			o.Properties.Area = 29
			setup.GameObjects[o.ID] = o.Properties
		}
		r = setup.Reactions["inTheMoor"]
	}
	return
}

func (object Object) Climb(area setup.Area) (r setup.Reaction, areaID int) {
	if area.ID == 31 {
		moves += 1
		r.OK = true
		areaID = 9
		return
	}
	if area.ID == 9 && object.ID == 27 {
		movement.RevealArea(31)
		moves += 1
		treetop := setup.GetAreaByID(31)
		treetop.Properties.Visited = true
		setup.GameAreas[treetop.ID] = treetop.Properties
		r.OK = true
		areaID = 31
		return
	}
	if !object.available(area) {
		r = setup.Reactions["dontSee"]
		areaID = area.ID
		return
	}
	if object.ID != 27 {
		r = setup.Reactions["silly"]
		areaID = area.ID
		return
	}
	return
}