package service

import (
	"context"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository"
)

type TourHighlightService struct {
	repo repository.TourHighlightRepository
}

func NewTourHighlightService(repo repository.TourHighlightRepository) *TourHighlightService {
	return &TourHighlightService{repo: repo}
}

func (s *TourHighlightService) Create(ctx context.Context, highlight *domain.TourHighlight) error {
	return s.repo.Create(ctx, highlight)
}

func (s *TourHighlightService) GetByTourID(ctx context.Context, tourID int) ([]*domain.TourHighlight, error) {
	return s.repo.GetByTourID(ctx, tourID)
}

func (s *TourHighlightService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
