package repository

import (
	"time"

	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDates(start time.Time, end time.Time, roomID int) (bool, error)
}
