package postgres

import (
	"context"
	"database/sql"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository/queries"
)

type TourHighlightRepositoryPostgres struct {
	db *sql.DB
}

func NewTourHighlightRepositoryPostgres(db *sql.DB) *TourHighlightRepositoryPostgres {
	return &TourHighlightRepositoryPostgres{db: db}
}

func (r *TourHighlightRepositoryPostgres) Create(ctx context.Context, highlight *domain.TourHighlight) error {
	return r.db.QueryRowContext(ctx, queries.CreateTourHighlight,
		highlight.TourID,
		highlight.ImageURL,
		highlight.SortOrder,
	).Scan(&highlight.ID, &highlight.CreatedAt)
}

func (r *TourHighlightRepositoryPostgres) GetByTourID(ctx context.Context, tourID int) ([]*domain.TourHighlight, error) {
	rows, err := r.db.QueryContext(ctx, queries.GetTourHighlightsByTourID, tourID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	highlights := make([]*domain.TourHighlight, 0)
	for rows.Next() {
		var h domain.TourHighlight
		err := rows.Scan(&h.ID, &h.TourID, &h.ImageURL, &h.SortOrder, &h.CreatedAt)
		if err != nil {
			return nil, err
		}
		highlights = append(highlights, &h)
	}
	return highlights, nil
}

func (r *TourHighlightRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, queries.DeleteTourHighlight, id)
	return err
}
