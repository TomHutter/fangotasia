package actions

import (
	"fangotasia/grid"
	"fangotasia/intro"
	"fangotasia/movement"
	"fangotasia/setup"
	"fangotasia/view"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (obj *Object) NewAreaID(areaID int) {
	obj.Properties.Area = areaID
	setup.GameObjects[obj.ID] = obj.Properties
}

func (obj *Object) NewCondition(condition string) {
	obj.Properties.Description.Long = condition
	setup.GameObjects[obj.ID] = obj.Properties
}

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
		letter := Object(setup.GetObjectByID(38))
		letter.NewAreaID(area.ID)
		ruby := Object(setup.GetObjectByID(39))
		ruby.NewAreaID(area.ID)
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
	case 10, 16, 18, 21, 22, 27, 36, 40, 42, 45:
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
	case 48:
		r = setup.GetReactionByName("takeImke")
		r.Statement[0] = fmt.Sprintf(r.Statement[0], "[#ff69b4::b]<IMKE>[green:black:-]")
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
			dwarf := Object(setup.GetObjectByID(18))
			dwarf.NewAreaID(0)
			r = setup.Reactions["stabDwarfHooded"]
		} else {
			r = setup.Reactions["stabDwarf"]
		}
		return
	case 36:
		if Object(setup.GetObjectByID(13)).inUse() {
			gnome := Object(setup.GetObjectByID(36))
			gnome.NewAreaID(0)
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

	obj.NewAreaID(setup.INUSE)

	switch obj.ID {
	case 13:
		r = setup.Reactions["useHood"]
	case 31:
		r = setup.Reactions["useShoes"]
	case 47:
		r = setup.Reactions["useMap"]
		setup.Verbs = setup.AddMapVerb(setup.Verbs)
	}
	return
}

func (obj Object) Throw(area setup.Area) (r setup.Reaction) {
	if !obj.inInventory() {
		r = setup.Reactions["dontHave"]
		return
	}
	// sphere?
	if obj.ID == 34 {
		// throwing sphere will always lead to loss
		obj.NewAreaID(0)
		gnome := Object(setup.GetObjectByID(36))
		// no gnome today?
		if !gnome.inArea(area) {
			r = setup.Reactions["brokenSphere"]
			return
		}
		r = setup.Reactions["squashed"]
		// gnome vanished
		gnome.NewAreaID(0)
		// golden sphere appears
		goldenSphere := Object(setup.GetObjectByID(46))
		goldenSphere.NewAreaID(area.ID)
		return
	}
	// on the tree throwing stone?
	if obj.ID == 20 && area.ID == 31 {
		m := Object(setup.GetObjectByID(47))
		// Map here?
		if !m.inArea(area) {
			r = setup.Reactions["throwStone"]
			obj.NewAreaID(9)
		}
		if !setup.Flags["MapMissed"] {
			r = setup.Reactions["missMap"]
			setup.Flags["MapMissed"] = true
			// stone falls to ground
			obj.NewAreaID(9)
			return
		}
		r = setup.Reactions["hitMap"]
		// stone and map fall to ground
		obj.NewAreaID(9)
		m.NewAreaID(9)
		return
	}
	if obj.ID == 46 || obj.ID == 20 {
		r = setup.Reactions["throwStone"]
		obj.NewAreaID(area.ID)
	}
	if obj.ID == 49 {
		for _, i := range []int{10, 14, 16, 18, 21, 29, 32, 36, 40, 42, 43, 45} {
			object := Object(setup.GetObjectByID(i))
			if object.inArea(area) {
				if strings.Contains(object.Properties.Description.Long, "::") {
					r = setup.Reactions["alreadyFangoed"]
					return
				} else {
					article := object.Properties.Description.Article
					parts := strings.Split(object.Properties.Description.Long, " ")[1:]
					long := strings.Join(parts, " ")
					newLong := fmt.Sprintf(setup.Conditions["fango"][article], long)
					object.NewCondition(newLong)
					r = setup.Reactions["hitWithFango"]
					return
				}
			}
		}
		r = setup.Reactions["throwFango"]
	}
	return
}

func (obj Object) Look(area setup.Area) (r setup.Reaction) {
	switch obj.ID {
	case 21:
		r = setup.Reactions["wallpainting"]
	default:
		r = setup.Reactions["looksGood"]
	}
	return
}

