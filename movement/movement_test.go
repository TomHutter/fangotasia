package movement_test

import (
	"fantasia/actions"
	"fantasia/movement"
	"fantasia/setup"
	"fantasia/view"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	path, _ := os.Getwd()
	setup.PathName = path + "/../"
	setup.Setup()
}

func TestSurroundings(t *testing.T) {
	area := setup.GetAreaByID(1)
	objects := setup.ObjectsInArea(area)
	s := strings.Join(view.Surroundings(area), "\n")
	assert.Contains(t, s, "Torbogen")
	assert.Contains(t, s, objects[0].Properties.Description.Long)
}

func TestDrawMap(t *testing.T) {
	// put map in use
	m := actions.Object(setup.GetObjectByID(47))
	m.NewAreaID(setup.INUSE)
	area := setup.GetAreaByID(1)
	s := strings.Join(movement.DrawMap(area), "\n")
	assert.Contains(t, s, area.Properties.Description.Short)
}

func TestRevealArea(t *testing.T) {
	movement.RevealArea(3)
	assert.True(t, setup.AreaVisible(3))
	movement.RevealArea(51)
	assert.True(t, setup.AreaVisible(50))
}
