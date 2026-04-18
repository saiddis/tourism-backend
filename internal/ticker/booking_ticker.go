package ticker

import (
	"context"
	"log"
	"time"

	"tourism-backend/internal/service"
)

type BookingTicker struct {
	bookingService *service.BookingService
	interval       time.Duration
	stopCh         chan struct{}
}

func NewBookingTicker(bookingService *service.BookingService, interval time.Duration) *BookingTicker {
	return &BookingTicker{
		bookingService: bookingService,
		interval:       interval,
		stopCh:         make(chan struct{}),
	}
}

func (t *BookingTicker) Start() {
	go func() {
		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()

		log.Println("Booking ticker started")

		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

				// Confirm pending bookings
				if err := t.bookingService.ConfirmAllDueBookings(ctx); err != nil {
					log.Printf("Booking ticker error (confirm): %v", err)
				} else {
					log.Println("Booking ticker: confirmed due bookings")
				}

				// Mark completed bookings
				if err := t.bookingService.MarkCompletedBookings(ctx); err != nil {
					log.Printf("Booking ticker error (complete): %v", err)
				} else {
					log.Println("Booking ticker: marked completed bookings")
				}

				cancel()
			case <-t.stopCh:
				log.Println("Booking ticker stopped")
				return
			}
		}
	}()
}

func (t *BookingTicker) Stop() {
	close(t.stopCh)
}
