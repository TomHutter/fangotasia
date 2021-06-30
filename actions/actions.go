package actions

import (
	"fangotasia/grid"
	"fangotasia/intro"
	"fangotasia/movement"
	"fangotasia/setup"
	"fangotasia/view"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (obj *Object) NewAreaID(areaID int) {
	obj.Properties.Area = areaID
	setup.GameObjects[obj.ID] = obj.Properties
}

func (obj *Object) NewCondition(condition map[string]string) {
	for lang := range condition {
		desc := obj.Properties.Description[lang]
		desc.Long = condition[lang]
		setup.GameObjects[obj.ID].Description[lang] = desc
		obj.Properties.Description[lang] = desc
	}
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
	// try to pick vanished hood?
	if object.ID == 13 && setup.Flags["HoodVanished"] {
		r = setup.Reactions["dontSee"]
		return
	}
	if !object.available(area) {
		r = setup.Reactions["dontSee"]
		return
	}
	if object.inInventory() || object.inUse() {
		r = setup.Reactions["haveAlready"]
		return
	}

	grub := Object(setup.GetObjectByID(10))
	if grub.inArea(area) {
		switch object.ID {
		case 10:
			r = setup.Reactions["silly"]
			return
			// let's pick stone
		case 20:
			return Object(setup.GetObjectByID(20)).pick()
		default:
			return object.snatchFrom(grub)
		}
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
		r.Statement[setup.Language][0] = fmt.Sprintf(r.Statement[setup.Language][0],
			"[#ff69b4::b]<IMKE>[green:black:-]",
			"[#ff69b4::b]<IMKE>[green:black:-]")
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
		hood := Object(setup.GetObjectByID(13))
		if hood.inUse() {
			r = setup.Reactions["stabDwarfHooded"]
			setup.Flags["HoodVanished"] = true
			hood.NewAreaID(object.Properties.Area)
		} else {
			r = setup.Reactions["stabDwarf"]
		}
		return
	case 36:
		hood := Object(setup.GetObjectByID(13))
		if hood.inUse() {
			r = setup.Reactions["stabGnomeHooded"]
			setup.Flags["HoodVanished"] = true
			hood.NewAreaID(object.Properties.Area)
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
	switch obj.ID {
	case 10, 14, 16, 18, 21, 22, 27, 29, 32, 36, 40, 45, 42, 43:
		r = setup.Reactions["throwFixedObject"]
		return
	}

	if !obj.inInventory() {
		r = setup.Reactions["dontHave"]
		return
	}

	// on the tree throwing stone?
	if obj.ID == 20 && area.ID == 31 {
		m := Object(setup.GetObjectByID(47))
		// Map here?
		if !m.inArea(area) {
			r = setup.Reactions["throwObject"]
			rand.Seed(time.Now().UnixNano())
			obj.NewAreaID(rand.Intn(51))
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

	switch obj.ID {
	// sphere?
	case 34:
		// throwing sphere will always lead to loss
		obj.NewAreaID(0)
		gnome := Object(setup.GetObjectByID(36))
		// no gnome today?
		if !gnome.inArea(area) {
			r = setup.Reactions["brokenSphere"]
			return
		}
		hood := Object(setup.GetObjectByID(13))
		if hood.inArea(area) {
			setup.Flags["HoodVanished"] = false
		}
		r = setup.Reactions["squashed"]
		// gnome vanished
		gnome.NewAreaID(0)
		// golden sphere appears
		goldenSphere := Object(setup.GetObjectByID(46))
		goldenSphere.NewAreaID(area.ID)
		return
	// Fanto Tango
	case 49:
		for _, i := range []int{10, 14, 16, 18, 21, 29, 32, 36, 40, 42, 43, 45} {
			object := Object(setup.GetObjectByID(i))
			if object.inArea(area) {
				if strings.Contains(object.Properties.Description[setup.Language].Long, "::") {
					r = setup.Reactions["alreadyFangoed"]
					return
				} else {
					var cond map[string]string
					cond = make(map[string]string)
					for lang := range setup.Conditions["fango"] {
						article := object.Properties.Description[lang].Article
						parts := strings.Split(object.Properties.Description[lang].Long, " ")[1:]
						long := strings.Join(parts, " ")
						cond[lang] = fmt.Sprintf(setup.Conditions["fango"][lang][article], long)
					}
					object.NewCondition(cond)
					r = setup.Reactions["hitWithFango"]
					return
				}
			}
		}
		r = setup.Reactions["throwFango"]
	default:
		re := regexp.MustCompile(`\\n.*$`)
		r = setup.Reactions["throwObject"]
		r.Statement = make(map[string][]string, len(setup.Reactions["throwObject"].Statement))
		rand.Seed(time.Now().UnixNano())
		newAreaID := rand.Intn(51)
		obj.NewAreaID(newAreaID)
		for lang := range setup.Reactions["throwObject"].Statement {
			r.Statement[lang] = make([]string, len(setup.Reactions["throwObject"].Statement[lang]))
			copy(r.Statement[lang], setup.Reactions["throwObject"].Statement[lang])
			newArea := setup.GetAreaByID(newAreaID)
			article := strings.Title(obj.Properties.Description[lang].Article)
			short := obj.Properties.Description[lang].Short
			location := newArea.Properties.Description[lang].Long
			location = string(re.ReplaceAll([]byte(location), []byte(" ")))
			for i, s := range setup.Reactions["throwObject"].Statement[lang] {
				r.Statement[lang][i] = fmt.Sprintf(s, article, short, location)
			}
		}
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
		r.Statement[setup.Language][0] = view.Highlight(r.Statement[setup.Language][0], "[green:black:-]")
	default:
		r = setup.Reactions["dontKnowHow"]
	}
	return
}

func (obj Object) Say(area setup.Area, word string) (r setup.Reaction) {
	switch strings.ToLower(word) {
	case "simsalabim":
		dwarf := Object(setup.GetObjectByID(18))
		if dwarf.inArea(area) {
			hood := Object(setup.GetObjectByID(13))
			if hood.inArea(area) {
				setup.Flags["HoodVanished"] = false
			}
			dwarf.NewAreaID(0)
			r = setup.Reactions["simsalabim"]
			return
		}
	case "fangotasia":
		if area.ID == 1 {
			scoreBoard(false, false)
		}
	}
	r = setup.GetReactionByName("say")
	r.Statement[setup.Language][0] = fmt.Sprintf(r.Statement[setup.Language][0], word)
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
		a := strings.Title(obj.Properties.Description[setup.Language].Article)
		desc := fmt.Sprintf("%s %s", a, obj.Properties.Description[setup.Language].Short)
		r.Statement[setup.Language][0] = fmt.Sprintf(r.Statement[setup.Language][0], desc)
		return
	default:
		r = setup.GetReactionByName("unusable")
		r.Statement[setup.Language][0] = fmt.Sprintf(r.Statement[setup.Language][0],
			view.Highlight(obj.Properties.Description[setup.Language].Long, "[red]"))
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
		a := strings.Title(obj.Properties.Description[setup.Language].Article)
		desc := fmt.Sprintf("%s %s", a, obj.Properties.Description[setup.Language].Short)
		r.Statement[setup.Language][0] = fmt.Sprintf(setup.Reactions["feed"].Statement[setup.Language][0], desc)
	case 16:
		berries := Object(setup.GetObjectByID(23))
		if berries.inInventory() {
			hood := Object(setup.GetObjectByID(13))
			if hood.inArea(area) {
				setup.Flags["HoodVanished"] = false
			}
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
			hood := Object(setup.GetObjectByID(13))
			if hood.inArea(setup.GetAreaByID(29)) {
				setup.Flags["HoodVanished"] = false
			}
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
	case 10, 14, 16, 18, 21, 22, 27, 29, 32, 36, 40, 45, 42, 43:
		if !obj.inArea(area) {
			r = setup.Reactions["dontSee"]
			return
		}
		r = setup.Reactions["cantEat"]
	default:
		if !obj.inInventory() {
			r = setup.Reactions["dontHave"]
			return
		}
		r = setup.Reactions["cantEat"]
	}
	return
}

func (jar Object) Drink(area setup.Area) (r setup.Reaction) {
	if !jar.inInventory() {
		r = setup.Reactions["noTool"]
		return
	}
	if jar.Properties.Description[setup.Language].Long ==
		setup.Conditions["jar"]["empty"][setup.Language] {
		r = setup.Reactions["jarEmpty"]
	}
	if jar.Properties.Description[setup.Language].Long ==
		setup.Conditions["jar"]["full"][setup.Language] {
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
		if obj.Properties.Description[setup.Language].Long ==
			setup.Conditions["ring"]["golden"][setup.Language] {
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
		if jar.Properties.Description[setup.Language].Long ==
			setup.Conditions["jar"]["empty"][setup.Language] {
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
		if obj.Properties.Description[setup.Language].Long ==
			setup.Conditions["bush"]["watered"][setup.Language] {
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
		a := strings.Title(obj.Properties.Description[setup.Language].Article)
		desc := fmt.Sprintf("%s %s", a, obj.Properties.Description[setup.Language].Short)
		r.Statement[setup.Language][0] = fmt.Sprintf(setup.Reactions["scare"].Statement[setup.Language][0], desc)
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
		SetLabel(fmt.Sprintf("%s \u23CE ", setup.TextElements["next"][setup.Language])).
		SetAcceptanceFunc(tview.InputFieldMaxLength(0)).
		SetDoneFunc(func(key tcell.Key) {
			grid.Grid.Clear()
			grid.Grid.AddItem(grid.InputGrid, 0, 0, 1, 1, 0, 0, false)
			grid.App.SetFocus(grid.InputField)
		})
	grid.App.SetFocus(grid.AreaField)
	r = setup.Reactions["ok"]
	return
}

func (obj Object) Lang(area setup.Area) (r setup.Reaction) {
	var keys []string
	keys = make([]string, 0)

	for k := range setup.TextElements["selectLanguage"] {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	grid.InputField.SetText("")
	grid.Response.SetText("")
	grid.Grid.Clear()
	grid.Grid.AddItem(grid.LanguageGrid, 0, 0, 1, 1, 0, 0, false)
	grid.LanguageSelect.
		SetLabel(setup.TextElements["selectLanguage"][setup.Language]).
		SetOptions(keys, nil).
		SetSelectedFunc(func(o string, i int) {
			setup.Language = o
			grid.Surroundings.SetText(strings.Join(view.Surroundings(area), "\n"))
			grid.Grid.Clear()
			grid.Grid.AddItem(grid.InputGrid, 0, 0, 1, 1, 0, 0, false)
			grid.App.SetFocus(grid.InputField)
		})
	grid.App.SetFocus(grid.LanguageSelect)
	r = setup.Reactions["ok"]
	return
}

func (obj Object) Help(area setup.Area) (r setup.Reaction) {
	grid.Grid.Clear()
	grid.Grid.AddItem(grid.AreaGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AreaMap.SetText("")
	grid.AreaField.SetText("")
	grid.App.SetFocus(grid.AreaField)
	intro.Intro()
	r = setup.Reactions["ok"]
	return
}
