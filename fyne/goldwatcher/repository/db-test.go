package repository

import "time"

type TestRepository struct{}

func NewTestRepository() *TestRepository {
	return &TestRepository{}
}

func (repo *TestRepository) Migrate() error {
	return nil
}

func (repo *TestRepository) InsertHolding(h Holdings) (*Holdings, error) {
	return &Holdings{
		ID:            1,
		Amount:        1,
		PurchaseDate:  time.Now(),
		PurchasePrice: 1,
	}, nil
}

func (repo *TestRepository) AllHoldings() ([]Holdings, error) {
	return []Holdings{{
		ID:            1,
		Amount:        1,
		PurchaseDate:  time.Now(),
		PurchasePrice: 1,
	}}, nil
}

func (repo *TestRepository) GetHoldingByID(id int64) (*Holdings, error) {
	return &Holdings{
		ID:            id,
		Amount:        1,
		PurchaseDate:  time.Now(),
		PurchasePrice: 1,
	}, nil
}

func (repo *TestRepository) UpdateHolding(id int64, updated Holdings) error {
	return nil
}

func (repo *TestRepository) DeleteHolding(id int64) error {
	return nil
}
