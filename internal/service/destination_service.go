package service

import (
	"context"
	"errors"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository"
)

type DestinationService struct {
	repo repository.DestinationRepository
}

func NewDestinationService(repo repository.DestinationRepository) *DestinationService {
	return &DestinationService{repo: repo}
}

func (s *DestinationService) Create(ctx context.Context, destination *domain.Destination) (*domain.Destination, error) {
	if err := s.repo.Create(ctx, destination); err != nil {
		return nil, err
	}
	return destination, nil
}

func (s *DestinationService) GetByID(ctx context.Context, id int) (*domain.Destination, error) {
	destination, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if destination == nil {
		return nil, errors.New("destination not found")
	}
	return destination, nil
}

func (s *DestinationService) GetAll(ctx context.Context) ([]*domain.Destination, error) {
	return s.repo.GetAll(ctx)
}

func (s *DestinationService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
