package service

import (
	"errors"
	"strings"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(name, email, password string) (*domain.User, error) {
	exist, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, errors.New("user with this email already exists")
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         domain.RoleClient,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(email, password string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	passwordMatches, passwordWasUpgraded, err := s.checkPassword(user, password)
	if err != nil {
		return nil, err
	}
	if !passwordMatches {
		return nil, errors.New("invalid password")
	}

	if passwordWasUpgraded {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			return nil, err
		}
		if err := s.repo.UpdatePasswordHash(user.ID, hashedPassword); err != nil {
			return nil, err
		}
		user.PasswordHash = hashedPassword
	}

	return user, nil
}

func (s *UserService) GetByID(id int) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) GetAll() ([]*domain.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) Update(id int, name, email string) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Name = name
	user.Email = email

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateAvatarURL(id int, url string) error {
	return s.repo.UpdateAvatarURL(id, url)
}

func (s *UserService) Deposit(id int, amount float64) (*domain.User, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	newBalance := user.Balance + amount
	if err := s.repo.UpdateBalance(id, newBalance); err != nil {
		return nil, err
	}

	user.Balance = newBalance
	return user, nil
}

func (s *UserService) DeductBalance(userID int, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return s.repo.DeductBalance(userID, amount)
}

func (s *UserService) RefundBalance(userID int, amount float64) (*domain.User, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	newBalance := user.Balance + amount
	if err := s.repo.UpdateBalance(userID, newBalance); err != nil {
		return nil, err
	}

	user.Balance = newBalance
	return user, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *UserService) checkPassword(user *domain.User, password string) (bool, bool, error) {
	if isBcryptHash(user.PasswordHash) {
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err == nil {
			return true, false, nil
		}
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, false, nil
		}
		return false, false, err
	}

	return user.PasswordHash == password, user.PasswordHash == password, nil
}

func isBcryptHash(passwordHash string) bool {
	return strings.HasPrefix(passwordHash, "$2a$") ||
		strings.HasPrefix(passwordHash, "$2b$") ||
		strings.HasPrefix(passwordHash, "$2y$")
}
