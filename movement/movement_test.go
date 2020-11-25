package movement_test

import (
	"fantasia/config"
	"fantasia/movement"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	config.Init()
}

func TestSurroundings(t *testing.T) {
	area := config.GetAreaByID(1)
	objects := config.ObjectsInArea(area)
	s := strings.Join(movement.Surroundings(area), "\n")
	assert.Contains(t, s, "Torbogen")
	assert.Contains(t, s, objects[0].Properties.Description.Long)
}

func TestDrawMap(t *testing.T) {
	area := config.GetAreaByID(1)
	s := strings.Join(movement.DrawMap(area), "\n")
	assert.Contains(t, s, area.Properties.Description.Short)
}

func TestRevealArea(t *testing.T) {
	movement.RevealArea(3)
	assert.True(t, config.AreaVisible(3))
	movement.RevealArea(51)
	assert.True(t, config.AreaVisible(50))
}
