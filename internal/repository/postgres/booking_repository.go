package postgres

import (
	"context"
	"database/sql"
	"tourism-backend/internal/domain"
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

func scanBooking(scanner rowScanner, booking *domain.Booking) error {
	var tourDescription sql.NullString
	var destinationDescription sql.NullString
	var destinationImageURL sql.NullString
	var tourCreatedAt sql.NullTime
	var destinationCreatedAt sql.NullTime
	var joinedTourID int
	var joinedDestinationID int
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
	return nil
}
