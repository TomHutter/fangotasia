package actions

import (
	"fangotasia/movement"
	"fangotasia/setup"
	"math/rand"
	"time"
)

func beads(areaID int) {
	if areaID == 1 {
		setup.Beads = 1
	} else {
		setup.Beads += areaID
	}

	if setup.Beads == 108 {
		imke := Object(setup.GetObjectByID(48))
		imke.Properties.Description.Long = "[#ff69b4::b]<IMKE>[blue:black:-] den pink Diamanten"
		//SetFieldTextColor(tcell.ColorHotPink).
		imke.Properties.Description.Short = "Imke"
		imke.NewAreaID(areaID)
	}
}

// As Move is called in context of object handling, Move reflects on Object even obj is not used.
func (obj Object) Move(area setup.Area, dir string) (r setup.Reaction, areaID int) {
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
	}

	beads(newArea)

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
		beads(area.ID)
		moves += 1
		r.OK = true
		areaID = 9
		return
	}
	if area.ID == 9 && object.ID == 27 {
		beads(area.ID)
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

func (object Object) Jump(area setup.Area) (r setup.Reaction, areaID int) {
	if area.ID == 31 {
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(4) > 0 {
			r = setup.Reactions["jumpTree"]
			return
		} else {
			beads(area.ID)
			moves += 1
			r.OK = true
			areaID = 9
			return
		}
	}
	r = setup.Reactions["jump"]
	return
}
