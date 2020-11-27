package setup_test

import (
	"fantasia/setup"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	setup.Init()
	assert.Greater(t, len(setup.GameObjects), 1)
	assert.Greater(t, len(setup.GameAreas), 1)
	assert.Greater(t, len(setup.Overwrites), 1)
	assert.Greater(t, len(setup.Reactions), 1)
	assert.Greater(t, len(setup.Verbs), 1)
	assert.Greater(t, setup.BoxLen, 1)
	assert.Greater(t, len(setup.Map), 1)
}

func TestGameObjects(t *testing.T) {
	//setup.Init()
	assert.Equal(t, "einen Zauberkuchen", setup.GameObjects[9].Description.Long, "The cake is a lie.")
}

func TestGetAreaByID(t *testing.T) {
	area := setup.GetAreaByID(9)
	assert.Equal(t, 9, area.ID, "Area ID not equal.")
	assert.Contains(t, area.Properties.Description.Long, area.Properties.Description.Short)
	assert.IsType(t, [4]int{}, area.Properties.Directions, "Area directions have to be type [4]int.")
	assert.IsType(t, 0, area.Properties.Coordinates.X, "Area x coordinates shoud by type int.")
	assert.IsType(t, 0, area.Properties.Coordinates.Y, "Area y coordinates should be type int.")
	assert.False(t, area.Properties.Visited, "Area should not be visited.")
}

func TestGetObjectByID(t *testing.T) {
	obj := setup.GetObjectByID(15)
	assert.Equal(t, 15, obj.ID, "Object ID not equal.")
	assert.Contains(t, obj.Properties.Description.Long, obj.Properties.Description.Short)
	assert.IsType(t, "das", obj.Properties.Description.Article, "Object article not type string.")
	assert.Equal(t, 20, obj.Properties.Area, "Object area not equal.")
	assert.Equal(t, 10, obj.Properties.Value, "Object value not equal.")
}

func TestGetObjectByName(t *testing.T) {
	obj := setup.GetObjectByName("Zauberkuchen")
	assert.Equal(t, 9, obj.ID, "Object ID not equal.")
	assert.Equal(t, "einen Zauberkuchen", obj.Properties.Description.Long, "The cake is a lie.")
}

func TestObjectsInArea(t *testing.T) {
	area := setup.GetAreaByID(1)
	objects := setup.ObjectsInArea(area)
	assert.Equal(t, 3, len(objects), "There sould be 3 object in area 1.")
}

func TestGetOverwriteByArea(t *testing.T) {
	o := setup.GetOverwriteByArea(52)
	assert.Equal(t, 52, o.Area, "Object ID not equal.")
	assert.Contains(t, o.Content[1], "Felswand")
}