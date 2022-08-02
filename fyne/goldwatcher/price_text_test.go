package main

import (
	"testing"
)

func TestConfig_getPriceText(t *testing.T) {
	open, _, _ := testApp.getPriceText()
	if open.Text != "Open: $1762.7324 USD" {
		t.Errorf("got %s; want %s", open.Text, "Open: $1762.7324 USD")
	}
}
