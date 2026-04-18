package repository

import (
	"context"
	"time"
	"tourism-backend/internal/domain"
)

type PendingBookingWithUser struct {
	BookingID int
	UserID    int
	TourID    int
	Status    domain.BookingStatus
	CreatedAt time.Time
	Email     string
	Name      string
}

type TourConfirmationInfo struct {
	TourID       int
	Name         string
	StartDate    time.Time
	Capacity     int
	PendingCount int
}

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	GetByID(ctx context.Context, id int) (*domain.Booking, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Booking, error)
	GetCompletedBookingsByUserID(ctx context.Context, userID int) ([]*domain.Booking, error)
	GetAll(ctx context.Context) ([]*domain.Booking, error)
	UpdateStatus(ctx context.Context, id int, status domain.BookingStatus) error
	Delete(ctx context.Context, id int) error
	GetPendingBookingsCountByTourID(ctx context.Context, tourID int) (int, error)
	GetActiveBookingsCount(ctx context.Context, tourID int) (int, error)
	GetPendingBookingsByTourID(ctx context.Context, tourID int) ([]PendingBookingWithUser, error)
	ConfirmPendingBookingsForTour(ctx context.Context, tourID int) error
	GetToursNeedingConfirmation(ctx context.Context) ([]TourConfirmationInfo, error)
	MarkBookingsCompleted(ctx context.Context) error
}
