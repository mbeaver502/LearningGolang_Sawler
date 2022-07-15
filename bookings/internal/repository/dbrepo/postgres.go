package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

const DEFAULT_TIMEOUT = 3 * time.Second

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
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
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
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
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
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
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
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	var rooms []models.Room

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

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

// GetRoomByID gets a room by the given id.
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	// Use a context so that we can timeout any transactions that exceed a certain time limit.
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	var room models.Room
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}

// GetUserByID returns a user by id.
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	// Use a context so that we can timeout any transactions that exceed a certain time limit.
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
				from users where id = $1`

	var u models.User
	err := m.DB.QueryRowContext(ctx, query,
		id,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates a user in the database.
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	// Use a context so that we can timeout any transactions that exceed a certain time limit.
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user.
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	// Use a context so that we can timeout any transactions that exceed a certain time limit.
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	// the ID of the authenticated user, if successful
	var id int
	var hashedPassword string

	query := `select id, password from users where email = $1`
	err := m.DB.QueryRowContext(ctx, query,
		email,
	).Scan(
		&id,
		&hashedPassword,
	)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// AllReservations returns a slice of all reservations.
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	// Use a context so that we can timeout any transactions that exceed a certain time limit.
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, 
			r.first_name, 
			r.last_name, 
			r.email, 
			r.phone, 
			r.start_date, 
			r.end_date, 
			r.room_id, 
			r.created_at, 
			r.updated_at, 
			r.processed,
			rm.id, 
			rm.room_name
		from reservations r left join rooms rm on (r.room_id = rm.id)
		order by r.start_date, r.end_date`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation

		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// AllNewReservations returns a slice of new reservations.
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	// Use a context so that we can timeout any transactions that exceed a certain time limit.
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, 
			r.first_name, 
			r.last_name, 
			r.email, 
			r.phone, 
			r.start_date, 
			r.end_date, 
			r.room_id, 
			r.created_at, 
			r.updated_at, 
			rm.id, 
			rm.room_name
		from reservations r left join rooms rm on (r.room_id = rm.id)
		where processed = 0
		order by r.start_date, r.end_date`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation

		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}
