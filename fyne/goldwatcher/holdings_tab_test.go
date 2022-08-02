package main

import "testing"

func TestConfig_currentHoldings(t *testing.T) {
	h, err := testApp.currentHoldings()
	if err != nil {
		t.Error(err)
	}

	if len(h) <= 1 {
		t.Error("invalid results")
	}

	if len(h) != 2 {
		t.Error("wrong record count")
	}
}

func TestConfig_getHoldingsSlice(t *testing.T) {
	slice := testApp.getHoldingsSlice()
	if len(slice) != 3 {
		t.Errorf("wrong row count: got %d; want %d", len(slice), 3)
	}
}
