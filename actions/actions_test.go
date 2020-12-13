package actions_test

import (
	"fangotasia/actions"
	"fangotasia/grid"
	"fangotasia/setup"
	"fangotasia/view"
	"fmt"
	"os"
	"regexp"
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
	grid.SetupGrid()
	// go to area 1
	area := setup.GetAreaByID(1)
	KO := actions.Parse("nimm zauberschuhe", &area)
	assert.Equal(t, 1, area.ID)
	KO = actions.Parse("trage zauberschuhe", &area)
	assert.Equal(t, 1, area.ID)
	KO = actions.Parse("o", &area)
	assert.Equal(t, 2, area.ID)
	assert.False(t, KO)
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

	// to to area 29
	area = setup.GetAreaByID(29)
	// pick up grub
	obj = actions.Object(setup.GetObjectByID(10))
	res = obj.Take(area)
	assert.Equal(t, setup.Reactions["silly"].Statement, res.Statement)
	assert.False(t, res.OK)

	// pick up stone
	obj = actions.Object(setup.GetObjectByID(20))
	obj.NewAreaID(area.ID)
	res = obj.Take(area)
	assert.Equal(t, setup.Reactions["ok"].Statement, res.Statement)
	assert.True(t, res.OK)

	// pick up sword
	obj = actions.Object(setup.GetObjectByID(15))
	grub := actions.Object(setup.GetObjectByID(10))
	obj.NewAreaID(area.ID)
	res = obj.Take(area)
	r := setup.GetReactionByName("wontLet")
	r.Statement[0] = fmt.Sprintf("%s %s %s",
		strings.Title(grub.Properties.Description.Article),
		grub.Properties.Description.Short,
		r.Statement[0])
	assert.Equal(t, r.Statement, res.Statement)
	assert.False(t, res.OK)

	// wear hood
	hood := actions.Object(setup.GetObjectByID(13))
	hood.NewAreaID(setup.INUSE)
	res = obj.Take(area)
	/*
		r := setup.GetReactionByName("wontLet")
		r.Statement[0] = fmt.Sprintf("%s %s %s",
			strings.Title(grub.Properties.Description.Article),
			grub.Properties.Description.Short,
			r.Statement[0])
		assert.Equal(t, r.Statement, res.Statement)
	*/
	assert.Equal(t, setup.Reactions["hoodInUse"].Statement, res.Statement)
	assert.False(t, res.OK)
}

