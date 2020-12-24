package actions

import (
	"fangotasia/movement"
	"fangotasia/setup"
	"fmt"
	"math/rand"
	"strings"
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
		for lang := range imke.Properties.Description {
			desc := imke.Properties.Description[lang]
			desc.Long = fmt.Sprintf("[#ff69b4::b]<IMKE>[blue:black:-] %s",
				setup.TextElements["imke"][lang])
			desc.Short = "Imke"
			setup.GameObjects[imke.ID].Description[lang] = desc
		}
		imke.NewAreaID(areaID)
	}
}

// As Move is called in context of object handling, Move reflects on Object even obj is not used.
func (obj Object) Move(area setup.Area, dir string) (r setup.Reaction, areaID int) {
	var char string
	var direction = map[string]int{}
	char = strings.ToLower(string(setup.TextElements["north"][setup.Language][0]))
	direction[char] = 0
	char = strings.ToLower(string(setup.TextElements["south"][setup.Language][0]))
	direction[char] = 1
	char = strings.ToLower(string(setup.TextElements["east"][setup.Language][0]))
	direction[char] = 2
	char = strings.ToLower(string(setup.TextElements["west"][setup.Language][0]))
	direction[char] = 3

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

	if area.ID == 29 {
		setup.Flags["Castle"] = true
	}
	// Area 30 and 25 are connected by a door. Is it open?
	if (area.ID == 30 || area.ID == 25 && direction[dir] == 0) && !setup.Flags["DoorOpen"] {
		r = setup.Reactions["locked"]
		areaID = area.ID
		return
	}
	movement.RevealArea(newArea)
	setup.Moves += 1
	r = setup.Reactions["ok"]
	areaID = newArea
	// Direction swamp?
	if newArea == 5 {
		setup.Flags["Swamp"] = true
		for _, o := range setup.ObjectsInArea(setup.GetAreaByID(setup.INVENTORY)) {
			obj := Object(o)
			obj.NewAreaID(29)
		}
		r = setup.Reactions["inTheSwamp"]
	}
	return
}

func (object Object) Climb(area setup.Area) (r setup.Reaction, areaID int) {
	if area.ID == 31 {
		beads(area.ID)
		setup.Moves += 1
		r = setup.Reactions["ok"]
		areaID = 9
		return
	}
	if area.ID == 9 && object.ID == 27 {
		beads(area.ID)
		movement.RevealArea(31)
		setup.Moves += 1
		setup.Flags["Tree"] = true
		r = setup.Reactions["ok"]
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
			setup.Moves += 1
			r = setup.Reactions["ok"]
			areaID = 9
			return
		}
	}
	r = setup.Reactions["jump"]
	areaID = area.ID
	return
}
