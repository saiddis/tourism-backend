package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository"
	"tourism-backend/internal/repository/queries"
)

type BookingRepositoryPostgres struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepositoryPostgres {
	return &BookingRepositoryPostgres{db: db}
}

func (r *BookingRepositoryPostgres) Create(ctx context.Context, booking *domain.Booking) error {
	return r.db.QueryRowContext(ctx, queries.CreateBooking,
		booking.UserID,
		booking.TourID,
		booking.Status,
	).Scan(&booking.ID, &booking.CreatedAt)
}

func (r *BookingRepositoryPostgres) GetByID(ctx context.Context, id int) (*domain.Booking, error) {
	var booking domain.Booking
	err := scanBooking(r.db.QueryRowContext(ctx, queries.GetBookingByID, id), &booking)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &booking, err
}

func (r *BookingRepositoryPostgres) GetByUserID(ctx context.Context, userID int) ([]*domain.Booking, error) {
	bookings := make([]*domain.Booking, 0)
	rows, err := r.db.QueryContext(ctx, queries.GetBookingsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var booking domain.Booking
		err = scanBooking(rows, &booking)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

func (r *BookingRepositoryPostgres) GetAll(ctx context.Context) ([]*domain.Booking, error) {
	bookings := make([]*domain.Booking, 0)
	rows, err := r.db.QueryContext(ctx, queries.GetAllBookings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var booking domain.Booking
		err = scanBooking(rows, &booking)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

func (r *BookingRepositoryPostgres) UpdateStatus(ctx context.Context, id int, status domain.BookingStatus) error {
	_, err := r.db.ExecContext(ctx, queries.UpdateBookingStatus, status, id)
	return err
}

func (r *BookingRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, queries.DeleteBooking, id)
	return err
}

func (r *BookingRepositoryPostgres) GetPendingBookingsCountByTourID(ctx context.Context, tourID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, queries.GetPendingBookingsCountByTourID, tourID).Scan(&count)
	return count, err
}

func (r *BookingRepositoryPostgres) GetActiveBookingsCount(ctx context.Context, tourID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, queries.GetActiveBookingsCountByTourID, tourID).Scan(&count)
	return count, err
}

func (r *BookingRepositoryPostgres) GetPendingBookingsByTourID(ctx context.Context, tourID int) ([]repository.PendingBookingWithUser, error) {
	rows, err := r.db.QueryContext(ctx, queries.GetPendingBookingsByTourID, tourID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]repository.PendingBookingWithUser, 0)
	for rows.Next() {
		var b repository.PendingBookingWithUser
		var status string
		err := rows.Scan(&b.BookingID, &b.UserID, &b.TourID, &status, &b.CreatedAt, &b.Email, &b.Name)
		if err != nil {
			return nil, err
		}
		b.Status = domain.BookingStatus(status)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *BookingRepositoryPostgres) ConfirmPendingBookingsForTour(ctx context.Context, tourID int) error {
	_, err := r.db.ExecContext(ctx, queries.ConfirmPendingBookingsForTour, tourID)
	return err
}

type TourConfirmationInfo struct {
	TourID       int
	Name         string
	StartDate    time.Time
	Capacity     int
	PendingCount int
}

func (r *BookingRepositoryPostgres) GetToursNeedingConfirmation(ctx context.Context) ([]repository.TourConfirmationInfo, error) {
	rows, err := r.db.QueryContext(ctx, queries.GetToursNeedingConfirmation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tours := make([]repository.TourConfirmationInfo, 0)
	for rows.Next() {
		var t repository.TourConfirmationInfo
		err := rows.Scan(&t.TourID, &t.Name, &t.StartDate, &t.Capacity, &t.PendingCount)
		if err != nil {
			return nil, err
		}
		tours = append(tours, t)
	}
	return tours, nil
}

func (r *BookingRepositoryPostgres) MarkBookingsCompleted(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, queries.MarkBookingsCompleted)
	return err
}

func (r *BookingRepositoryPostgres) GetCompletedBookingsByUserID(ctx context.Context, userID int) ([]*domain.Booking, error) {
	rows, err := r.db.QueryContext(ctx, queries.GetCompletedBookingsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]*domain.Booking, 0)
	for rows.Next() {
		var booking domain.Booking
		err := scanBooking(rows, &booking)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

func scanBooking(scanner rowScanner, booking *domain.Booking) error {
	var tourDescription sql.NullString
	var destinationDescription sql.NullString
	var destinationImageURL sql.NullString
	var tourCreatedAt sql.NullTime
	var destinationCreatedAt sql.NullTime
	var joinedTourID int
	var joinedDestinationID int
	var highlightsJSON []byte
	err := scanner.Scan(
		&booking.ID,
		&booking.UserID,
		&booking.TourID,
		&booking.Status,
		&booking.CreatedAt,
		&joinedTourID,
		&joinedDestinationID,
		&booking.TourName,
		&tourDescription,
		&booking.TourPrice,
		&booking.TourStartDate,
		&booking.TourEndDate,
		&booking.TourCapacity,
		&tourCreatedAt,
		&joinedDestinationID,
		&booking.DestinationName,
		&destinationDescription,
		&destinationImageURL,
		&destinationCreatedAt,
		&booking.RemainingSpots,
		&highlightsJSON,
	)
	if err != nil {
		return err
	}
	if booking.TourID == 0 {
		booking.TourID = joinedTourID
	}
	booking.DestinationID = joinedDestinationID
	booking.TourDescription = tourDescription.String
	booking.DestinationDescription = destinationDescription.String
	booking.DestinationImageURL = destinationImageURL.String

	if highlightsJSON != nil {
		var highlights []domain.TourHighlight
		if err := json.Unmarshal(highlightsJSON, &highlights); err == nil {
			booking.TourHighlights = make([]*domain.TourHighlight, len(highlights))
			for i := range highlights {
				booking.TourHighlights[i] = &highlights[i]
			}
		}
	}

	return nil
}