func TestStab(t *testing.T) {
	setup.Setup()
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
	obj = actions.Object(setup.GetObjectByID(31))
	obj.Use(area)
	obj = actions.Object(setup.GetObjectByID(31))
	res, areaID = obj.Move(area, "o")
	assert.Equal(t, 2, areaID)
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
	setup.Setup()
	// throw sphere
	// not area 4 - sphere breaks
	area := setup.GetAreaByID(1)
	sphere := actions.Object(setup.GetObjectByID(34))
	sphere.NewAreaID(setup.INVENTORY)
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
	sphere.NewAreaID(setup.INVENTORY)
	res = sphere.Throw(area)
	sphere = actions.Object(setup.GetObjectByID(46))
	statement := make([]string, len(setup.Reactions["throwObject"].Statement))
	copy(statement, setup.Reactions["throwObject"].Statement)
	newArea := setup.GetAreaByID(sphere.Properties.Area)
	article := strings.Title(sphere.Properties.Description.Article)
	short := sphere.Properties.Description.Short
	location := newArea.Properties.Description.Long
	re := regexp.MustCompile(`\\n.*$`)
	location = string(re.ReplaceAll([]byte(location), []byte(" ")))
	for i, s := range statement {
		statement[i] = fmt.Sprintf(s, article, short, location)
	}
	assert.Equal(t, statement, res.Statement)
	assert.True(t, res.OK)
	assert.LessOrEqual(t, setup.GetObjectByID(46).Properties.Area, 51)

	// stone
	area = setup.GetAreaByID(9)
	stone := actions.Object(setup.GetObjectByID(20))
	stone.NewAreaID(setup.INVENTORY)
	res = stone.Throw(area)
	stone = actions.Object(setup.GetObjectByID(20))
	statement = make([]string, len(setup.Reactions["throwObject"].Statement))
	copy(statement, setup.Reactions["throwObject"].Statement)
	newArea = setup.GetAreaByID(stone.Properties.Area)
	article = strings.Title(stone.Properties.Description.Article)
	short = stone.Properties.Description.Short
	location = newArea.Properties.Description.Long
	location = string(re.ReplaceAll([]byte(location), []byte(" ")))
	for i, s := range statement {
		statement[i] = fmt.Sprintf(s, article, short, location)
	}
	assert.Equal(t, statement, res.Statement)
	assert.True(t, res.OK)
	assert.LessOrEqual(t, setup.GetObjectByID(46).Properties.Area, 51)

	// stone on tree - map present by setup
	stone.NewAreaID(setup.INVENTORY)
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

	// fango
	area = setup.GetAreaByID(9)
	fango := actions.Object(setup.GetObjectByID(49))
	fango.NewAreaID(setup.INVENTORY)
	res = fango.Throw(area)
	assert.Equal(t, setup.Reactions["throwFango"].Statement, res.Statement)
	assert.True(t, res.OK)

	// fango at dwarf
	area = setup.GetAreaByID(18)
	//fango := actions.Object(setup.GetObjectByID(49))
	res = fango.Throw(area)
	assert.Equal(t, setup.Reactions["hitWithFango"].Statement, res.Statement)
	assert.True(t, res.OK)

	// book
	book := actions.Object(setup.GetObjectByID(11))
	area = setup.GetAreaByID(2)
	res = book.Throw(area)
	assert.Equal(t, setup.Reactions["dontHave"].Statement, res.Statement)
	assert.False(t, res.OK)

	book.NewAreaID(setup.INVENTORY)
	res = book.Throw(area)
	book = actions.Object(setup.GetObjectByID(11))
	statement = make([]string, len(setup.Reactions["throwObject"].Statement))
	copy(statement, setup.Reactions["throwObject"].Statement)
	newArea = setup.GetAreaByID(book.Properties.Area)
	article = strings.Title(book.Properties.Description.Article)
	short = book.Properties.Description.Short
	location = newArea.Properties.Description.Long
	location = string(re.ReplaceAll([]byte(location), []byte(" ")))
	for i, s := range statement {
		statement[i] = fmt.Sprintf(s, article, short, location)
	}
	assert.Equal(t, statement, res.Statement)
	assert.True(t, res.OK)
	assert.LessOrEqual(t, setup.GetObjectByID(46).Properties.Area, 51)

}

func TestRead(t *testing.T) {
	// read panel
	area := setup.GetAreaByID(1)
	panel := actions.Object(setup.GetObjectByID(32))
	res := panel.Read(area)
	assert.Equal(t, view.Highlight(setup.Reactions["panel"].Statement[0], "[green:black:-]"), res.Statement[0])
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
	say := setup.GetReactionByName("say")
	s := fmt.Sprintf(say.Statement[0], "blubb")
	assert.Equal(t, s, res.Statement[0])
	assert.True(t, res.OK)

	/*
		// say "fangotasia"
		res = obj.Say(area, "Fangotasia")
		fangotasia := setup.GetReactionByName("fangotasia")
		s = fmt.Sprintf(fangotasia.Statement[0], sword.Properties.Value+dagger.Properties.Value)
		assert.Equal(t, s, res.Statement[0])
		assert.True(t, res.OK)
	*/

	// say "simsalabim"
	area = setup.GetAreaByID(18)
	res = obj.Say(area, "SimSalabim")
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
	unusable := setup.GetReactionByName("unusable")
	s := view.Highlight(fmt.Sprintf(unusable.Statement[0], sword.Properties.Description.Long), "[red]")
	assert.Equal(t, s, res.Statement[0])
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
	unsuitable := setup.GetReactionByName("unsuitable")
	s = fmt.Sprintf(unsuitable.Statement[0], desc)
	assert.Equal(t, s, res.Statement[0])
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
	feed := setup.GetReactionByName("feed")
	assert.Equal(t, fmt.Sprintf(feed.Statement[0], desc), res.Statement[0])
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
	tree := actions.Object(setup.GetObjectByID(27))

	// eat Tree
	res := tree.Eat(area)
	assert.Equal(t, setup.Reactions["dontSee"].Statement, res.Statement)
	assert.False(t, res.OK)

	// eat Tree
	area = setup.GetAreaByID(9)
	res = tree.Eat(area)
	assert.Equal(t, setup.Reactions["cantEat"].Statement, res.Statement)
	assert.False(t, res.OK)

	// eat jar
	res = jar.Eat(area)
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
