package service

import (
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

func (s *TourService) Create(tour *domain.Tour) (*domain.Tour, error) {
	if err := s.repo.Create(tour); err != nil {
		return nil, err
	}
	return tour, nil
}

func (s *TourService) GetByID(id int) (*domain.Tour, error) {
	tour, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if tour == nil {
		return nil, errors.New("tour not found")
	}
	return tour, nil
}

func (s *TourService) GetAll() ([]*domain.Tour, error) {
	return s.repo.GetAll()
}

func (s *TourService) GetByDestinationID(destinationID int) ([]*domain.Tour, error) {
	return s.repo.GetByDestinationID(destinationID)
}

func (s *TourService) Update(tour *domain.Tour) (*domain.Tour, error) {
	if err := s.repo.Update(tour); err != nil {
		return nil, err
	}
	return tour, nil
}

func (s *TourService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *TourService) DecrementCapacity(id int) error {
	return s.repo.DecrementCapacity(id)
}

func (s *TourService) IncrementCapacity(id int) error {
	return s.repo.IncrementCapacity(id)
}
