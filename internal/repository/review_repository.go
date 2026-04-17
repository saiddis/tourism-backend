package repository

import (
	"context"
	"tourism-backend/internal/domain"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *domain.Review) error
	GetByTourID(ctx context.Context, tourID int) ([]*domain.Review, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Review, error)
	Delete(ctx context.Context, id int) error
}
