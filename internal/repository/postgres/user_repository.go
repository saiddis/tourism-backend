package postgres

import (
	"database/sql"
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
	err := r.db.QueryRow(queries.GetUserByEmail, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
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
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}
