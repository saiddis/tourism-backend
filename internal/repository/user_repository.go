package repository

import (
	"context"
	"tourism-backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePasswordHash(ctx context.Context, id int, passwordHash string) error
	Delete(ctx context.Context, id int) error
	UpdateAvatarURL(ctx context.Context, id int, url string) error
	UpdateBalance(ctx context.Context, id int, balance float64) error
	DeductBalance(ctx context.Context, id int, amount float64) error
}
