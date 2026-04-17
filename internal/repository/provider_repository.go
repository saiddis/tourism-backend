package repository

import (
	"context"
	"tourism-backend/internal/domain"
)

type ProviderRepository interface {
	Create(ctx context.Context, provider *domain.Provider) error
	GetByID(ctx context.Context, id int) (*domain.Provider, error)
	GetByUserID(ctx context.Context, userID int) (*domain.Provider, error)
	GetAll(ctx context.Context) ([]*domain.Provider, error)
	GetActive(ctx context.Context) ([]*domain.Provider, error)
	Update(ctx context.Context, provider *domain.Provider) error
	UpdateActive(ctx context.Context, id int, active bool) error
}
