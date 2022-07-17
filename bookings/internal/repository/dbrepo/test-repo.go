package dbrepo

import (
	"errors"
	"time"

	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the database.
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 2, errors.New("some error")
	}

	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database.
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 2 {
		return errors.New("some error")
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for a given room ID, otherwise false.
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start time.Time, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for a given date range.
func (m *testDBRepo) SearchAvailabilityForAllRooms(start time.Time, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID gets a room by the given id.
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}

// GetUserByID returns a user by id.
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

// UpdateUser updates a user in the database.
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Authenticate authenticates a user.
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

// AllReservations returns a slice of all reservations.
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	return []models.Reservation{}, nil
}

// AllNewReservations returns a slice of new reservations.
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	return []models.Reservation{}, nil
}

// GetReservationByID returns one reservation by id.
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	return models.Reservation{}, nil
}

// UpdateReservation updates a reservation in the database.
func (m *testDBRepo) UpdateReservation(res models.Reservation) error {
	return nil
}

// DeleteReservation deletes one reservation by id.
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdateProcessedForReservation updates processed field for a reservation by id.
func (m *testDBRepo) UpdateProcessedForReservation(id int, processed int) error {
	return nil
}

// AllRooms returns all rooms in the database.
func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room

	return rooms, nil
}

// GetRestrctionsForRoomByDate returns all restrictions for a given room over a given date range.
func (m *testDBRepo) GetRestrctionsForRoomByDate(id int, start time.Time, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction

	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction.
func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	return nil
}

// DeleteBlockByID deletes a room restriction.
func (m *testDBRepo) DeleteBlockByID(id int) error {
	return nil
}