func (obj Object) Read(area setup.Area) (r setup.Reaction) {
	var reaction = map[int]string{
		11: "book",
		12: "paper",
		17: "shield",
		21: "wallpainting",
		28: "parchment",
		32: "panel",
		43: "rottenPanel",
		38: "letter",
		47: "readMap",
	}
	switch obj.ID {
	case 11, 12, 17, 28, 38, 47:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
			return
		}
	}
	switch obj.ID {
	case 11, 12, 17, 21, 28, 38, 43, 47:
		r = setup.Reactions[reaction[obj.ID]]
	case 32:
		r = setup.GetReactionByName(reaction[obj.ID])
		r.Statement[0] = view.Highlight(r.Statement[0], "[green:black:-]")
	default:
		r = setup.Reactions["dontKnowHow"]
	}
	return
}

func (obj Object) Say(area setup.Area, word string) (r setup.Reaction) {
	switch word {
	case "simsalabim":
		dwarf := Object(setup.GetObjectByID(18))
		if dwarf.inArea(area) {
			dwarf.NewAreaID(0)
			r = setup.Reactions["simsalabim"]
			return
		}
	case "fangotasia":
		if area.ID == 1 {
			var points int
			for _, o := range setup.ObjectsInArea(area) {
				points = points + int(o.Properties.Value)
			}
			r = setup.GetReactionByName("fangotasia")
			r.Statement[0] = fmt.Sprintf(r.Statement[0], points)
			return
		}
	}
	r = setup.GetReactionByName("say")
	r.Statement[0] = fmt.Sprintf(r.Statement[0], word)
	return
}

func (obj Object) Put(area setup.Area) (r setup.Reaction) {
	if !obj.inInventory() {
		r = setup.Reactions["dontHave"]
		return
	}
	obj.NewAreaID(area.ID)
	r = setup.Reactions["ok"]
	return
}

func (obj Object) Fill(area setup.Area) (r setup.Reaction) {
	if !obj.inInventory() {
		r = setup.Reactions["dontHave"]
		return
	}
	switch obj.ID {
	case 30:
		r = setup.Reactions["ok"]
	case 35, 44:
		r = setup.GetReactionByName("unsuitable")
		a := strings.Title(obj.Properties.Description.Article)
		desc := fmt.Sprintf("%s %s", a, obj.Properties.Description.Short)
		r.Statement[0] = fmt.Sprintf(r.Statement[0], desc)
		return
	default:
		r = setup.GetReactionByName("unusable")
		r.Statement[0] = fmt.Sprintf(r.Statement[0], view.Highlight(obj.Properties.Description.Long, "[red]"))
		return
	}

	switch area.ID {
	case 3, 35:
		r = setup.Reactions["waterUnreachable"]
	case 17:
		r = setup.Reactions["ok"]
		obj.NewCondition(setup.Conditions["jar"]["full"])
	default:
		r = setup.Reactions["noWater"]
	}
	return
}

func (obj Object) Feed(area setup.Area) (r setup.Reaction) {
	switch obj.ID {
	case 10, 18, 36, 24:
		r = setup.GetReactionByName("feed")
		a := strings.Title(obj.Properties.Description.Article)
		desc := fmt.Sprintf("%s %s", a, obj.Properties.Description.Short)
		r.Statement[0] = fmt.Sprintf(setup.Reactions["feed"].Statement[0], desc)
	case 16:
		berries := Object(setup.GetObjectByID(23))
		if berries.inInventory() {
			r = setup.Reactions["feedBaerWithBerries"]
			obj.NewAreaID(0)
			berries.NewAreaID(0)
		} else {
			r = setup.Reactions["feedBaer"]
		}
	default:
		r = setup.Reactions["dontKnowHow"]
	}
	return
}

func (obj Object) Cut(area setup.Area) (r setup.Reaction) {
	if obj.ID != 14 {
		r = setup.Reactions["dontKnowHow"]
		return
	}
	ring := Object(setup.GetObjectByID(37))
	// ring already present?
	if ring.Properties.Area != 0 {
		r = setup.Reactions["fruitAlreadyCut"]
		return
	}
	dagger := Object(setup.GetObjectByID(33))
	if !dagger.inInventory() {
		r = setup.Reactions["noTool"]
		return
	}
	r = setup.Reactions["cutFruit"]
	ring.NewAreaID(area.ID)
	return
}

func (obj Object) Catapult(area setup.Area) (r setup.Reaction) {
	if !obj.inInventory() {
		r = setup.Reactions["dontHave"]
		return
	}
	switch obj.ID {
	case 34:
		r = setup.Reactions["tryThrow"]
	case 20:
		// Catapult around?
		catapult := Object(setup.GetObjectByID(29))
		if !catapult.inArea(area) {
			r = setup.Reactions["noTool"]
		} else {
			obj.NewAreaID(29)
			grub := Object(setup.GetObjectByID(10))
			grub.NewAreaID(0)
			r = setup.Reactions["ok"]
		}
	default:
		r = setup.Reactions["dontKnowHow"]
	}
	return
}

