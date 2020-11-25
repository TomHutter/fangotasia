package view_test

import (
	"fantasia/config"
	"fantasia/view"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	config.Init()
}

func TestFlashNotice(t *testing.T) {
	assert.False(t, view.FlashNotice())
	view.AddFlashNotice("test", 3, config.BLUE)
	assert.True(t, view.FlashNotice())
}
