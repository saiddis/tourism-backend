package handler

import (
	middleware "tourism-backend/internal/middlewares"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	userHandler *UserHandler,
	tourHandler *TourHandler,
	bookingHandler *BookingHandler,
	paymentHandler *PaymentHandler,
	reviewHandler *ReviewHandler,
	destinationHandler *DestinationHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORSMiddleware)

	// 🔓 Открытые endpoints
	r.Post("/auth/register", userHandler.Register)
	r.Post("/auth/login", userHandler.Login)
	r.Post("/auth/refresh", userHandler.Refresh)
	r.Delete("/auth/logout", userHandler.Logout)
	r.Get("/tours", tourHandler.GetAll)
	r.Get("/tours/{id}", tourHandler.GetByID)
	r.Get("/tours/destination/{id}", tourHandler.GetByDestinationID)
	r.Get("/destinations", destinationHandler.GetAll)
	r.Get("/destinations/{id}", destinationHandler.GetByID)

	// 🔐 Защищённые endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Client и выше
		r.Get("/users/me", userHandler.Me)
		r.Post("/bookings", bookingHandler.Create)
		r.Get("/bookings/user/{id}", bookingHandler.GetByUserID)
		r.Post("/reviews", reviewHandler.Create)
		r.Post("/payments", paymentHandler.Create)

		// Только manager и admin
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("manager", "admin"))
			r.Get("/bookings", bookingHandler.GetAll)
			r.Put("/bookings/{id}/status", bookingHandler.UpdateStatus)
			r.Delete("/bookings/{id}", bookingHandler.Delete)
			r.Post("/tours", tourHandler.Create)
			r.Put("/tours/{id}", tourHandler.Update)
			r.Delete("/tours/{id}", tourHandler.Delete)
			r.Post("/destinations", destinationHandler.Create)
			r.Delete("/destinations/{id}", destinationHandler.Delete)
			r.Put("/payments/{id}/status", paymentHandler.UpdateStatus)
			r.Get("/payments/{id}", paymentHandler.GetByID)
			r.Get("/reviews/tour/{id}", reviewHandler.GetByTourID)
			r.Delete("/reviews/{id}", reviewHandler.Delete)
		})

		// Только admin
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin"))
			r.Get("/users", userHandler.GetAll)
			r.Get("/users/{id}", userHandler.GetByID)
		})
	})

	return r
}
