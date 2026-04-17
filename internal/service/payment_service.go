// internal/service/payment_service.go
package service

import (
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

func (s *PaymentService) Create(payment *domain.Payment, statusOverride ...domain.PaymentStatus) (*domain.Payment, error) {
	if len(statusOverride) > 0 {
		payment.Status = statusOverride[0]
	} else {
		payment.Status = domain.PaymentStatusPending
	}
	if payment.Currency == "" {
		payment.Currency = "TJS"
	}
	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *PaymentService) GetByID(id int) (*domain.Payment, error) {
	payment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errors.New("payment not found")
	}
	return payment, nil
}

func (s *PaymentService) GetByBookingID(bookingID int) (*domain.Payment, error) {
	return s.repo.GetByBookingID(bookingID)
}

func (s *PaymentService) UpdateStatus(id int, status domain.PaymentStatus) error {
	return s.repo.UpdateStatus(id, status)
}
