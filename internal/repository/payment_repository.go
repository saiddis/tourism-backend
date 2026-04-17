package repository

import (
	"context"
	"tourism-backend/internal/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, id int) (*domain.Payment, error)
	GetByBookingID(ctx context.Context, bookingID int) (*domain.Payment, error)
	UpdateStatus(ctx context.Context, id int, status domain.PaymentStatus) error
}
