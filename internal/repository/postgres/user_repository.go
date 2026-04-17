package postgres

import (
	"context"
	"database/sql"
	"errors"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository/queries"
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) Create(ctx context.Context, user *domain.User) error {
	return r.db.QueryRowContext(ctx, queries.CreateUser, user.Name, user.Email, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepositoryPostgres) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, queries.GetUserByEmail, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.AvatarURL,
		&user.Balance,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepositoryPostgres) GetAll(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.QueryContext(ctx, queries.GetAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*domain.User, 0)
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.AvatarURL,
			&user.Balance,
			&user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepositoryPostgres) Update(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateUser, user.Name, user.Email, user.ID)
	return err
}

func (r *UserRepositoryPostgres) UpdatePasswordHash(ctx context.Context, id int, passwordHash string) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateUserPassword, passwordHash, id)
	return err
}

func (r *UserRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, queries.DeleteUser, id)
	return err
}

func (r *UserRepositoryPostgres) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, queries.GetUserByID, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.AvatarURL,
		&user.Balance,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepositoryPostgres) UpdateAvatarURL(ctx context.Context, id int, url string) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateUserAvatarURL, url, id)
	return err
}

func (r *UserRepositoryPostgres) UpdateBalance(ctx context.Context, id int, balance float64) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateUserBalance, balance, id)
	return err
}

func (r *UserRepositoryPostgres) DeductBalance(ctx context.Context, id int, amount float64) error {
	result, err := r.db.ExecContext(ctx, queries.DeductUserBalance, amount, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("insufficient balance")
	}
	return nil
}

func (r *UserRepositoryPostgres) UpdateRole(ctx context.Context, id int, role domain.UserRole) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateUserRole, role, id)
	return err
}
