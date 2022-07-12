package repository

import "github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}
