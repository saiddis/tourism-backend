package repository

import (
	"context"
	"tourism-backend/internal/domain"
)

type TourRepository interface {
	Create(ctx context.Context, tour *domain.Tour) error
	GetByID(ctx context.Context, id int) (*domain.Tour, error)
	GetAll(ctx context.Context) ([]*domain.Tour, error)
	GetByDestinationID(ctx context.Context, destinationID int) ([]*domain.Tour, error)
	Update(ctx context.Context, tour *domain.Tour) error
	Delete(ctx context.Context, id int) error
	DecrementCapacity(ctx context.Context, id int) error
	IncrementCapacity(ctx context.Context, id int) error
}
