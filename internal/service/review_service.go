// internal/service/review_service.go
package service

import (
	"context"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository"
)

type ReviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) Create(ctx context.Context, review *domain.Review) (*domain.Review, error) {
	if err := s.repo.Create(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) GetByTourID(ctx context.Context, tourID int) ([]*domain.Review, error) {
	return s.repo.GetByTourID(ctx, tourID)
}

func (s *ReviewService) GetByUserID(ctx context.Context, userID int) ([]*domain.Review, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *ReviewService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
