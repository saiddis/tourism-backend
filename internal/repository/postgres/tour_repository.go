package postgres

import (
	"database/sql"
	"errors"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository/queries"
)

type TourRepositoryPostgres struct {
	db *sql.DB
}

func NewTourRepositoryPostgres(db *sql.DB) *TourRepositoryPostgres {
	return &TourRepositoryPostgres{db: db}
}

func (r *TourRepositoryPostgres) Create(tour *domain.Tour) error {
	return r.db.QueryRow(queries.CreateTour, tour.DestinationID, tour.Name, tour.Description, tour.Price, tour.StartDate, tour.EndDate, tour.Capacity).Scan(&tour.ID, &tour.CreatedAt)
}

func (r *TourRepositoryPostgres) GetByID(id int) (*domain.Tour, error) {
	var tour domain.Tour
	err := scanTour(r.db.QueryRow(queries.GetTourByID, id), &tour)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &tour, err
}

func (r *TourRepositoryPostgres) GetAll() ([]*domain.Tour, error) {
	tours := make([]*domain.Tour, 0)
	rows, err := r.db.Query(queries.GetAllTours)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tour domain.Tour
		err = scanTour(rows, &tour)
		if err != nil {
			return nil, err
		}
		tours = append(tours, &tour)
	}
	return tours, nil
}

func (r *TourRepositoryPostgres) GetByDestinationID(destinationID int) ([]*domain.Tour, error) {
	tours := make([]*domain.Tour, 0)
	rows, err := r.db.Query(queries.GetToursByDestinationID, destinationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tour domain.Tour
		err = scanTour(rows, &tour)
		if err != nil {
			return nil, err
		}
		tours = append(tours, &tour)
	}
	return tours, nil
}

func (r *TourRepositoryPostgres) Update(tour *domain.Tour) error {
	return r.db.QueryRow(queries.UpdateTour,
		tour.DestinationID,
		tour.Name,
		tour.Description,
		tour.Price,
		tour.StartDate,
		tour.EndDate,
		tour.Capacity,
		tour.ID,
	).Scan(&tour.CreatedAt)
}

func (r *TourRepositoryPostgres) Delete(id int) error {
	_, err := r.db.Exec(queries.DeleteTour, id)
	return err
}

func (r *TourRepositoryPostgres) DecrementCapacity(id int) error {
	result, err := r.db.Exec(queries.DecrementTourCapacity, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no capacity available")
	}
	return nil
}

func (r *TourRepositoryPostgres) IncrementCapacity(id int) error {
	_, err := r.db.Exec(queries.IncrementTourCapacity, id)
	return err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTour(scanner rowScanner, tour *domain.Tour) error {
	var description sql.NullString
	var destinationDescription sql.NullString
	var destinationImageURL sql.NullString
	var destinationCreatedAt sql.NullTime
	var joinedDestinationID int
	err := scanner.Scan(
		&tour.ID,
		&tour.DestinationID,
		&tour.Name,
		&description,
		&tour.Price,
		&tour.StartDate,
		&tour.EndDate,
		&tour.Capacity,
		&tour.CreatedAt,
		&joinedDestinationID,
		&tour.DestinationName,
		&destinationDescription,
		&destinationImageURL,
		&destinationCreatedAt,
	)
	if err != nil {
		return err
	}
	tour.Description = description.String
	tour.DestinationDescription = destinationDescription.String
	tour.DestinationImageURL = destinationImageURL.String
	if tour.DestinationID == 0 {
		tour.DestinationID = joinedDestinationID
	}
	return nil
}
