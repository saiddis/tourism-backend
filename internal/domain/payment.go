package domain

import "time"

type PaymentStatus string

const (
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusPending PaymentStatus = "pending"
)

type Payment struct {
	ID        int           `json:"id"`
	BookingID int           `json:"booking_id"`
	Amount    float64       `json:"amount"`
	Currency  string        `json:"currency"`
	Status    PaymentStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}
