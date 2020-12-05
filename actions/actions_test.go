package actions_test

import (
	"fantasia/actions"
	"fantasia/setup"
	"fantasia/view"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	path, _ := os.Getwd()
	setup.PathName = path + "/../"
	fmt.Println(setup.PathName)
	setup.Setup()
}

func TestParse(t *testing.T) {
	setup.Setup()
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
	setup.Setup()
	// go to area 1
	area := setup.GetAreaByID(1)
	// pick up dwarf
	obj := actions.Object(setup.GetObjectByID(18))
	res := obj.Take(area)
	// dwarf not in area
	assert.Equal(t, setup.Reactions["dontSee"].Statement, res.Statement)
	assert.False(t, res.OK)

	// pick magic shoes
	obj = actions.Object(setup.GetObjectByID(31))
	res = obj.Take(area)
	// reload obj to get new area => should be setup.INVENTORY by taken into inv
	obj = actions.Object(setup.GetObjectByID(31))
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)
	// look up inventory
	inv := setup.ObjectsInArea(setup.GetAreaByID(setup.INVENTORY))
	assert.Contains(t, inv, setup.Object(obj))

	// try to pick shoes again
	res = obj.Take(area)
	assert.Equal(t, setup.Reactions["haveAlready"].Statement, res.Statement)
	assert.True(t, res.OK)

	// to to area 4
	area = setup.GetAreaByID(4)
	// pick up gnome
	obj = actions.Object(setup.GetObjectByID(36))
	res = obj.Take(area)
	assert.Equal(t, setup.Reactions["silly"].Statement, res.Statement)
	assert.False(t, res.OK)
}

func TestStab(t *testing.T) {
	// go to area 4
	area := setup.GetAreaByID(4)
	// stab gnome
	gnome := actions.Object(setup.GetObjectByID(36))
	res := gnome.Stab(area)
	assert.Equal(t, setup.Reactions["noTool"].Statement, res.Statement)
	assert.False(t, res.OK)

	// go to area 3
	area = setup.GetAreaByID(3)
	// pick up dwarf dagger
	obj := actions.Object(setup.GetObjectByID(33))
	res = obj.Take(area)
	// go to area 4
	area = setup.GetAreaByID(4)
	// stab gnome
	res = gnome.Stab(area)
	assert.Equal(t, setup.Reactions["stabGnome"].Statement, res.Statement)
	assert.False(t, res.OK)
	assert.True(t, res.KO)
}

func TestMove(t *testing.T) {
	// go to area 4
	area := setup.GetAreaByID(1)
	obj := actions.Object{}
	// move north
	res, areaID := obj.Move(area, "n")
	assert.Equal(t, setup.Reactions["noWay"].Statement, res.Statement)
	assert.False(t, res.OK)
	assert.Equal(t, areaID, 1)

	// move east
	res, areaID = obj.Move(area, "o")
	assert.Equal(t, setup.Reactions["noShoes"].Statement, res.Statement)
	assert.False(t, res.OK)
	assert.True(t, res.KO)
	assert.Equal(t, areaID, 1)

	// pick magic shoes
	obj = actions.Object(setup.GetObjectByID(31))
	obj.Take(area)
	obj.Use(area)
	res, areaID = obj.Move(area, "o")
	assert.Equal(t, areaID, 2)
	assert.True(t, res.OK)
}

func TestUse(t *testing.T) {
	setup.Setup()
	// go to area 1
	area := setup.GetAreaByID(1)
	// pick magic shoes
	shoes := actions.Object(setup.GetObjectByID(31))
	shoes.Take(area)
	// reload obj to get new area => should be setup.INVENTORY (inventory)
	shoes = actions.Object(setup.GetObjectByID(31))
	res := shoes.Use(area)
	// reload obj to get new area => should be setup.INUSE (inUse)
	shoes = actions.Object(setup.GetObjectByID(31))
	assert.Equal(t, setup.Reactions["useShoes"].Statement, res.Statement)
	assert.True(t, res.OK)
	// look up inUse
	inUse := setup.ObjectsInArea(setup.GetAreaByID(setup.INUSE))
	assert.Contains(t, inUse, setup.Object(shoes))

	// go to area 3
	area = setup.GetAreaByID(3)
	// pick magic shoes
	dagger := actions.Object(setup.GetObjectByID(33))
	dagger.Take(area)
	// reload obj to get new area => should be setup.INVENTORY (inventory)
	dagger = actions.Object(setup.GetObjectByID(33))
	res = dagger.Use(area)
	// reload obj to get new area => should be setup.INUSE (inUse)
	dagger = actions.Object(setup.GetObjectByID(33))
	assert.Equal(t, setup.Reactions["dontKnowHow"].Statement, res.Statement)
	assert.False(t, res.OK)

	// use map
	m := actions.Object(setup.GetObjectByID(47))
	m.NewAreaID(setup.INVENTORY)
	res = m.Use(area)
	assert.Equal(t, setup.Reactions["useMap"].Statement, res.Statement)
	m = actions.Object(setup.GetObjectByID(47))
	assert.True(t, res.OK)
	assert.Equal(t, setup.INUSE, m.Properties.Area)
}

