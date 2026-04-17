package postgres

import (
	"context"
	"database/sql"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository/queries"
)

type ReviewRepositoryPostgres struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepositoryPostgres {
	return &ReviewRepositoryPostgres{db: db}
}

func (r *ReviewRepositoryPostgres) Create(ctx context.Context, review *domain.Review) error {
	return r.db.QueryRowContext(ctx, queries.CreateReview,
		review.UserID,
		review.TourID,
		review.Rating,
		review.Comment,
	).Scan(&review.ID, &review.CreatedAt)
}

func (r *ReviewRepositoryPostgres) GetByTourID(ctx context.Context, tourID int) ([]*domain.Review, error) {
	reviews := make([]*domain.Review, 0)
	rows, err := r.db.QueryContext(ctx, queries.GetReviewsByTourID, tourID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var review domain.Review
		err = rows.Scan(
			&review.ID,
			&review.UserID,
			&review.TourID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}
	return reviews, nil
}

func (r *ReviewRepositoryPostgres) GetByUserID(ctx context.Context, userID int) ([]*domain.Review, error) {
	reviews := make([]*domain.Review, 0)
	rows, err := r.db.QueryContext(ctx, queries.GetReviewsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var review domain.Review
		err = rows.Scan(
			&review.ID,
			&review.UserID,
			&review.TourID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}
	return reviews, nil
}

func (r *ReviewRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, queries.DeleteReview, id)
	return err
}