func (obj Object) Eat(area setup.Area) (r setup.Reaction) {
	switch obj.ID {
	case 9:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
			return
		}
		obj.NewAreaID(setup.INUSE)
		r = setup.Reactions["eatCake"]
	case 23:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
			return
		}
		obj.NewAreaID(0)
		r = setup.Reactions["ok"]
	default:
		r = setup.Reactions["cantEat"]
	}
	return
}

func (jar Object) Drink(area setup.Area) (r setup.Reaction) {
	if !jar.inInventory() {
		r = setup.Reactions["noTool"]
		return
	}
	if jar.Properties.Description.Long == setup.Conditions["jar"]["empty"] {
		r = setup.Reactions["jarEmpty"]
	}
	if jar.Properties.Description.Long == setup.Conditions["jar"]["full"] {
		jar.NewCondition(setup.Conditions["jar"]["empty"])
		r = setup.Reactions["drinkJar"]
	}
	return
}

func (obj Object) Spin(area setup.Area) (r setup.Reaction) {
	switch obj.ID {
	case 10, 14, 16, 18, 21, 22, 27, 29, 32, 36, 40, 45, 42, 43:
		r = setup.Reactions["dontSpin"]
	case 37:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
			return
		}
		if obj.Properties.Description.Long == setup.Conditions["ring"]["golden"] {
			r = setup.Reactions["spin"]
		} else {
			obj.NewCondition(setup.Conditions["ring"]["golden"])
			obj.NewAreaID(area.ID)
			r = setup.Reactions["spinRing"]
		}
	default:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
			return
		}
		r = setup.Reactions["spin"]
	}
	return
}

func (obj Object) Water(area setup.Area) (r setup.Reaction) {
	switch obj.ID {
	case 14, 22, 27:
		jar := Object(setup.GetObjectByID(30))
		if !jar.inInventory() {
			r = setup.Reactions["noJar"]
			return
		}
		if jar.Properties.Description.Long == setup.Conditions["jar"]["empty"] {
			r = setup.Reactions["jarEmpty"]
			return
		}
	default:
		r = setup.Reactions["silly"]
		return
	}

	switch obj.ID {
	case 14, 27:
		r = setup.Reactions["ok"]
	case 22:
		jar := Object(setup.GetObjectByID(30))
		jar.NewCondition(setup.Conditions["jar"]["empty"])
		if obj.Properties.Description.Long == setup.Conditions["bush"]["watered"] {
			r = setup.Reactions["ok"]
		} else {
			r = setup.Reactions["waterBush"]
			obj.NewCondition(setup.Conditions["bush"]["watered"])
			berries := Object(setup.GetObjectByID(23))
			berries.NewAreaID(area.ID)
			leaves := Object(setup.GetObjectByID(24))
			leaves.NewAreaID(area.ID)
		}
	}
	return
}

func (obj Object) Scare(area setup.Area) (r setup.Reaction) {
	switch obj.ID {
	case 10, 16, 18, 36, 42:
		r = setup.GetReactionByName("scare")
		a := strings.Title(obj.Properties.Description.Article)
		desc := fmt.Sprintf("%s %s", a, obj.Properties.Description.Short)
		r.Statement[0] = fmt.Sprintf(setup.Reactions["scare"].Statement[0], desc)
	default:
		r = setup.Reactions["silly"]
	}
	return
}

func (obj Object) Map(area setup.Area) (r setup.Reaction) {
	grid.InputField.SetText("")
	grid.AreaMap.SetText(strings.Join(movement.DrawMap(area), "\n"))
	grid.Grid.Clear()
	grid.Grid.AddItem(grid.AreaGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AreaField.SetText("").
		SetLabel("Weiter \u23CE ").
		SetAcceptanceFunc(tview.InputFieldMaxLength(0)).
		SetDoneFunc(func(key tcell.Key) {
			grid.Grid.Clear()
			grid.Grid.AddItem(grid.InputGrid, 0, 0, 1, 1, 0, 0, false)
			grid.App.SetFocus(grid.InputField)
		})
	grid.App.SetFocus(grid.AreaField)
	return
}

func (obj Object) Help(area setup.Area) (r setup.Reaction) {
	grid.Grid.Clear()
	grid.Grid.AddItem(grid.AreaGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AreaMap.SetText("")
	grid.AreaField.SetText("")
	grid.App.SetFocus(grid.AreaField)
	intro.Intro()
	return
}
