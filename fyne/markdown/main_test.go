package main

import (
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func TestConfig_makeUI(t *testing.T) {
	var testCfg config

	edit, preview := testCfg.makeUI()

	// simulate typing inside the edit box
	test.Type(edit, "Hello")

	if preview.String() != "Hello" {
		t.Errorf("got %s; want %s", preview.String(), "Hello")
	}
}

func TestConfig_RunApp(t *testing.T) {
	var testCfg config

	testApp := test.NewApp()
	testWin := testApp.NewWindow("Test Window")

	edit, preview := testCfg.makeUI()

	testCfg.createMenuItems(testWin)

	testWin.SetContent(container.NewHSplit(edit, preview))

	testWin.ShowAndRun()

	test.Type(edit, "Hello")
	if preview.String() != "Hello" {
		t.Errorf("got %s; want %s", preview.String(), "Hello")
	}
}