func TestClimb(t *testing.T) {
	// go to area 9
	area := setup.GetAreaByID(9)
	tree := actions.Object(setup.GetObjectByID(27))
	res, areaID := tree.Climb(area)
	assert.Equal(t, areaID, 31)
	assert.True(t, res.OK)

	// got to area 31
	area = setup.GetAreaByID(31)
	tree = actions.Object(setup.GetObjectByID(27))
	res, areaID = tree.Climb(area)
	assert.Equal(t, areaID, 9)
	assert.True(t, res.OK)

	// got to area 1
	area = setup.GetAreaByID(1)
	tree = actions.Object(setup.GetObjectByID(32))
	res, areaID = tree.Climb(area)
	assert.Equal(t, setup.Reactions["silly"].Statement, res.Statement)
	assert.Equal(t, areaID, 1)
	assert.False(t, res.OK)
}

func TestThrow(t *testing.T) {
	// not area 4 - sphere breaks
	area := setup.GetAreaByID(1)
	sphere := actions.Object(setup.GetObjectByID(34))
	res := sphere.Throw(area)
	assert.Equal(t, setup.Reactions["brokenSphere"].Statement, res.Statement)
	assert.Equal(t, 0, setup.GetObjectByID(34).Properties.Area)
	assert.False(t, res.OK)

	// area 4 - gnome breaks
	area = setup.GetAreaByID(4)
	res = sphere.Throw(area)
	assert.Equal(t, setup.Reactions["squashed"].Statement, res.Statement)
	assert.Equal(t, 0, setup.GetObjectByID(36).Properties.Area)
	assert.True(t, res.OK)

	// golden sphere
	area = setup.GetAreaByID(1)
	sphere = actions.Object(setup.GetObjectByID(46))
	res = sphere.Throw(area)
	assert.Equal(t, setup.Reactions["throw"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, 1, setup.GetObjectByID(46).Properties.Area)

	// stone
	area = setup.GetAreaByID(9)
	stone := actions.Object(setup.GetObjectByID(20))
	res = stone.Throw(area)
	assert.Equal(t, setup.Reactions["throw"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, 9, setup.GetObjectByID(20).Properties.Area)

	// stone on tree - map present by setup
	area = setup.GetAreaByID(31)
	assert.False(t, setup.Flags["MapMissed"])
	res = stone.Throw(area)
	assert.Equal(t, setup.Reactions["missMap"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.True(t, setup.Flags["MapMissed"])
	assert.Equal(t, 9, setup.GetObjectByID(20).Properties.Area)

	// stone on tree - map present by setup, 2nd try
	res = stone.Throw(area)
	assert.Equal(t, setup.Reactions["hitMap"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, 9, setup.GetObjectByID(20).Properties.Area)
	assert.Equal(t, 9, setup.GetObjectByID(47).Properties.Area)
}

func TestRead(t *testing.T) {
	// read panel
	area := setup.GetAreaByID(1)
	panel := actions.Object(setup.GetObjectByID(32))
	res := panel.Read(area)
	assert.Equal(t, view.Highlight(setup.Reactions["panel"].Statement, "[green]"), res.Statement)
	assert.True(t, res.OK)

	// read letter
	letter := actions.Object(setup.GetObjectByID(12))
	res = letter.Read(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	// read letter in inv
	letter.Properties.Area = setup.INVENTORY
	res = letter.Read(area)
	assert.Equal(t, setup.Reactions["paper"].Statement, res.Statement)
	assert.True(t, res.OK)
}

func TestOpen(t *testing.T) {
	// open panel
	area := setup.GetAreaByID(1)
	panel := actions.Object(setup.GetObjectByID(32))
	res := panel.Open(area)
	assert.Equal(t, setup.Reactions["dontKnowHow"].Statement, res.Statement)
	assert.False(t, res.OK)

	// open red box
	box := actions.Object(setup.GetObjectByID(35))
	res = box.Open(area)
	assert.Equal(t, setup.Reactions["dontSee"].Statement, res.Statement)
	assert.False(t, res.OK)

	// correct area
	area = setup.GetAreaByID(4)
	res = box.Open(area)
	assert.Equal(t, setup.Reactions["noKey"].Statement, res.Statement)
	assert.False(t, res.OK)
	assert.False(t, setup.Flags["BoxOpen"])

	// key in inv
	key := actions.Object(setup.GetObjectByID(26))
	key.NewAreaID(setup.INVENTORY)
	res = box.Open(area)
	assert.Equal(t, setup.Reactions["openBox"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.True(t, setup.Flags["BoxOpen"])

	// reopen box
	res = box.Open(area)
	assert.Equal(t, setup.Reactions["alreadyOpen"].Statement, res.Statement)
	assert.True(t, res.OK)

	// open door
	area = setup.GetAreaByID(25)
	door := actions.Object(setup.GetObjectByID(40))
	key.NewAreaID(9)
	res = door.Open(area)
	assert.Equal(t, setup.Reactions["noKey"].Statement, res.Statement)
	assert.False(t, res.OK)
	assert.False(t, setup.Flags["DoorOpen"])

	// key in inv
	key.NewAreaID(setup.INVENTORY)
	res = door.Open(area)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.True(t, setup.Flags["DoorOpen"])

	// reopen door (from the other side)
	area = setup.GetAreaByID(30)
	door = actions.Object(setup.GetObjectByID(45))
	res = door.Open(area)
	assert.Equal(t, setup.Reactions["alreadyOpen"].Statement, res.Statement)
	assert.True(t, res.OK)
}

func TestPut(t *testing.T) {
	setup.Setup()
	area := setup.GetAreaByID(1)
	sword := actions.Object(setup.GetObjectByID(15))
	dagger := actions.Object(setup.GetObjectByID(33))

	// put sword into inv
	sword.NewAreaID(setup.INVENTORY)

	// try to put dagger
	res := dagger.Put(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	// to put dagger
	res = sword.Put(area)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, 1, setup.GetObjectByID(15).Properties.Area)
}

func TestSay(t *testing.T) {
	setup.Setup()
	area := setup.GetAreaByID(1)
	sword := actions.Object(setup.GetObjectByID(15))
	dagger := actions.Object(setup.GetObjectByID(33))
	obj := actions.Object{}

	// put sword in area 1
	sword.NewAreaID(1)

	// put dagger in area 1
	dagger.NewAreaID(1)

	// say "blubb"
	res := obj.Say(area, "blubb")
	s := fmt.Sprintf(setup.Reactions["say"].Statement, "blubb")
	assert.Equal(t, s, res.Statement)
	assert.True(t, res.OK)

	// say "fangotasia"
	res = obj.Say(area, "fangotasia")
	s = fmt.Sprintf(setup.Reactions["fangotasia"].Statement, sword.Properties.Value+dagger.Properties.Value)
	assert.Equal(t, s, res.Statement)
	assert.True(t, res.OK)

	// say "simsalabim"
	area = setup.GetAreaByID(18)
	res = obj.Say(area, "simsalabim")
	assert.Equal(t, setup.Reactions["simsalabim"].Statement, res.Statement)
	assert.True(t, res.OK)
}

func TestFill(t *testing.T) {
	setup.Setup()
	area := setup.GetAreaByID(1)
	sword := actions.Object(setup.GetObjectByID(15))

	// put sword into inv
	sword.NewAreaID(setup.INVENTORY)

	res := sword.Fill(area)
	s := view.Highlight(fmt.Sprintf(setup.Reactions["unusable"].Statement, sword.Properties.Description.Long), "[red]")
	assert.Equal(t, s, res.Statement)
	assert.False(t, res.OK)

	// fill jar
	jar := actions.Object(setup.GetObjectByID(30))
	res = jar.Fill(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	// put jar into inv
	jar.NewAreaID(setup.INVENTORY)
	res = jar.Fill(area)
	assert.Equal(t, setup.Reactions["noWater"].Statement, res.Statement)
	assert.False(t, res.OK)

	// fill jar at pond
	area = setup.GetAreaByID(3)
	res = jar.Fill(area)
	assert.Equal(t, setup.Reactions["waterUnreachable"].Statement, res.Statement)
	assert.False(t, res.OK)

	// fill jar at well
	area = setup.GetAreaByID(17)
	res = jar.Fill(area)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)

	// fill goblet at well
	goblet := actions.Object(setup.GetObjectByID(44))
	// put goblet into inv
	goblet.NewAreaID(setup.INVENTORY)
	res = goblet.Fill(area)
	a := strings.Title(goblet.Properties.Description.Article)
	desc := fmt.Sprintf("%s %s", a, goblet.Properties.Description.Short)
	s = fmt.Sprintf(setup.Reactions["unsuitable"].Statement, desc)
	assert.Equal(t, s, res.Statement)
	assert.False(t, res.OK)
}

func TestFeed(t *testing.T) {
	setup.Setup()
	area := setup.GetAreaByID(1)
	panel := actions.Object(setup.GetObjectByID(32))

	// feed panel
	res := panel.Feed(area)
	assert.Equal(t, setup.Reactions["dontKnowHow"].Statement, res.Statement)
	assert.False(t, res.OK)

	// feed panel
	area = setup.GetAreaByID(18)
	dwarf := actions.Object(setup.GetObjectByID(18))
	res = dwarf.Feed(area)
	a := strings.Title(dwarf.Properties.Description.Article)
	desc := fmt.Sprintf("%s %s", a, dwarf.Properties.Description.Short)
	assert.Equal(t, fmt.Sprintf(setup.Reactions["feed"].Statement, desc), res.Statement)
	assert.False(t, res.OK)

	// feed baer
	area = setup.GetAreaByID(19)
	baer := actions.Object(setup.GetObjectByID(16))
	res = baer.Feed(area)
	assert.Equal(t, setup.Reactions["feedBaer"].Statement, res.Statement)
	assert.False(t, res.OK)

	berries := actions.Object(setup.GetObjectByID(23))
	// put berries into inv
	berries.NewAreaID(setup.INVENTORY)
	// feed baer
	res = baer.Feed(area)
	assert.Equal(t, setup.Reactions["feedBaerWithBerries"].Statement, res.Statement)
	assert.True(t, res.OK)
	baer = actions.Object(setup.GetObjectByID(16))
	assert.Equal(t, 0, baer.Properties.Area)
	berries = actions.Object(setup.GetObjectByID(23))
	assert.Equal(t, 0, berries.Properties.Area)
}

func TestCut(t *testing.T) {
	area := setup.GetAreaByID(1)
	panel := actions.Object(setup.GetObjectByID(32))

	// cut panel
	res := panel.Cut(area)
	assert.Equal(t, setup.Reactions["dontKnowHow"].Statement, res.Statement)
	assert.False(t, res.OK)

	// cut fruit
	area = setup.GetAreaByID(26)
	fruit := actions.Object(setup.GetObjectByID(14))
	res = fruit.Cut(area)
	assert.Equal(t, setup.Reactions["noTool"].Statement, res.Statement)
	assert.False(t, res.OK)

	// put dagger into inv
	dagger := actions.Object(setup.GetObjectByID(33))
	dagger.NewAreaID(setup.INVENTORY)
	res = fruit.Cut(area)
	assert.Equal(t, setup.Reactions["cutFruit"].Statement, res.Statement)
	assert.True(t, res.OK)
	// ring present?
	ring := actions.Object(setup.GetObjectByID(37))
	assert.Equal(t, area.ID, ring.Properties.Area)

	res = fruit.Cut(area)
	assert.Equal(t, setup.Reactions["fruitAlreadyCut"].Statement, res.Statement)
	assert.True(t, res.OK)
}

func TestCatapult(t *testing.T) {
	setup.Setup()
	area := setup.GetAreaByID(1)
	stone := actions.Object(setup.GetObjectByID(20))
	dagger := actions.Object(setup.GetObjectByID(33))

	// catapult dagger
	res := dagger.Catapult(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	// put stone into inv
	dagger.NewAreaID(setup.INVENTORY)
	// catapult dagger
	res = dagger.Catapult(area)
	assert.Equal(t, setup.Reactions["dontKnowHow"].Statement, res.Statement)
	assert.False(t, res.OK)

	// catapult stone
	res = stone.Catapult(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	// put stone into inv
	stone.NewAreaID(setup.INVENTORY)
	// catapult stone
	res = stone.Catapult(area)
	assert.Equal(t, setup.Reactions["noTool"].Statement, res.Statement)
	assert.False(t, res.OK)

	area = setup.GetAreaByID(12)
	// catapult stone
	res = stone.Catapult(area)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, 29, actions.Object(setup.GetObjectByID(stone.ID)).Properties.Area)
	grub := actions.Object(setup.GetObjectByID(10))
	assert.Equal(t, 0, grub.Properties.Area)
}

func TestDrink(t *testing.T) {
	setup.Setup()
	area := setup.GetAreaByID(1)
	jar := actions.Object(setup.GetObjectByID(30))

	// drinking only from jar
	res := jar.Drink(area)
	assert.Equal(t, setup.Reactions["noTool"].Statement, res.Statement)
	assert.False(t, res.OK)

	// put jar into inv
	jar.NewAreaID(setup.INVENTORY)
	res = jar.Drink(area)
	assert.Equal(t, setup.Reactions["jarEmpty"].Statement, res.Statement)
	assert.False(t, res.OK)

	// fill jar
	area = setup.GetAreaByID(17)
	jar.Fill(area)
	// reload jar
	jar = actions.Object(setup.GetObjectByID(30))
	res = jar.Drink(area)
	assert.Equal(t, setup.Reactions["drinkJar"].Statement, res.Statement)
	jar = actions.Object(setup.GetObjectByID(30))
	assert.True(t, res.OK)
	assert.Equal(t, setup.Conditions["jar"]["empty"], jar.Properties.Description.Long)
}

func TestEat(t *testing.T) {
	setup.Setup()
	area := setup.GetAreaByID(17)
	jar := actions.Object(setup.GetObjectByID(30))
	cake := actions.Object(setup.GetObjectByID(9))
	berries := actions.Object(setup.GetObjectByID(23))

	// eat jar
	res := jar.Eat(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	// put jar into inv
	jar.NewAreaID(setup.INVENTORY)
	res = jar.Eat(area)
	assert.Equal(t, setup.Reactions["cantEat"].Statement, res.Statement)
	assert.False(t, res.OK)

	// eat cake
	cake.NewAreaID(setup.INVENTORY)
	res = cake.Eat(area)
	assert.Equal(t, setup.Reactions["eatCake"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, setup.INUSE, actions.Object(setup.GetObjectByID(cake.ID)).Properties.Area)

	// eat berries
	berries.NewAreaID(setup.INVENTORY)
	res = berries.Eat(area)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, 0, actions.Object(setup.GetObjectByID(berries.ID)).Properties.Area)
}

func TestSpin(t *testing.T) {
	setup.Setup()
	area := setup.GetAreaByID(27)
	jar := actions.Object(setup.GetObjectByID(30))
	ring := actions.Object(setup.GetObjectByID(37))

	// spin jar
	res := jar.Spin(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	// put jar into inv
	jar.NewAreaID(setup.INVENTORY)
	res = jar.Spin(area)
	assert.Equal(t, setup.Reactions["spin"].Statement, res.Statement)
	assert.True(t, res.OK)

	// spin iron ring
	ring.NewAreaID(setup.INVENTORY)
	res = ring.Spin(area)
	assert.Equal(t, setup.Reactions["spinRing"].Statement, res.Statement)
	assert.True(t, res.OK)
	assert.Equal(t, setup.Conditions["ring"]["golden"], setup.GetObjectByID(ring.ID).Properties.Description.Long)
	assert.Equal(t, area.ID, setup.GetObjectByID(ring.ID).Properties.Area)

	// spin golden ring
	ring = actions.Object(setup.GetObjectByID(37))
	ring.NewAreaID(setup.INVENTORY)
	res = ring.Spin(area)
	assert.Equal(t, setup.Reactions["spin"].Statement, res.Statement)
	assert.True(t, res.OK)
}

func TestWater(t *testing.T) {
	setup.Setup()
	area := setup.GetAreaByID(14)
	jar := actions.Object(setup.GetObjectByID(30))
	bush := actions.Object(setup.GetObjectByID(22))
	tree := actions.Object(setup.GetObjectByID(27))
	panel := actions.Object(setup.GetObjectByID(32))

	// water bush no Jar
	res := bush.Water(area)
	assert.Equal(t, setup.Reactions["noJar"].Statement, res.Statement)
	assert.False(t, res.OK)

	// put jar into inv
	jar.NewAreaID(setup.INVENTORY)
	res = bush.Water(area)
	assert.Equal(t, setup.Reactions["jarEmpty"].Statement, res.Statement)
	assert.False(t, res.OK)

	// fill jar
	jar.NewCondition(setup.Conditions["jar"]["full"])
	res = bush.Water(area)
	assert.Equal(t, setup.Reactions["waterBush"].Statement, res.Statement)
	assert.True(t, res.OK)

	// water again
	res = bush.Water(area)
	assert.Equal(t, setup.Reactions["jarEmpty"].Statement, res.Statement)
	assert.False(t, res.OK)

	// water tree
	area = setup.GetAreaByID(9)
	jar.NewCondition(setup.Conditions["jar"]["full"])
	res = tree.Water(area)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)

	// water panel
	area = setup.GetAreaByID(1)
	jar.NewCondition(setup.Conditions["jar"]["full"])
	res = panel.Water(area)
	assert.Equal(t, setup.Reactions["silly"].Statement, res.Statement)
	assert.False(t, res.OK)
}
