package service

import (
	"context"
	"errors"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository"
)

type PaymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) Create(ctx context.Context, payment *domain.Payment, statusOverride ...domain.PaymentStatus) (*domain.Payment, error) {
	if len(statusOverride) > 0 {
		payment.Status = statusOverride[0]
	} else {
		payment.Status = domain.PaymentStatusPending
	}
	if payment.Currency == "" {
		payment.Currency = "TJS"
	}
	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *PaymentService) GetByID(ctx context.Context, id int) (*domain.Payment, error) {
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errors.New("payment not found")
	}
	return payment, nil
}

func (s *PaymentService) GetByBookingID(ctx context.Context, bookingID int) (*domain.Payment, error) {
	return s.repo.GetByBookingID(ctx, bookingID)
}

func (s *PaymentService) UpdateStatus(ctx context.Context, id int, status domain.PaymentStatus) error {
	return s.repo.UpdateStatus(ctx, id, status)
}
