package postgres

import (
	"context"
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

func (r *TourRepositoryPostgres) Create(ctx context.Context, tour *domain.Tour) error {
	return r.db.QueryRowContext(ctx, queries.CreateTour, tour.DestinationID, tour.ProviderID, tour.Name, tour.Description, tour.Price, tour.StartDate, tour.EndDate, tour.Capacity).Scan(&tour.ID, &tour.CreatedAt)
}

func (r *TourRepositoryPostgres) GetByID(ctx context.Context, id int) (*domain.Tour, error) {
	var tour domain.Tour
	err := scanTour(r.db.QueryRowContext(ctx, queries.GetTourByID, id), &tour)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	tour.Highlights, err = r.GetHighlightsByTourID(ctx, id)
	return &tour, err
}

func (r *TourRepositoryPostgres) GetAll(ctx context.Context) ([]*domain.Tour, error) {
	tours := make([]*domain.Tour, 0)
	rows, err := r.db.QueryContext(ctx, queries.GetAllTours)
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

	// Fetch highlights for all tours
	if len(tours) > 0 {
		tourIDs := make([]int, len(tours))
		for i, t := range tours {
			tourIDs[i] = t.ID
		}
		highlightsMap, err := r.GetHighlightsByTourIDs(ctx, tourIDs)
		if err != nil {
			return nil, err
		}
		for _, tour := range tours {
			tour.Highlights = highlightsMap[tour.ID]
		}
	}

	return tours, nil
}

func (r *TourRepositoryPostgres) GetByDestinationID(ctx context.Context, destinationID int) ([]*domain.Tour, error) {
	tours := make([]*domain.Tour, 0)
	rows, err := r.db.QueryContext(ctx, queries.GetToursByDestinationID, destinationID)
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

	// Fetch highlights for all tours
	if len(tours) > 0 {
		tourIDs := make([]int, len(tours))
		for i, t := range tours {
			tourIDs[i] = t.ID
		}
		highlightsMap, err := r.GetHighlightsByTourIDs(ctx, tourIDs)
		if err != nil {
			return nil, err
		}
		for _, tour := range tours {
			tour.Highlights = highlightsMap[tour.ID]
		}
	}

	return tours, nil
}

func (r *TourRepositoryPostgres) Update(ctx context.Context, tour *domain.Tour) error {
	return r.db.QueryRowContext(ctx, queries.UpdateTour,
		tour.DestinationID,
		tour.ProviderID,
		tour.Name,
		tour.Description,
		tour.Price,
		tour.StartDate,
		tour.EndDate,
		tour.Capacity,
		tour.ID,
	).Scan(&tour.CreatedAt)
}

func (r *TourRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, queries.DeleteTour, id)
	return err
}

func (r *TourRepositoryPostgres) DecrementCapacity(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, queries.DecrementTourCapacity, id)
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

func (r *TourRepositoryPostgres) IncrementCapacity(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, queries.IncrementTourCapacity, id)
	return err
}

func (r *TourRepositoryPostgres) GetRemainingSpots(ctx context.Context, tourID int) (int, error) {
	var remaining int
	err := r.db.QueryRowContext(ctx, queries.GetRemainingSpotsByTourID, tourID).Scan(&remaining)
	return remaining, err
}

func (r *TourRepositoryPostgres) GetHighlightsByTourID(ctx context.Context, tourID int) ([]*domain.TourHighlight, error) {
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

func (r *TourRepositoryPostgres) GetHighlightsByTourIDs(ctx context.Context, tourIDs []int) (map[int][]*domain.TourHighlight, error) {
	if len(tourIDs) == 0 {
		return make(map[int][]*domain.TourHighlight), nil
	}

	result := make(map[int][]*domain.TourHighlight)
	for _, id := range tourIDs {
		result[id] = make([]*domain.TourHighlight, 0)
	}

	rows, err := r.db.QueryContext(ctx, queries.GetTourHighlightsByTourIDs, toSlice(tourIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var h domain.TourHighlight
		err := rows.Scan(&h.ID, &h.TourID, &h.ImageURL, &h.SortOrder, &h.CreatedAt)
		if err != nil {
			return nil, err
		}
		result[h.TourID] = append(result[h.TourID], &h)
	}

	// Sort each tour's highlights by sort_order
	for _, highlights := range result {
		sortHighlights(highlights)
	}

	return result, nil
}

func sortHighlights(highlights []*domain.TourHighlight) {
	for i := 0; i < len(highlights)-1; i++ {
		for j := i + 1; j < len(highlights); j++ {
			if highlights[i].SortOrder > highlights[j].SortOrder {
				highlights[i], highlights[j] = highlights[j], highlights[i]
			}
		}
	}
}

func toSlice(ids []int) []interface{} {
	result := make([]interface{}, len(ids))
	for i, id := range ids {
		result[i] = id
	}
	return result
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
	var providerID sql.NullInt64
	err := scanner.Scan(
		&tour.ID,
		&tour.DestinationID,
		&providerID,
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
	if providerID.Valid {
		tour.ProviderID = new(int)
		*tour.ProviderID = int(providerID.Int64)
	}
	tour.Description = description.String
	tour.DestinationDescription = destinationDescription.String
	tour.DestinationImageURL = destinationImageURL.String
	if tour.DestinationID == 0 {
		tour.DestinationID = joinedDestinationID
	}
	return nil
}
