package postgres

import (
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

func (r *UserRepositoryPostgres) Create(user *domain.User) error {
	return r.db.QueryRow(queries.CreateUser, user.Name, user.Email, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepositoryPostgres) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(queries.GetUserByEmail, email).Scan(
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

func (r *UserRepositoryPostgres) GetAll() ([]*domain.User, error) {
	rows, err := r.db.Query(queries.GetAllUsers)
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

func (r *UserRepositoryPostgres) Update(user *domain.User) error {
	_, err := r.db.Exec(queries.UpdateUser, user.Name, user.Email, user.ID)
	return err
}

func (r *UserRepositoryPostgres) UpdatePasswordHash(id int, passwordHash string) error {
	_, err := r.db.Exec(queries.UpdateUserPassword, passwordHash, id)
	return err
}

func (r *UserRepositoryPostgres) Delete(id int) error {
	_, err := r.db.Exec(queries.DeleteUser, id)
	return err
}

func (r *UserRepositoryPostgres) GetByID(id int) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(queries.GetUserByID, id).Scan(
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

func (r *UserRepositoryPostgres) UpdateAvatarURL(id int, url string) error {
	_, err := r.db.Exec(queries.UpdateUserAvatarURL, url, id)
	return err
}

func (r *UserRepositoryPostgres) UpdateBalance(id int, balance float64) error {
	_, err := r.db.Exec(queries.UpdateUserBalance, balance, id)
	return err
}

func (r *UserRepositoryPostgres) DeductBalance(id int, amount float64) error {
	result, err := r.db.Exec(queries.DeductUserBalance, amount, id)
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
