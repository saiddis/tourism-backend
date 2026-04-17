package repository

import (
	"context"
	"tourism-backend/internal/domain"
)

type TourHighlightRepository interface {
	Create(ctx context.Context, highlight *domain.TourHighlight) error
	GetByTourID(ctx context.Context, tourID int) ([]*domain.TourHighlight, error)
	Delete(ctx context.Context, id int) error
}
