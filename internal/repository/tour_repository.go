package repository

import "tourism-backend/internal/domain"

type TourRepository interface {
	Create(tour *domain.Tour) error
	GetByID(id int) (*domain.Tour, error)
	GetAll() ([]*domain.Tour, error)
	GetByDestinationID(destinationID int) ([]*domain.Tour, error)
	Update(tour *domain.Tour) error
	Delete(id int) error
	DecrementCapacity(id int) error
	IncrementCapacity(id int) error
}
