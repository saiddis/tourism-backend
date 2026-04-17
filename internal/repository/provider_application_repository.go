package repository

import (
	"context"
	"tourism-backend/internal/domain"
)

type ProviderApplicationRepository interface {
	Create(ctx context.Context, app *domain.ProviderApplication) error
	GetByID(ctx context.Context, id int) (*domain.ProviderApplication, error)
	GetByUserID(ctx context.Context, userID int) (*domain.ProviderApplication, error)
	GetByToken(ctx context.Context, token string) (*domain.ProviderApplication, error)
	GetAll(ctx context.Context) ([]*domain.ProviderApplication, error)
	GetPending(ctx context.Context) ([]*domain.ProviderApplication, error)
	UpdateStatus(ctx context.Context, id int, status domain.ApplicationStatus, adminNote *string) error
}
