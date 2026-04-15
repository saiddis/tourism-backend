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
	err := r.db.QueryRow(queries.GetBookingByID, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.TourID,
		&booking.Status,
		&booking.CreatedAt,
	)
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
		err = rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.TourID,
			&booking.Status,
			&booking.CreatedAt,
		)
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
		err = rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.TourID,
			&booking.Status,
			&booking.CreatedAt,
		)
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
