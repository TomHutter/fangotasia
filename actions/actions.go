package actions

import (
	"fantasia/setup"
	"fantasia/view"
)

func (object Object) Open(area setup.Area) (r setup.Reaction) {
	if !object.available(area) {
		r = setup.Reactions["dontSee"]
		return
	}

	switch object.ID {
	case 35, 40, 45:
		// key present?
		key := Object(setup.GetObjectByID(26))
		if !key.inInventory() {
			r = setup.Reactions["noKey"]
			return
		}
	default:
		r = setup.Reactions["dontKnowHow"]
		return
	}

	switch object.ID {
	// red box
	case 35:
		if setup.Flags["BoxOpen"] {
			r = setup.Reactions["alreadyOpen"]
			return
		}
		r = setup.Reactions["openBox"]
		letter := setup.GetObjectByID(38)
		letter.Properties.Area = area.ID
		setup.GameObjects[letter.ID] = letter.Properties
		ruby := setup.GetObjectByID(39)
		ruby.Properties.Area = area.ID
		setup.GameObjects[ruby.ID] = ruby.Properties
		setup.Flags["BoxOpen"] = true

		// door
	case 40, 45:
		if setup.Flags["DoorOpen"] {
			r = setup.Reactions["alreadyOpen"]
			return
		}
		setup.Flags["DoorOpen"] = true
		r = setup.Reactions["ok"]
	}
	/*
	   	446 f=0:gosub605:iffl=1thenfl=0:goto280
	   447 ifno=40thenprint"versuche 'sperre'.":goto280
	   496 f=0:gosub605:iffl=1thenfl=0:goto280
	   497 ifno<>40andno<>35thenprinta$(2):goto280
	   498 ifno=35thenprint"versuche 'oeffne'.":goto280
	   499 iftu=1thenprint"ist schon offen !":goto280
	   500 ifge(26)<>-1thenprint"ich habe keinen schluessel.":goto280
	   501 print"gut.":tu=1:goto281


	   448 f=1:gosub607:iffl=1thenfl=0:goto280
	   449 ifno<>35thenprinta$(2):goto280
	   450 ifge(38)<>0thenprint"gut. es ist leer.":goto281
	   451 ifge(26)<>-1thenprinta$(1):goto280
	   452 print"zwei dinge fallen heraus. sag 'sieh'.":ge(38)=oa:ge(39)=oa:goto281
	*/
	return
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
	case 47:
		if area.ID == 31 {
			r = setup.Reactions["unreachable"]
			return
		}
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
	switch obj.ID {
	case 13, 31, 47:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
			return
		}
	default:
		r = setup.Reactions["dontKnowHow"]
		return
	}

	obj.Properties.Area = setup.INUSE
	setup.GameObjects[obj.ID] = obj.Properties

	switch obj.ID {
	case 13:
		r = setup.Reactions["useHood"]
	case 31:
		r = setup.Reactions["useShoes"]
	case 47:
		r = setup.Reactions["useMap"]
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
		goldenSphere := Object(setup.GetObjectByID(46))
		goldenSphere.Properties.Area = area.ID
		setup.GameObjects[goldenSphere.ID] = goldenSphere.Properties
		return
	}
	// on the tree trhowing stone?
	if obj.ID == 20 && area.ID == 31 {
		m := Object(setup.GetObjectByID(47))
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
	if obj.ID == 46 || obj.ID == 20 {
		r = setup.Reactions["throw"]
		obj.Properties.Area = area.ID
		setup.GameObjects[obj.ID] = obj.Properties
	}
	return
}

func (obj Object) Read(area setup.Area) (r setup.Reaction) {
	var reaction = map[int]string{
		21: "wallpainting",
		32: "panel",
		43: "rottenPanel",
		12: "paper",
		17: "shield",
		28: "parchment",
		38: "letter",
		47: "readMap",
	}
	switch obj.ID {
	case 12, 17, 28, 38, 47:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
			return
		}
	}
	switch obj.ID {
	case 12, 17, 28, 38, 43, 47:
		r = setup.Reactions[reaction[obj.ID]]
	case 32:
		r = setup.Reactions[reaction[obj.ID]]
		r.Statement = view.Highlight(r.Statement, setup.GREEN)
	default:
		r = setup.Reactions["dontKnowHow"]
	}
	return
}
