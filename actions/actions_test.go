package actions_test

import (
	"fantasia/actions"
	"fantasia/setup"
	"fantasia/view"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	setup.Init()
}

func TestParse(t *testing.T) {
	setup.Init()
	var text = []string{}
	// go to area 1
	area1 := setup.GetAreaByID(1)
	area := actions.Parse("nimm zauberschuhe", area1, text)
	assert.Equal(t, area1, area)
	area = actions.Parse("trage zauberschuhe", area1, text)
	assert.Equal(t, area1, area)
	area = actions.Parse("o", area1, text)
	area2 := setup.GetAreaByID(2)
	assert.Equal(t, area2, area)
}
func TestTake(t *testing.T) {
	setup.Init()
	// go to area 1
	area := setup.GetAreaByID(1)
	// pick up dwarf
	obj := setup.GetObjectByID(18)
	res := actions.Object(obj).Take(area)
	// dwarf not in area
	assert.Equal(t, setup.Reactions["dontSee"].Statement, res.Statement)
	assert.False(t, res.OK)

	// pick magic shoes
	obj = setup.GetObjectByID(31)
	res = actions.Object(obj).Take(area)
	// reload obj to get new area => should be setup.INVENTORY by taken into inv
	obj = setup.GetObjectByID(31)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)
	// look up inventory
	inv := setup.ObjectsInArea(setup.GetAreaByID(setup.INVENTORY))
	assert.Contains(t, inv, obj)

	// try to pick shoes again
	res = actions.Object(obj).Take(area)
	assert.Equal(t, setup.Reactions["haveAlready"].Statement, res.Statement)
	assert.True(t, res.OK)

	// to to area 4
	area = setup.GetAreaByID(4)
	// pick up gnome
	obj = setup.GetObjectByID(36)
	res = actions.Object(obj).Take(area)
	assert.Equal(t, setup.Reactions["silly"].Statement, res.Statement)
	assert.False(t, res.OK)
}

func TestStab(t *testing.T) {
	// go to area 4
	area := setup.GetAreaByID(4)
	// stab gnome
	gnome := setup.GetObjectByID(36)
	res := actions.Object(gnome).Stab(area)
	assert.Equal(t, setup.Reactions["noTool"].Statement, res.Statement)
	assert.False(t, res.OK)

	// go to area 3
	area = setup.GetAreaByID(3)
	// pick up dwarf dagger
	obj := setup.GetObjectByID(33)
	res = actions.Object(obj).Take(area)
	// go to area 4
	area = setup.GetAreaByID(4)
	// stab gnome
	res = actions.Object(gnome).Stab(area)
	assert.Equal(t, setup.Reactions["stabGnome"].Statement, res.Statement)
	assert.False(t, res.OK)
	assert.True(t, res.KO)
}

func TestMove(t *testing.T) {
	// go to area 4
	area := setup.GetAreaByID(1)
	obj := setup.Object{}
	// move north
	res, areaID := actions.Object(obj).Move(area, "n")
	assert.Equal(t, setup.Reactions["noWay"].Statement, res.Statement)
	assert.False(t, res.OK)
	assert.Equal(t, areaID, 1)

	// move east
	res, areaID = actions.Object(obj).Move(area, "o")
	assert.Equal(t, setup.Reactions["noShoes"].Statement, res.Statement)
	assert.False(t, res.OK)
	assert.True(t, res.KO)
	assert.Equal(t, areaID, 1)

	// pick magic shoes
	obj = setup.GetObjectByID(31)
	actions.Object(obj).Take(area)
	actions.Object(obj).Use(area)
	res, areaID = actions.Object(obj).Move(area, "o")
	assert.Equal(t, areaID, 2)
	assert.True(t, res.OK)
}

func TestUse(t *testing.T) {
	setup.Init()
	// go to area 1
	area := setup.GetAreaByID(1)
	// pick magic shoes
	shoes := setup.GetObjectByID(31)
	actions.Object(shoes).Take(area)
	// reload obj to get new area => should be setup.INVENTORY (inventory)
	shoes = setup.GetObjectByID(31)
	res := actions.Object(shoes).Use(area)
	// reload obj to get new area => should be setup.INUSE (inUse)
	shoes = setup.GetObjectByID(31)
	assert.Equal(t, setup.Reactions["useShoes"].Statement, res.Statement)
	assert.True(t, res.OK)
	// look up inUse
	inUse := setup.ObjectsInArea(setup.GetAreaByID(setup.INUSE))
	assert.Contains(t, inUse, shoes)

	// go to area 3
	area = setup.GetAreaByID(3)
	// pick magic shoes
	dagger := setup.GetObjectByID(33)
	actions.Object(dagger).Take(area)
	// reload obj to get new area => should be setup.INVENTORY (inventory)
	dagger = setup.GetObjectByID(33)
	res = actions.Object(dagger).Use(area)
	// reload obj to get new area => should be setup.INUSE (inUse)
	dagger = setup.GetObjectByID(33)
	assert.Equal(t, setup.Reactions["dontKnowHow"].Statement, res.Statement)
	assert.False(t, res.OK)

	// use map
	m := setup.GetObjectByID(47)
	m.Properties.Area = setup.INVENTORY
	setup.GameObjects[m.ID] = m.Properties
	res = actions.Object(m).Use(area)
	assert.Equal(t, setup.Reactions["useMap"].Statement, res.Statement)
	m = setup.GetObjectByID(47)
	assert.True(t, res.OK)
	assert.Equal(t, setup.INUSE, m.Properties.Area)
}

