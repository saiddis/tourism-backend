package postgres

import (
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

func (r *BookingRepositoryPostgres) Create(booking *domain.Booking) error {
	return r.db.QueryRow(queries.CreateBooking,
		booking.UserID,
		booking.TourID,
		booking.Status,
	).Scan(&booking.ID, &booking.CreatedAt)
}

func (r *BookingRepositoryPostgres) GetByID(id int) (*domain.Booking, error) {
	var booking domain.Booking
	err := scanBooking(r.db.QueryRow(queries.GetBookingByID, id), &booking)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &booking, err
}

func (r *BookingRepositoryPostgres) GetByUserID(userID int) ([]*domain.Booking, error) {
	bookings := make([]*domain.Booking, 0)
	rows, err := r.db.Query(queries.GetBookingsByUserID, userID)
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

func (r *BookingRepositoryPostgres) GetAll() ([]*domain.Booking, error) {
	bookings := make([]*domain.Booking, 0)
	rows, err := r.db.Query(queries.GetAllBookings)
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

func (r *BookingRepositoryPostgres) UpdateStatus(id int, status domain.BookingStatus) error {
	_, err := r.db.Exec(queries.UpdateBookingStatus, status, id)
	return err
}

func (r *BookingRepositoryPostgres) Delete(id int) error {
	_, err := r.db.Exec(queries.DeleteBooking, id)
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
