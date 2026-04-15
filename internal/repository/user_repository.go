package repository

import "tourism-backend/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id int) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Update(user *domain.User) error
	UpdatePasswordHash(id int, passwordHash string) error
	Delete(id int) error
}