func TestClimb(t *testing.T) {
	// go to area 9
	area := setup.GetAreaByID(9)
	tree := setup.GetObjectByID(27)
	res, areaID := actions.Object(tree).Climb(area)
	assert.Equal(t, areaID, 31)
	assert.True(t, res.OK)

	// got to area 31
	area = setup.GetAreaByID(31)
	tree = setup.GetObjectByID(27)
	res, areaID = actions.Object(tree).Climb(area)
	assert.Equal(t, areaID, 9)
	assert.True(t, res.OK)

	// got to area 1
	area = setup.GetAreaByID(1)
	tree = setup.GetObjectByID(32)
	res, areaID = actions.Object(tree).Climb(area)
	assert.Equal(t, setup.Reactions["silly"].Statement, res.Statement)
	assert.Equal(t, areaID, 1)
	assert.False(t, res.OK)
}

func TestThrow(t *testing.T) {
	// not area 4 - sphere breaks
	area := setup.GetAreaByID(1)
	sphere := setup.GetObjectByID(34)
	res := actions.Object(sphere).Throw(area)
	assert.Equal(t, setup.Reactions["brokenSphere"].Statement, res.Statement)
	assert.Equal(t, 0, setup.GetObjectByID(34).Properties.Area)
	assert.False(t, res.OK)

	// area 4 - gnome breaks
	area = setup.GetAreaByID(4)
	res = actions.Object(sphere).Throw(area)
	assert.Equal(t, setup.Reactions["squashed"].Statement, res.Statement)
	assert.Equal(t, 0, setup.GetObjectByID(36).Properties.Area)
	assert.True(t, res.OK)

	// golden sphere
	area = setup.GetAreaByID(1)
	sphere = setup.GetObjectByID(46)
	res = actions.Object(sphere).Throw(area)
	assert.Equal(t, setup.Reactions["throw"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, 1, setup.GetObjectByID(46).Properties.Area)

	// stone
	area = setup.GetAreaByID(9)
	stone := setup.GetObjectByID(20)
	res = actions.Object(stone).Throw(area)
	assert.Equal(t, setup.Reactions["throw"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, 9, setup.GetObjectByID(20).Properties.Area)

	// stone on tree - map present by setup
	area = setup.GetAreaByID(31)
	assert.False(t, setup.Flags["MapMissed"])
	res = actions.Object(stone).Throw(area)
	assert.Equal(t, setup.Reactions["missMap"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.True(t, setup.Flags["MapMissed"])
	assert.Equal(t, 9, setup.GetObjectByID(20).Properties.Area)

	// stone on tree - map present by setup, 2nd try
	res = actions.Object(stone).Throw(area)
	assert.Equal(t, setup.Reactions["hitMap"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, 9, setup.GetObjectByID(20).Properties.Area)
	assert.Equal(t, 9, setup.GetObjectByID(47).Properties.Area)
}

func TestRead(t *testing.T) {
	// read panel
	area := setup.GetAreaByID(1)
	panel := setup.GetObjectByID(32)
	res := actions.Object(panel).Read(area)
	assert.Equal(t, view.Highlight(setup.Reactions["panel"].Statement, setup.GREEN), res.Statement)
	assert.True(t, res.OK)

	// read letter
	letter := setup.GetObjectByID(12)
	res = actions.Object(letter).Read(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	// read letter in inv
	letter.Properties.Area = setup.INVENTORY
	res = actions.Object(letter).Read(area)
	assert.Equal(t, setup.Reactions["paper"].Statement, res.Statement)
	assert.True(t, res.OK)
}

func TestOpen(t *testing.T) {
	// open panel
	area := setup.GetAreaByID(1)
	panel := setup.GetObjectByID(32)
	res := actions.Object(panel).Open(area)
	assert.Equal(t, setup.Reactions["dontKnowHow"].Statement, res.Statement)
	assert.False(t, res.OK)

	// open red box
	box := setup.GetObjectByID(35)
	res = actions.Object(box).Open(area)
	assert.Equal(t, setup.Reactions["dontSee"].Statement, res.Statement)
	assert.False(t, res.OK)

	// correct area
	area = setup.GetAreaByID(4)
	res = actions.Object(box).Open(area)
	assert.Equal(t, setup.Reactions["noKey"].Statement, res.Statement)
	assert.False(t, res.OK)
	assert.False(t, setup.Flags["BoxOpen"])

	// key in inv
	key := setup.GetObjectByID(26)
	key.Properties.Area = setup.INVENTORY
	setup.GameObjects[key.ID] = key.Properties
	res = actions.Object(box).Open(area)
	assert.Equal(t, setup.Reactions["openBox"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.True(t, setup.Flags["BoxOpen"])

	// reopen box
	res = actions.Object(box).Open(area)
	assert.Equal(t, setup.Reactions["alreadyOpen"].Statement, res.Statement)
	assert.True(t, res.OK)

	// open door
	area = setup.GetAreaByID(25)
	door := setup.GetObjectByID(40)
	key.Properties.Area = 9
	setup.GameObjects[key.ID] = key.Properties
	res = actions.Object(door).Open(area)
	assert.Equal(t, setup.Reactions["noKey"].Statement, res.Statement)
	assert.False(t, res.OK)
	assert.False(t, setup.Flags["DoorOpen"])

	// key in inv
	key.Properties.Area = setup.INVENTORY
	setup.GameObjects[key.ID] = key.Properties
	res = actions.Object(door).Open(area)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.True(t, setup.Flags["DoorOpen"])

	// reopen door (from the other side)
	area = setup.GetAreaByID(30)
	door = setup.GetObjectByID(45)
	res = actions.Object(door).Open(area)
	assert.Equal(t, setup.Reactions["alreadyOpen"].Statement, res.Statement)
	assert.True(t, res.OK)
}

func TestPut(t *testing.T) {
	setup.Init()
	area := setup.GetAreaByID(1)
	sword := setup.GetObjectByID(15)
	dagger := setup.GetObjectByID(33)

	// put sword into inv
	sword.Properties.Area = setup.INVENTORY
	setup.GameObjects[sword.ID] = sword.Properties

	// try to put dagger
	res := actions.Object(dagger).Put(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	// to put dagger
	res = actions.Object(sword).Put(area)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, 1, setup.GetObjectByID(15).Properties.Area)
}

func TestSay(t *testing.T) {
	setup.Init()
	area := setup.GetAreaByID(1)
	sword := setup.GetObjectByID(15)
	dagger := setup.GetObjectByID(33)
	obj := actions.Object{}

	// put sword in area 1
	sword.Properties.Area = 1
	setup.GameObjects[sword.ID] = sword.Properties

	// put dagger in area 1
	dagger.Properties.Area = 1
	setup.GameObjects[dagger.ID] = dagger.Properties

	// say "blubb"
	res := actions.Object(obj).Say(area, "blubb")
	s := fmt.Sprintf(setup.Reactions["say"].Statement, "blubb")
	assert.Equal(t, s, res.Statement)
	assert.True(t, res.OK)

	// say "fantasia"
	res = actions.Object(obj).Say(area, "fantasia")
	s = fmt.Sprintf(setup.Reactions["fantasia"].Statement, sword.Properties.Value+dagger.Properties.Value)
	assert.Equal(t, s, res.Statement)
	assert.True(t, res.OK)

	// say "simsalabim"
	area = setup.GetAreaByID(18)
	res = actions.Object(obj).Say(area, "simsalabim")
	assert.Equal(t, setup.Reactions["simsalabim"].Statement, res.Statement)
	assert.True(t, res.OK)
}

func TestFill(t *testing.T) {
	setup.Init()
	area := setup.GetAreaByID(1)
	sword := setup.GetObjectByID(15)

	// put sword into inv
	sword.Properties.Area = setup.INVENTORY
	setup.GameObjects[sword.ID] = sword.Properties

	res := actions.Object(sword).Fill(area)
	s := view.Highlight(fmt.Sprintf(setup.Reactions["unusable"].Statement, sword.Properties.Description.Long), setup.RED)
	assert.Equal(t, s, res.Statement)
	assert.False(t, res.OK)

	// fill jar
	jar := setup.GetObjectByID(30)
	res = actions.Object(jar).Fill(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	// put jar into inv
	jar.Properties.Area = setup.INVENTORY
	setup.GameObjects[jar.ID] = jar.Properties
	res = actions.Object(jar).Fill(area)
	assert.Equal(t, setup.Reactions["noWater"].Statement, res.Statement)
	assert.False(t, res.OK)

	// fill jar at pond
	area = setup.GetAreaByID(3)
	res = actions.Object(jar).Fill(area)
	assert.Equal(t, setup.Reactions["waterUnreachable"].Statement, res.Statement)
	assert.False(t, res.OK)

	// fill jar at well
	area = setup.GetAreaByID(17)
	res = actions.Object(jar).Fill(area)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)

	// fill goblet at well
	goblet := setup.GetObjectByID(44)
	// put goblet into inv
	goblet.Properties.Area = setup.INVENTORY
	setup.GameObjects[goblet.ID] = goblet.Properties
	res = actions.Object(goblet).Fill(area)
	a := strings.Title(goblet.Properties.Description.Article)
	desc := fmt.Sprintf("%s %s", a, goblet.Properties.Description.Short)
	s = fmt.Sprintf(setup.Reactions["unsuitable"].Statement, desc)
	assert.Equal(t, s, res.Statement)
	assert.False(t, res.OK)
}
