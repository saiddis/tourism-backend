package postgres

import (
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

func (r *DestinationRepositoryPostgres) Create(destination *domain.Destination) error {
	return r.db.QueryRow(queries.CreateDestination,
		destination.Name,
		destination.Description,
		destination.ImageURL,
	).Scan(&destination.ID, &destination.ImageURL, &destination.CreatedAt)
}

func (r *DestinationRepositoryPostgres) GetByID(id int) (*domain.Destination, error) {
	var destination domain.Destination
	err := scanDestination(r.db.QueryRow(queries.GetDestinationByID, id), &destination)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &destination, err
}

func (r *DestinationRepositoryPostgres) GetAll() ([]*domain.Destination, error) {
	destinations := make([]*domain.Destination, 0)
	rows, err := r.db.Query(queries.GetAllDestinations)
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

func (r *DestinationRepositoryPostgres) Delete(id int) error {
	_, err := r.db.Exec(queries.DeleteDestination, id)
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
