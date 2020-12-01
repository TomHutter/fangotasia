package view_test

import (
	"fantasia/setup"
	"fantasia/view"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	path, _ := os.Getwd()
	setup.PathName = path + "/../"
	setup.Setup()
}

func TestFlashNotice(t *testing.T) {
	assert.False(t, view.FlashNotice())
	view.AddFlashNotice("test", 3, setup.BLUE)
	assert.True(t, view.FlashNotice())
}
