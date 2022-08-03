package main

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestConfig_getToolbar(t *testing.T) {
	tb := testApp.getToolbar()
	if len(tb.Items) != 4 {
		t.Errorf("toolbar items: got %d; want %d", len(tb.Items), 4)
	}
}

func TestConfig_addHoldingsDialog(t *testing.T) {
	_ = testApp.addHoldingsDialog()

	test.Type(testApp.AddHoldingsPurchaseAmountEntry, "1")
	test.Type(testApp.AddHoldingsPurchasePriceEntry, "1234")
	test.Type(testApp.AddHoldingsPurchaseDateEntry, "2022-08-02")

	if testApp.AddHoldingsPurchaseDateEntry.Text != "2022-08-02" {
		t.Errorf("got %s; want %s", testApp.AddHoldingsPurchaseDateEntry.Text, "2022-08-02")
	}

	if testApp.AddHoldingsPurchasePriceEntry.Text != "1234" {
		t.Errorf("got %s; want %s", testApp.AddHoldingsPurchasePriceEntry.Text, "1234")
	}

	if testApp.AddHoldingsPurchaseAmountEntry.Text != "1" {
		t.Errorf("got %s; want %s", testApp.AddHoldingsPurchaseAmountEntry.Text, "1")
	}
}
