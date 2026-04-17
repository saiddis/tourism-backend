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
	UpdateAvatarURL(id int, url string) error
	UpdateBalance(id int, balance float64) error
	DeductBalance(id int, amount float64) error
}
