package actions_test

import (
	"fantasia/actions"
	"fantasia/setup"
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
	assert.Contains(t, res.Statement, setup.Reactions["dontSee"].Statement)
	assert.False(t, res.OK)

	// pick magic shoes
	obj = setup.GetObjectByID(31)
	res = actions.Object(obj).Take(area)
	// reload obj to get new area => should be 1000 by taken into inv
	obj = setup.GetObjectByID(31)
	assert.Contains(t, res.Statement, setup.Reactions["ok"].Statement)
	assert.True(t, res.OK)
	// look up inventory
	inv := setup.ObjectsInArea(setup.GetAreaByID(1000))
	assert.Contains(t, inv, obj)

	// try to pick shoes again
	res = actions.Object(obj).Take(area)
	assert.Contains(t, res.Statement, setup.Reactions["haveAlready"].Statement)
	assert.True(t, res.OK)

	// to to area 4
	area = setup.GetAreaByID(4)
	// pick up gnome
	obj = setup.GetObjectByID(36)
	res = actions.Object(obj).Take(area)
	assert.Contains(t, res.Statement, setup.Reactions["silly"].Statement)
	assert.False(t, res.OK)
}

func TestStab(t *testing.T) {
	// go to area 4
	area := setup.GetAreaByID(4)
	// stab gnome
	gnome := setup.GetObjectByID(36)
	res := actions.Object(gnome).Stab(area)
	assert.Contains(t, res.Statement, setup.Reactions["noTool"].Statement)
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
	assert.Contains(t, res.Statement, setup.Reactions["stabGnome"].Statement)
	assert.False(t, res.OK)
	assert.True(t, res.KO)
}

func TestMove(t *testing.T) {
	// go to area 4
	area := setup.GetAreaByID(1)
	obj := setup.Object{}
	// move north
	res, areaID := actions.Object(obj).Move(area, "n")
	assert.Contains(t, res.Statement, setup.Reactions["noWay"].Statement)
	assert.False(t, res.OK)
	assert.Equal(t, areaID, 1)

	// move east
	res, areaID = actions.Object(obj).Move(area, "o")
	assert.Contains(t, res.Statement, setup.Reactions["noShoes"].Statement)
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
	// reload obj to get new area => should be 1000 (inventory)
	shoes = setup.GetObjectByID(31)
	res := actions.Object(shoes).Use(area)
	// reload obj to get new area => should be 2000 (inUse)
	shoes = setup.GetObjectByID(31)
	assert.Contains(t, res.Statement, setup.Reactions["useShoes"].Statement)
	assert.True(t, res.OK)
	// look up inUse
	inUse := setup.ObjectsInArea(setup.GetAreaByID(2000))
	assert.Contains(t, inUse, shoes)

	// go to area 3
	area = setup.GetAreaByID(3)
	// pick magic shoes
	dagger := setup.GetObjectByID(33)
	actions.Object(dagger).Take(area)
	// reload obj to get new area => should be 1000 (inventory)
	dagger = setup.GetObjectByID(33)
	res = actions.Object(dagger).Use(area)
	// reload obj to get new area => should be 2000 (inUse)
	dagger = setup.GetObjectByID(33)
	assert.Contains(t, res.Statement, setup.Reactions["dontKnowHow"].Statement)
	assert.False(t, res.OK)
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
	assert.Contains(t, res.Statement, setup.Reactions["silly"].Statement)
	assert.Equal(t, areaID, 1)
	assert.False(t, res.OK)
}
