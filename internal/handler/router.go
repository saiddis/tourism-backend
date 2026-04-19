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
	providerHandler *ProviderHandler,
	providerAppHandler *ProviderApplicationHandler,
	tourHighlightHandler *TourHighlightHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORSMiddleware)

	// 🔓 Открытые endpoints (with /api prefix)
	r.Post("/api/auth/register", userHandler.Register)
	r.Post("/api/auth/login", userHandler.Login)
	r.Post("/api/auth/refresh", userHandler.Refresh)
	r.Delete("/api/auth/logout", userHandler.Logout)
	r.Get("/api/tours", tourHandler.GetAll)
	r.Get("/api/tours/{id}", tourHandler.GetByID)
	r.Get("/api/tours/destination/{id}", tourHandler.GetByDestinationID)
	r.Get("/api/reviews/tour/{id}", reviewHandler.GetByTourID)
	r.Get("/api/destinations", destinationHandler.GetAll)
	r.Get("/api/destinations/{id}", destinationHandler.GetByID)
	r.Get("/api/providers/active", providerHandler.GetActive)
	r.Get("/api/providers/user/{userId}", providerHandler.GetByUserID)
	r.Get("/api/providers/{id}", providerHandler.GetByID)
	r.Get("/api/tours/{id}/highlights", tourHighlightHandler.GetByTourID)
	r.Get("/api/provider-applications/{id}/accept", providerAppHandler.AcceptByToken)
	r.Get("/api/provider-applications/{id}/reject", providerAppHandler.RejectByToken)

	// 🔐 Защищённые endpoints (with /api prefix)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Client и выше
		r.Get("/api/users/me", userHandler.Me)
		r.Patch("/api/users/me", userHandler.UpdateMe)
		r.Post("/api/users/me/avatar", userHandler.UploadAvatar)
		r.Put("/api/users/me/avatar-url", userHandler.SetAvatarURL)
		r.Post("/api/users/deposit", userHandler.Deposit)
		r.Post("/api/bookings", bookingHandler.Create)
		r.Get("/api/bookings/user/{id}", bookingHandler.GetByUserID)
		r.Put("/api/bookings/{id}/status", bookingHandler.UpdateStatus)
		r.Post("/api/reviews", reviewHandler.Create)
		r.Post("/api/payments", paymentHandler.Create)
		r.Post("/api/provider-applications", providerAppHandler.Submit)
		r.Get("/api/provider-applications/me", providerAppHandler.GetByUserID)
		r.Get("/api/providers/me", providerHandler.GetByUserID)
		r.Put("/api/providers/{id}", providerHandler.Update)
		r.Post("/api/providers/me/tours", tourHandler.CreateForProvider)
		r.Put("/api/tours/{id}/renew", tourHandler.Renew)
		r.Post("/api/tours/{id}/highlights", tourHighlightHandler.Create)
		r.Delete("/api/tours/{id}/highlights/{highlightId}", tourHighlightHandler.Delete)

		// Только manager и admin
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("provider", "admin"))
			r.Get("/api/bookings", bookingHandler.GetAll)
			r.Delete("/api/bookings/{id}", bookingHandler.Delete)
			r.Post("/api/tours", tourHandler.Create)
			r.Put("/api/tours/{id}", tourHandler.Update)
			r.Delete("/api/tours/{id}", tourHandler.Delete)
			r.Post("/api/destinations", destinationHandler.Create)
			r.Delete("/api/destinations/{id}", destinationHandler.Delete)
			r.Put("/api/payments/{id}/status", paymentHandler.UpdateStatus)
			r.Get("/api/payments/{id}", paymentHandler.GetByID)
			r.Delete("/api/reviews/{id}", reviewHandler.Delete)
		})

		// Только admin
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin"))
			r.Get("/api/users", userHandler.GetAll)
			r.Get("/api/users/{id}", userHandler.GetByID)
			r.Post("/api/providers", providerHandler.Create)
			r.Get("/api/providers", providerHandler.GetAll)
			r.Put("/api/providers/{id}/active", providerHandler.ToggleActive)
			r.Get("/api/provider-applications", providerAppHandler.GetAll)
			r.Put("/api/provider-applications/{id}/accept", providerAppHandler.Accept)
			r.Put("/api/provider-applications/{id}/reject", providerAppHandler.Reject)
		})
	})

	return r
}
