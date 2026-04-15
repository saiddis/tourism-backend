package postgres

import (
	"database/sql"
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

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTour(scanner rowScanner, tour *domain.Tour) error {
	destination := &domain.Destination{}
	err := scanner.Scan(
		&tour.ID,
		&tour.DestinationID,
		&tour.Name,
		&tour.Description,
		&tour.Price,
		&tour.StartDate,
		&tour.EndDate,
		&tour.Capacity,
		&tour.CreatedAt,
		&destination.ID,
		&destination.Name,
		&destination.Description,
		&destination.ImageURL,
		&destination.CreatedAt,
	)
	if err != nil {
		return err
	}
	tour.Destination = destination
	return nil
}
