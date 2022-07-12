package dbrepo

import (
	"context"
	"time"

	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the database.
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	stmt := `insert into reservations (first_name, last_name, email, phone, 
									   start_date, end_date, room_id, 
									   created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				returning id`

	// Use a context so that we can timeout any transactions that exceed a certain time limit.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database.
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, 
											restriction_id, created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6, $7)`

	// Use a context so that we can timeout any transactions that exceed a certain time limit.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for a given room ID, otherwise false.
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start time.Time, end time.Time, roomID int) (bool, error) {
	query := `select count(id) from room_restrictions where $1 < end_date and $2 > start_date and room_id = $3`

	// Use a context so that we can timeout any transactions that exceed a certain time limit.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	err := m.DB.QueryRowContext(ctx, query, start, end, roomID).Scan(&numRows)
	if err != nil {
		return false, err
	}

	return (numRows == 0), nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for a given date range.
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start time.Time, end time.Time) ([]models.Room, error) {
	query := `select r.id, r.room_name from rooms r where r.id not in (
				select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)`

	// Use a context so that we can timeout any transactions that exceed a certain time limit.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room

		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}
