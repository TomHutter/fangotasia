package view_test

import (
	"fangotasia/setup"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	path, _ := os.Getwd()
	setup.PathName = path + "/../"
	setup.Setup()
}
