package postgres

import (
	"context"
	"database/sql"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository/queries"
)

type PaymentRepositoryPostgres struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepositoryPostgres {
	return &PaymentRepositoryPostgres{db: db}
}

func (r *PaymentRepositoryPostgres) Create(ctx context.Context, payment *domain.Payment) error {
	return r.db.QueryRowContext(ctx, queries.CreatePayment,
		payment.BookingID,
		payment.Amount,
		payment.Currency,
		payment.Status,
	).Scan(&payment.ID, &payment.CreatedAt)
}

func (r *PaymentRepositoryPostgres) GetByID(ctx context.Context, id int) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.QueryRowContext(ctx, queries.GetPaymentByID, id).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &payment, err
}

func (r *PaymentRepositoryPostgres) GetByBookingID(ctx context.Context, bookingID int) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.QueryRowContext(ctx, queries.GetPaymentByBookingID, bookingID).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &payment, err
}

func (r *PaymentRepositoryPostgres) UpdateStatus(ctx context.Context, id int, status domain.PaymentStatus) error {
	_, err := r.db.ExecContext(ctx, queries.UpdatePaymentStatus, status, id)
	return err
}
