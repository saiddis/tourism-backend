package repository

import (
	"context"
	"tourism-backend/internal/domain"
)

type DestinationRepository interface {
	Create(ctx context.Context, destination *domain.Destination) error
	GetByID(ctx context.Context, id int) (*domain.Destination, error)
	GetAll(ctx context.Context) ([]*domain.Destination, error)
	Delete(ctx context.Context, id int) error
}
