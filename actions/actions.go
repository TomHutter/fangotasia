package actions

import (
	"fantasia/movement"
	"fantasia/setup"
)

func (object Object) Open(area int) (ok bool, answer []string) {
	answer = append(answer, "lässt sich öffnen")
	return true, answer
}

func (object Object) Take(area setup.Area) (r setup.Reaction) {
	if !object.available(area) {
		r = setup.Reactions["dontSee"]
		return
	}
	if object.inInventory() || object.inUse() {
		r = setup.Reactions["haveAlready"]
		return
	}

	switch object.ID {
	case 10, 16, 18, 21, 22, 27, 36, 40, 42:
		r = setup.Reactions["silly"]
		return
	case 29, 14:
		r = setup.Reactions["tooHeavy"]
		return
	case 34:
		if !Object(setup.GetObjectByID(9)).inUse() {
			r = setup.Reactions["tooHeavy"]
			return
		}
	case 17:
		return object.snatchFrom(Object(setup.GetObjectByID(16)))
	case 19:
		return object.snatchFrom(Object(setup.GetObjectByID(18)))
	case 35:
		return object.snatchFrom(Object(setup.GetObjectByID(36)))
	case 44:
		return object.snatchFrom(Object(setup.GetObjectByID(42)))
	case 32, 43:
		r = setup.Reactions["unreachable"]
		return
	}
	return object.snatchFrom(Object(setup.GetObjectByID(10)))
}

func (object Object) Stab(area setup.Area) (r setup.Reaction) {
	if !Object(setup.GetObjectByID(15)).inInventory() &&
		!Object(setup.GetObjectByID(25)).inInventory() &&
		!Object(setup.GetObjectByID(33)).inInventory() {
		r = setup.Reactions["noTool"]
		return
	}
	switch object.ID {
	case 14:
		r = setup.Reactions["tryCut"]
		return
	case 10:
		r = setup.Reactions["stabGrub"]
		return
	case 16:
		r = setup.Reactions["stabBaer"]
		return
	case 18:
		if Object(setup.GetObjectByID(13)).inUse() {
			dwarf := setup.GetObjectByID(18)
			dwarf.Properties.Area = -1
			setup.GameObjects[dwarf.ID] = dwarf.Properties
			r = setup.Reactions["stabDwarfHooded"]
		} else {
			r = setup.Reactions["stabDwarf"]
		}
		return
	case 36:
		if Object(setup.GetObjectByID(13)).inUse() {
			gnome := setup.GetObjectByID(36)
			gnome.Properties.Area = -1
			setup.GameObjects[gnome.ID] = gnome.Properties
			r = setup.Reactions["stabGnomeHooded"]
		} else {
			r = setup.Reactions["stabGnome"]
		}
		return
	case 42:
		r = setup.Reactions["stabDragon"]
		return
	}
	r = setup.Reactions["dontKnowHow"]
	return
}

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
		for _, o := range setup.ObjectsInArea(setup.GetAreaByID(1000)) {
			o.Properties.Area = 29
			setup.GameObjects[o.ID] = o.Properties
		}
		r = setup.Reactions["inTheMoor"]
	}
	return
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

func (obj Object) Use(area setup.Area) (r setup.Reaction) {
	r.KO = false
	r.Sleep = 2
	switch obj.ID {

	case 13:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
		} else {
			obj.Properties.Area = 2000
			setup.GameObjects[obj.ID] = obj.Properties
			r = setup.Reactions["useHood"]
		}
	case 31:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
		} else {
			obj.Properties.Area = 2000
			setup.GameObjects[obj.ID] = obj.Properties
			r = setup.Reactions["useShoes"]
		}
	default:
		r = setup.Reactions["dontKnowHow"]
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

func (obj Object) Throw(area setup.Area) (r setup.Reaction) {
	// sphere?
	if obj.ID == 34 {
		// throwing sphere will always lead to loss
		obj.Properties.Area = 0
		setup.GameObjects[obj.ID] = obj.Properties
		gnome := Object(setup.GetObjectByID(36))
		// no gnome today?
		if !gnome.inArea(area) {
			r = setup.Reactions["brokenSphere"]
			return
		}
		r = setup.Reactions["squashed"]
		// gnome vanished
		gnome.Properties.Area = 0
		setup.GameObjects[gnome.ID] = gnome.Properties
		// golden sphere appears
		goldenSphere := Object(setup.GetObjectByID(45))
		goldenSphere.Properties.Area = area.ID
		setup.GameObjects[goldenSphere.ID] = goldenSphere.Properties
		return
	}
	// on the tree trhowing stone?
	if obj.ID == 20 && area.ID == 31 {
		m := Object(setup.GetObjectByID(46))
		// Map here?
		if !m.inArea(area) {
			r = setup.Reactions["throw"]
			obj.Properties.Area = 9
			setup.GameObjects[obj.ID] = obj.Properties
		}
		if !setup.Flags["MapMissed"] {
			r = setup.Reactions["missMap"]
			setup.Flags["MapMissed"] = true
			// stone falls to ground
			obj.Properties.Area = 9
			setup.GameObjects[obj.ID] = obj.Properties
			return
		}
		r = setup.Reactions["hitMap"]
		// stone and map fall to ground
		obj.Properties.Area = 9
		setup.GameObjects[obj.ID] = obj.Properties
		m.Properties.Area = 9
		setup.GameObjects[m.ID] = m.Properties
		return
	}
	if obj.ID == 45 || obj.ID == 20 {
		r = setup.Reactions["throw"]
		obj.Properties.Area = area.ID
		setup.GameObjects[obj.ID] = obj.Properties
	}
	return
}
