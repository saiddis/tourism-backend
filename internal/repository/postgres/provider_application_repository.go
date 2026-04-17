package postgres

import (
	"context"
	"database/sql"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository/queries"
)

type ProviderApplicationRepositoryPostgres struct {
	db *sql.DB
}

func NewProviderApplicationRepository(db *sql.DB) *ProviderApplicationRepositoryPostgres {
	return &ProviderApplicationRepositoryPostgres{db: db}
}

func (r *ProviderApplicationRepositoryPostgres) Create(ctx context.Context, app *domain.ProviderApplication) error {
	return r.db.QueryRowContext(ctx, queries.CreateProviderApplication,
		app.UserID,
		app.Phone,
		app.ProviderType,
		app.InstagramURL,
		app.TelegramURL,
		app.FacebookURL,
		app.YearsExperience,
		app.Bio,
		domain.ApplicationStatusPending,
	).Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt)
}

func (r *ProviderApplicationRepositoryPostgres) GetByID(ctx context.Context, id int) (*domain.ProviderApplication, error) {
	app := &domain.ProviderApplication{}
	err := r.db.QueryRowContext(ctx, queries.GetProviderApplicationByID, id).Scan(
		&app.ID,
		&app.UserID,
		&app.Phone,
		&app.ProviderType,
		&app.InstagramURL,
		&app.TelegramURL,
		&app.FacebookURL,
		&app.YearsExperience,
		&app.Bio,
		&app.Status,
		&app.AdminNote,
		&app.CreatedAt,
		&app.UpdatedAt,
		&app.UserName,
		&app.UserEmail,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (r *ProviderApplicationRepositoryPostgres) GetByUserID(ctx context.Context, userID int) (*domain.ProviderApplication, error) {
	app := &domain.ProviderApplication{}
	err := r.db.QueryRowContext(ctx, queries.GetProviderApplicationByUserID, userID).Scan(
		&app.ID,
		&app.UserID,
		&app.Phone,
		&app.ProviderType,
		&app.InstagramURL,
		&app.TelegramURL,
		&app.FacebookURL,
		&app.YearsExperience,
		&app.Bio,
		&app.Status,
		&app.AdminNote,
		&app.CreatedAt,
		&app.UpdatedAt,
		&app.UserName,
		&app.UserEmail,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (r *ProviderApplicationRepositoryPostgres) GetAll(ctx context.Context) ([]*domain.ProviderApplication, error) {
	rows, err := r.db.QueryContext(ctx, queries.GetAllProviderApplications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*domain.ProviderApplication
	for rows.Next() {
		app := &domain.ProviderApplication{}
		err := rows.Scan(
			&app.ID,
			&app.UserID,
			&app.Phone,
			&app.ProviderType,
			&app.InstagramURL,
			&app.TelegramURL,
			&app.FacebookURL,
			&app.YearsExperience,
			&app.Bio,
			&app.Status,
			&app.AdminNote,
			&app.CreatedAt,
			&app.UpdatedAt,
			&app.UserName,
			&app.UserEmail,
		)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (r *ProviderApplicationRepositoryPostgres) GetPending(ctx context.Context) ([]*domain.ProviderApplication, error) {
	rows, err := r.db.QueryContext(ctx, queries.GetPendingProviderApplications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*domain.ProviderApplication
	for rows.Next() {
		app := &domain.ProviderApplication{}
		err := rows.Scan(
			&app.ID,
			&app.UserID,
			&app.Phone,
			&app.ProviderType,
			&app.InstagramURL,
			&app.TelegramURL,
			&app.FacebookURL,
			&app.YearsExperience,
			&app.Bio,
			&app.Status,
			&app.AdminNote,
			&app.CreatedAt,
			&app.UpdatedAt,
			&app.UserName,
			&app.UserEmail,
		)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (r *ProviderApplicationRepositoryPostgres) UpdateStatus(ctx context.Context, id int, status domain.ApplicationStatus, adminNote *string) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateProviderApplicationStatus, id, status, adminNote)
	return err
}
