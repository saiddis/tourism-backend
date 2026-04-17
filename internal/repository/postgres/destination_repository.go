package postgres

import (
	"context"
	"database/sql"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository/queries"
)

type DestinationRepositoryPostgres struct {
	db *sql.DB
}

func NewDestinationRepository(db *sql.DB) *DestinationRepositoryPostgres {
	return &DestinationRepositoryPostgres{db: db}
}

func (r *DestinationRepositoryPostgres) Create(ctx context.Context, destination *domain.Destination) error {
	return r.db.QueryRowContext(ctx, queries.CreateDestination,
		destination.Name,
		destination.Description,
		destination.ImageURL,
	).Scan(&destination.ID, &destination.ImageURL, &destination.CreatedAt)
}

func (r *DestinationRepositoryPostgres) GetByID(ctx context.Context, id int) (*domain.Destination, error) {
	var destination domain.Destination
	err := scanDestination(r.db.QueryRowContext(ctx, queries.GetDestinationByID, id), &destination)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &destination, err
}

func (r *DestinationRepositoryPostgres) GetAll(ctx context.Context) ([]*domain.Destination, error) {
	destinations := make([]*domain.Destination, 0)
	rows, err := r.db.QueryContext(ctx, queries.GetAllDestinations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var destination domain.Destination
		err = scanDestination(rows, &destination)
		if err != nil {
			return nil, err
		}
		destinations = append(destinations, &destination)
	}
	return destinations, nil
}

func (r *DestinationRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, queries.DeleteDestination, id)
	return err
}

func scanDestination(scanner rowScanner, destination *domain.Destination) error {
	var description sql.NullString
	var imageURL sql.NullString
	err := scanner.Scan(
		&destination.ID,
		&destination.Name,
		&description,
		&imageURL,
		&destination.CreatedAt,
	)
	if err != nil {
		return err
	}
	destination.Description = description.String
	destination.ImageURL = imageURL.String
	return nil
}
