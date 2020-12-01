package actions

import (
	"fantasia/movement"
	"fantasia/setup"
)

func cycles(areaID int) {
	if areaID == 1 {
		setup.Cycles = 1
	} else {
		setup.Cycles += areaID
	}

	if setup.Cycles == 108 {
		imke := Object(setup.GetObjectByID(48))
		imke.Properties.Description.Long = "\033[01;95mIMKE\033[0m den pink Diamanten"
		imke.Properties.Description.Short = "Imke"
		imke.NewAreaID(areaID)
		//setup.GameObjects[imke.ID] = imke.Properties
	}
}

// As Move is called in context of object handling, Move reflects on Object even obj is not used.
func (obj Object) Move(area setup.Area, dir string) (r setup.Reaction, areaID int) {
	//func Move(area int, direction int, text []string) int {
	//if direction == 0 {
	//	return 0, "Ich brauche eine Richtung."
	//}
	// wearing the hood?
	hood := Object(setup.GetObjectByID(13))
	if Object(hood).inUse() {
		hood.NewAreaID(0)
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

	cycles(newArea)

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
			obj := Object(o)
			obj.NewAreaID(29)
		}
		r = setup.Reactions["inTheMoor"]
	}
	return
}

func (object Object) Climb(area setup.Area) (r setup.Reaction, areaID int) {
	if area.ID == 31 {
		cycles(area.ID)
		moves += 1
		r.OK = true
		areaID = 9
		return
	}
	if area.ID == 9 && object.ID == 27 {
		cycles(area.ID)
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
