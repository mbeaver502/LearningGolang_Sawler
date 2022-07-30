package data

import "testing"

func Test_Ping(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("failed to ping DB", err)
	}
}

func TestBook_GetAll(t *testing.T) {
	all, err := models.Book.GetAll()
	if err != nil {
		t.Error(err)
	}

	if len(all) != 1 {
		t.Errorf("got %d; want %d records", len(all), 1)
	}
}

func TestBook_GetOneByID(t *testing.T) {
	b, err := models.Book.GetOneById(1)
	if err != nil {
		t.Error(err)
	}

	if b.Title != "My Book" {
		t.Errorf("got %s; want %s", b.Title, "My Book`")
	}
}

func TestBook_GetOneBySlug(t *testing.T) {
	b, err := models.Book.GetOneBySlug("my-book")
	if err != nil {
		t.Error(err)
	}

	if b.Title != "My Book" {
		t.Errorf("got %s; want %s", b.Title, "My Book`")
	}

	_, err = models.Book.GetOneBySlug("does-not-exist")
	if err == nil {
		t.Errorf("got %v; want %s", err, "some error")
	}
}
