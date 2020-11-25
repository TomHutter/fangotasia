package config_test

import (
	"fantasia/config"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	config.Init()
	assert.Greater(t, len(config.GameObjects), 1)
	assert.Greater(t, len(config.GameAreas), 1)
	assert.Greater(t, len(config.Overwrites), 1)
	assert.Greater(t, len(config.Answers), 1)
	assert.Greater(t, len(config.Verbs), 1)
	assert.Greater(t, config.BoxLen, 1)
	assert.Greater(t, len(config.Map), 1)
}

func TestGameObjects(t *testing.T) {
	//config.Init()
	assert.Equal(t, "einen Zauberkuchen", config.GameObjects[9].Description.Long, "The cake is a lie.")
}

func TestGetAreaByID(t *testing.T) {
	area := config.GetAreaByID(9)
	assert.Equal(t, 9, area.ID, "Area ID not equal.")
	assert.True(t, strings.Contains(area.Properties.Description.Long,
		area.Properties.Description.Short), "Short area description is not in long description.")
	assert.IsType(t, [4]int{}, area.Properties.Directions, "Area directions have to be type [4]int.")
	assert.IsType(t, 0, area.Properties.Coordinates.X, "Area x coordinates shoud by type int.")
	assert.IsType(t, 0, area.Properties.Coordinates.Y, "Area y coordinates should be type int.")
	assert.False(t, area.Properties.Visited, "Area should not be visited.")
}

func TestGetObjectByID(t *testing.T) {
	obj := config.GetObjectByID(15)
	assert.Equal(t, 15, obj.ID, "Object ID not equal.")
	assert.True(t, strings.Contains(obj.Properties.Description.Long,
		obj.Properties.Description.Short), "Short object description is not in long description.")
	assert.IsType(t, "das", obj.Properties.Description.Article, "Object article not type string.")
	assert.Equal(t, 20, obj.Properties.Area, "Object area not equal.")
	assert.Equal(t, 10, obj.Properties.Value, "Object value not equal.")
}

func TestGetObjectByName(t *testing.T) {
	obj := config.GetObjectByName("Zauberkuchen")
	assert.Equal(t, 9, obj.ID, "Object ID not equal.")
	assert.Equal(t, "einen Zauberkuchen", obj.Properties.Description.Long, "The cake is a lie.")
}

func TestObjectsInArea(t *testing.T) {
	area := config.GetAreaByID(1)
	objects := config.ObjectsInArea(area)
	assert.Equal(t, 3, len(objects), "There sould be 3 object in area 1.")
}

func TestGetOverwriteByArea(t *testing.T) {
	o := config.GetOverwriteByArea(52)
	assert.Equal(t, 52, o.Area, "Object ID not equal.")
	assert.True(t, strings.Contains(o.Content[1], "Felswand"))
}
