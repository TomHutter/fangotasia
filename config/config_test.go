package config_test

import (
	"fantasia/config"
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
	assert.Equal(t, "auf einer blühenden Heide.", area.Properties.Description.Long, "Long area description not equal.")
	assert.Equal(t, "Heide", area.Properties.Description.Short, "Short area description not equal.")
	assert.Equal(t, [4]int{15, 3, 10, 8}, area.Properties.Directions, "Area directions not equal.")
	assert.Equal(t, 2, area.Properties.Coordinates.X, "Area x coordinates not equal.")
	assert.Equal(t, 10, area.Properties.Coordinates.Y, "Area y coordinates not equal.")
	assert.Equal(t, false, area.Properties.Visited, "Area should not be visited.")
}

func TestGetObjectByID(t *testing.T) {
	obj := config.GetObjectByID(15)
	assert.Equal(t, 15, obj.ID, "Object ID not equal.")
	assert.Equal(t, "::ein Feuerschwert::", obj.Properties.Description.Long, "Long object description not equal.")
	assert.Equal(t, "Feuerschwert", obj.Properties.Description.Short, "Short object description not equal.")
	assert.Equal(t, "das", obj.Properties.Description.Article, "Object article not equal.")
	assert.Equal(t, 20, obj.Properties.Area, "Object area not equal.")
	assert.Equal(t, 10, obj.Properties.Value, "Object value not equal.")
}
func TestObjectsInArea(t *testing.T) {
	assert.Equal(t, "einen Zauberkuchen", config.GameObjects[9].Description.Long, "The cake is a lie.")
}
