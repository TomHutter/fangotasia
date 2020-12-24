package setup_test

import (
	"fangotasia/setup"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	path, _ := os.Getwd()
	setup.PathName = path + "/../"
	setup.Setup()
	assert.Greater(t, len(setup.GameObjects), 1)
	assert.Greater(t, len(setup.GameAreas), 1)
	assert.Greater(t, len(setup.Overwrites), 1)
	assert.Greater(t, len(setup.Reactions), 1)
	assert.Greater(t, len(setup.Verbs), 1)
	assert.Greater(t, setup.BoxLen, 1)
	assert.Greater(t, len(setup.Map), 1)
}

func TestGameObjects(t *testing.T) {
	assert.IsType(t, setup.ObjectProperties{}, setup.GameObjects[9])
}

func TestGetAreaByID(t *testing.T) {
	area := setup.GetAreaByID(9)
	assert.Equal(t, 9, area.ID, "Area ID not equal.")
	assert.Contains(t, area.Properties.Description[setup.Language].Long,
		area.Properties.Description[setup.Language].Short)
	assert.IsType(t, [4]int{}, area.Properties.Directions, "Area directions have to be type [4]int.")
	assert.IsType(t, 0, area.Properties.Coordinates.X, "Area x coordinates shoud by type int.")
	assert.IsType(t, 0, area.Properties.Coordinates.Y, "Area y coordinates should be type int.")
}

func TestGetObjectByID(t *testing.T) {
	obj := setup.GetObjectByID(15)
	assert.Equal(t, 15, obj.ID, "Object ID not equal.")
	assert.Contains(t, obj.Properties.Description[setup.Language].Long, obj.Properties.Description[setup.Language].Short)
	assert.IsType(t, "das", obj.Properties.Description[setup.Language].Article, "Object article not type string.")
	assert.Equal(t, 20, obj.Properties.Area, "Object area not equal.")
	assert.Equal(t, 10, obj.Properties.Value, "Object value not equal.")
}

func TestObjectsInArea(t *testing.T) {
	area := setup.GetAreaByID(1)
	objects := setup.ObjectsInArea(area)
	assert.Equal(t, 3, len(objects), "There sould be 3 object in area 1.")
}

func TestGetOverwriteByArea(t *testing.T) {
	_, o := setup.GetOverwriteByArea(54)
	assert.Equal(t, 54, o.Area, "Object ID not equal.")
	assert.Contains(t, o.Content[setup.Language][0], "┏━━━━━━━━┻━━━━━━┻━┓")
}
