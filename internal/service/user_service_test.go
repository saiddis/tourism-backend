package service

import (
	"testing"
	"tourism-backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type userRepoStub struct {
	usersByID           map[int]*domain.User
	usersByEmail        map[string]*domain.User
	createdUser         *domain.User
	updatedPasswordID   int
	updatedPasswordHash string
}

func newUserRepoStub(users ...*domain.User) *userRepoStub {
	repo := &userRepoStub{
		usersByID:    make(map[int]*domain.User),
		usersByEmail: make(map[string]*domain.User),
	}
	for _, user := range users {
		repo.usersByID[user.ID] = user
		repo.usersByEmail[user.Email] = user
	}
	return repo
}

func (r *userRepoStub) Create(user *domain.User) error {
	r.createdUser = user
	if user.ID == 0 {
		user.ID = len(r.usersByID) + 1
	}
	r.usersByID[user.ID] = user
	r.usersByEmail[user.Email] = user
	return nil
}

func (r *userRepoStub) GetByID(id int) (*domain.User, error) {
	return r.usersByID[id], nil
}

func (r *userRepoStub) GetByEmail(email string) (*domain.User, error) {
	return r.usersByEmail[email], nil
}

func (r *userRepoStub) GetAll() ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(r.usersByID))
	for _, user := range r.usersByID {
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepoStub) Update(user *domain.User) error {
	r.usersByID[user.ID] = user
	r.usersByEmail[user.Email] = user
	return nil
}

func (r *userRepoStub) UpdatePasswordHash(id int, passwordHash string) error {
	r.updatedPasswordID = id
	r.updatedPasswordHash = passwordHash
	if user, ok := r.usersByID[id]; ok {
		user.PasswordHash = passwordHash
		r.usersByEmail[user.Email] = user
	}
	return nil
}

func (r *userRepoStub) Delete(id int) error {
	delete(r.usersByID, id)
	return nil
}

func TestRegisterHashesPassword(t *testing.T) {
	repo := newUserRepoStub()
	service := NewUserService(repo)

	user, err := service.Register("Alice", "alice@example.com", "secret123")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if repo.createdUser == nil {
		t.Fatal("expected Create() to be called")
	}
	if repo.createdUser.PasswordHash == "secret123" {
		t.Fatal("password was stored as plain text")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.createdUser.PasswordHash), []byte("secret123")); err != nil {
		t.Fatalf("stored password hash is invalid: %v", err)
	}
	if user.PasswordHash == "secret123" {
		t.Fatal("returned user still has a plain text password")
	}
}

func TestLoginAcceptsHashedPassword(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	user := &domain.User{
		ID:           1,
		Email:        "alice@example.com",
		PasswordHash: string(passwordHash),
	}

	repo := newUserRepoStub(user)
	service := NewUserService(repo)

	loggedInUser, err := service.Login("alice@example.com", "secret123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if loggedInUser.ID != user.ID {
		t.Fatalf("Login() returned user ID %d, want %d", loggedInUser.ID, user.ID)
	}
	if repo.updatedPasswordHash != "" {
		t.Fatal("hashed password should not be rewritten")
	}
}

func TestLoginUpgradesLegacyPlainTextPassword(t *testing.T) {
	user := &domain.User{
		ID:           7,
		Email:        "legacy@example.com",
		PasswordHash: "legacy-pass",
	}

	repo := newUserRepoStub(user)
	service := NewUserService(repo)

	_, err := service.Login("legacy@example.com", "legacy-pass")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if repo.updatedPasswordID != user.ID {
		t.Fatalf("UpdatePasswordHash() called for user %d, want %d", repo.updatedPasswordID, user.ID)
	}
	if repo.updatedPasswordHash == "" || repo.updatedPasswordHash == "legacy-pass" {
		t.Fatal("legacy password was not upgraded to a bcrypt hash")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.updatedPasswordHash), []byte("legacy-pass")); err != nil {
		t.Fatalf("upgraded password hash is invalid: %v", err)
	}
}
