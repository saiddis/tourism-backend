package postgres

import (
	"context"
	"database/sql"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository/queries"
)

type ProviderRepositoryPostgres struct {
	db *sql.DB
}

func NewProviderRepository(db *sql.DB) *ProviderRepositoryPostgres {
	return &ProviderRepositoryPostgres{db: db}
}

func (r *ProviderRepositoryPostgres) Create(ctx context.Context, provider *domain.Provider) error {
	return r.db.QueryRowContext(ctx, queries.CreateProvider,
		provider.UserID,
		provider.Phone,
		provider.ProviderType,
		provider.InstagramURL,
		provider.TelegramURL,
		provider.FacebookURL,
		provider.YearsExperience,
		provider.Bio,
		true,
	).Scan(&provider.ID, &provider.CreatedAt)
}

func (r *ProviderRepositoryPostgres) GetByID(ctx context.Context, id int) (*domain.Provider, error) {
	provider := &domain.Provider{}
	err := r.db.QueryRowContext(ctx, queries.GetProviderByID, id).Scan(
		&provider.ID,
		&provider.UserID,
		&provider.Phone,
		&provider.ProviderType,
		&provider.InstagramURL,
		&provider.TelegramURL,
		&provider.FacebookURL,
		&provider.YearsExperience,
		&provider.Bio,
		&provider.Active,
		&provider.CreatedAt,
		&provider.UserName,
		&provider.UserEmail,
		&provider.UserAvatarURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	provider.Active = true
	return provider, nil
}

func (r *ProviderRepositoryPostgres) GetByUserID(ctx context.Context, userID int) (*domain.Provider, error) {
	provider := &domain.Provider{}
	err := r.db.QueryRowContext(ctx, queries.GetProviderByUserID, userID).Scan(
		&provider.ID,
		&provider.UserID,
		&provider.Phone,
		&provider.ProviderType,
		&provider.InstagramURL,
		&provider.TelegramURL,
		&provider.FacebookURL,
		&provider.YearsExperience,
		&provider.Bio,
		&provider.Active,
		&provider.CreatedAt,
		&provider.UserName,
		&provider.UserEmail,
		&provider.UserAvatarURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (r *ProviderRepositoryPostgres) GetAll(ctx context.Context) ([]*domain.Provider, error) {
	rows, err := r.db.QueryContext(ctx, queries.GetAllProviders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*domain.Provider
	for rows.Next() {
		provider := &domain.Provider{}
		err := rows.Scan(
			&provider.ID,
			&provider.UserID,
			&provider.Phone,
			&provider.ProviderType,
			&provider.InstagramURL,
			&provider.TelegramURL,
			&provider.FacebookURL,
			&provider.YearsExperience,
			&provider.Bio,
			&provider.Active,
			&provider.CreatedAt,
			&provider.UserName,
			&provider.UserEmail,
			&provider.UserAvatarURL,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func (r *ProviderRepositoryPostgres) GetActive(ctx context.Context) ([]*domain.Provider, error) {
	rows, err := r.db.QueryContext(ctx, queries.GetActiveProviders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*domain.Provider
	for rows.Next() {
		provider := &domain.Provider{}
		err := rows.Scan(
			&provider.ID,
			&provider.UserID,
			&provider.Phone,
			&provider.ProviderType,
			&provider.InstagramURL,
			&provider.TelegramURL,
			&provider.FacebookURL,
			&provider.YearsExperience,
			&provider.Bio,
			&provider.Active,
			&provider.CreatedAt,
			&provider.UserName,
			&provider.UserEmail,
			&provider.UserAvatarURL,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func (r *ProviderRepositoryPostgres) Update(ctx context.Context, provider *domain.Provider) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateProvider,
		provider.ID,
		provider.Phone,
		provider.ProviderType,
		provider.InstagramURL,
		provider.TelegramURL,
		provider.FacebookURL,
		provider.YearsExperience,
		provider.Bio,
	)
	return err
}

func (r *ProviderRepositoryPostgres) UpdateActive(ctx context.Context, id int, active bool) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateProviderActive, id, active)
	return err
}
