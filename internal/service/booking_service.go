package service

import (
	"context"
	"errors"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository"
)

type BookingService struct {
	repo repository.BookingRepository
}

func NewBookingService(repo repository.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) Create(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	if booking.Status == "" {
		booking.Status = domain.BookingStatusPending
	}
	if err := s.repo.Create(ctx, booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *BookingService) GetByID(ctx context.Context, id int) (*domain.Booking, error) {
	booking, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}
	return booking, nil
}

func (s *BookingService) GetByUserID(ctx context.Context, userID int) ([]*domain.Booking, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *BookingService) GetAll(ctx context.Context) ([]*domain.Booking, error) {
	return s.repo.GetAll(ctx)
}

func (s *BookingService) UpdateStatus(ctx context.Context, id int, status domain.BookingStatus) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *BookingService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
