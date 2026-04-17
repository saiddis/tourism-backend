package repository

import (
	"context"
	"tourism-backend/internal/domain"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	GetByID(ctx context.Context, id int) (*domain.Booking, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Booking, error)
	GetAll(ctx context.Context) ([]*domain.Booking, error)
	UpdateStatus(ctx context.Context, id int, status domain.BookingStatus) error
	Delete(ctx context.Context, id int) error
}
