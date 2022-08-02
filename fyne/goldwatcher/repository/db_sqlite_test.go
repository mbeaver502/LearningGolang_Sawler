package repository

import (
	"testing"
	"time"
)

func TestSQLiteRepository_Migrate(t *testing.T) {
	err := testRepo.Migrate()
	if err != nil {
		t.Error(err)
	}
}

func TestSQLiteRepository_InsertHolding(t *testing.T) {
	h := Holdings{
		Amount:        1,
		PurchaseDate:  time.Now(),
		PurchasePrice: 1000,
	}

	res, err := testRepo.InsertHolding(h)
	if err != nil {
		t.Error(err)
	}

	if res.ID < 1 {
		t.Errorf("got %d; want >= 1", res.ID)
	}
}

func TestSQLiteRepository_AllHoldings(t *testing.T) {
	res, err := testRepo.AllHoldings()
	if err != nil {
		t.Error(err)
	}

	if len(res) < 1 {
		t.Error("results empty")
	}

	if len(res) != 1 {
		t.Errorf("wrong rows: got %d; want %d", len(res), 1)
	}
}

func TestSQLiteRepository_GetHoldingByID(t *testing.T) {
	res, err := testRepo.GetHoldingByID(1)
	if err != nil {
		t.Error(err)
	}

	if res.ID != 1 {
		t.Errorf("invalid id: got %d; want %d", res.ID, 1)
	}

	if res.PurchasePrice != 1000 {
		t.Errorf("invalid purchase price: got %d; want %d", res.PurchasePrice, 1000)
	}

	_, err = testRepo.GetHoldingByID(2)
	if err == nil {
		t.Errorf("got %v; want error", err)
	}
}

func TestSQLiteRepository_UpdateHolding(t *testing.T) {
	h, err := testRepo.GetHoldingByID(1)
	if err != nil {
		t.Error(err)
	}

	h.PurchasePrice = 2000

	err = testRepo.UpdateHolding(h.ID, *h)
	if err != nil {
		t.Error(err)
	}

	h, err = testRepo.GetHoldingByID(1)
	if err != nil {
		t.Error(err)
	}

	if h.PurchasePrice != 2000 {
		t.Errorf("invalid purchase price: got %d; want %d", h.PurchasePrice, 2000)
	}
}

func TestSQLiteRepository_DeleteHolding(t *testing.T) {
	err := testRepo.DeleteHolding(1)
	if err != nil {
		t.Error(err)
		if err != errDeleteFailed {
			t.Errorf("wrong err: got %v; want %v", err, errDeleteFailed)
		}
	}

	err = testRepo.DeleteHolding(2)
	if err != errDeleteFailed {
		t.Error("no error when deleting invalid id")
	}
}
