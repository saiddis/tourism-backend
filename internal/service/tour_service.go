package service

import (
	"context"
	"errors"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository"
)

type TourService struct {
	repo repository.TourRepository
}

func NewTourService(repo repository.TourRepository) *TourService {
	return &TourService{repo: repo}
}

func (s *TourService) Create(ctx context.Context, tour *domain.Tour) (*domain.Tour, error) {
	if err := s.repo.Create(ctx, tour); err != nil {
		return nil, err
	}
	return tour, nil
}

func (s *TourService) GetByID(ctx context.Context, id int) (*domain.Tour, error) {
	tour, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tour == nil {
		return nil, errors.New("tour not found")
	}
	return tour, nil
}

func (s *TourService) GetAll(ctx context.Context) ([]*domain.Tour, error) {
	return s.repo.GetAll(ctx)
}

func (s *TourService) GetByDestinationID(ctx context.Context, destinationID int) ([]*domain.Tour, error) {
	return s.repo.GetByDestinationID(ctx, destinationID)
}

func (s *TourService) Update(ctx context.Context, tour *domain.Tour) (*domain.Tour, error) {
	if err := s.repo.Update(ctx, tour); err != nil {
		return nil, err
	}
	return tour, nil
}

func (s *TourService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *TourService) DecrementCapacity(ctx context.Context, id int) error {
	return s.repo.DecrementCapacity(ctx, id)
}

func (s *TourService) IncrementCapacity(ctx context.Context, id int) error {
	return s.repo.IncrementCapacity(ctx, id)
}

func (s *TourService) CalculateRemainingSpots(ctx context.Context, tours []*domain.Tour) error {
	for _, tour := range tours {
		remaining, err := s.repo.GetRemainingSpots(ctx, tour.ID)
		if err != nil {
			return err
		}
		tour.RemainingSpots = remaining
	}
	return nil
}

func (s *TourService) GetRemainingSpots(ctx context.Context, tourID int) (int, error) {
	return s.repo.GetRemainingSpots(ctx, tourID)
}
