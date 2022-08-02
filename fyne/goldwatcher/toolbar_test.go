package main

import "testing"

func TestConfig_getToolbar(t *testing.T) {
	tb := testApp.getToolbar()
	if len(tb.Items) != 4 {
		t.Errorf("toolbar items: got %d; want %d", len(tb.Items), 4)
	}
}
