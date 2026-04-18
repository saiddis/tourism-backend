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

	// 🔓 Открытые endpoints
	r.Post("/auth/register", userHandler.Register)
	r.Post("/auth/login", userHandler.Login)
	r.Post("/auth/refresh", userHandler.Refresh)
	r.Delete("/auth/logout", userHandler.Logout)
	r.Get("/tours", tourHandler.GetAll)
	r.Get("/tours/{id}", tourHandler.GetByID)
	r.Get("/tours/destination/{id}", tourHandler.GetByDestinationID)
	r.Get("/reviews/tour/{id}", reviewHandler.GetByTourID)
	r.Get("/destinations", destinationHandler.GetAll)
	r.Get("/destinations/{id}", destinationHandler.GetByID)
	r.Get("/providers/active", providerHandler.GetActive)
	r.Get("/providers/user/{userId}", providerHandler.GetByUserID)
	r.Get("/providers/{id}", providerHandler.GetByID)
	r.Get("/tours/{id}/highlights", tourHighlightHandler.GetByTourID)
	r.Get("/provider-applications/{id}/accept", providerAppHandler.AcceptByToken)
	r.Get("/provider-applications/{id}/reject", providerAppHandler.RejectByToken)

	// 🔐 Защищённые endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Client и выше
		r.Get("/users/me", userHandler.Me)
		r.Patch("/users/me", userHandler.UpdateMe)
		r.Post("/users/me/avatar", userHandler.UploadAvatar)
		r.Put("/users/me/avatar-url", userHandler.SetAvatarURL)
		r.Post("/users/deposit", userHandler.Deposit)
		r.Post("/bookings", bookingHandler.Create)
		r.Get("/bookings/user/{id}", bookingHandler.GetByUserID)
		r.Put("/bookings/{id}/status", bookingHandler.UpdateStatus)
		r.Post("/reviews", reviewHandler.Create)
		r.Post("/payments", paymentHandler.Create)
		r.Post("/provider-applications", providerAppHandler.Submit)
		r.Get("/provider-applications/me", providerAppHandler.GetByUserID)
		r.Get("/providers/me", providerHandler.GetByUserID)
		r.Put("/providers/{id}", providerHandler.Update)
		r.Post("/providers/me/tours", tourHandler.CreateForProvider)
		r.Put("/tours/{id}/renew", tourHandler.Renew)
		r.Post("/tours/{id}/highlights", tourHighlightHandler.Create)
		r.Delete("/tours/{id}/highlights/{highlightId}", tourHighlightHandler.Delete)

		// Только manager и admin
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("provider", "admin"))
			r.Get("/bookings", bookingHandler.GetAll)
			r.Delete("/bookings/{id}", bookingHandler.Delete)
			r.Post("/tours", tourHandler.Create)
			r.Put("/tours/{id}", tourHandler.Update)
			r.Delete("/tours/{id}", tourHandler.Delete)
			r.Post("/destinations", destinationHandler.Create)
			r.Delete("/destinations/{id}", destinationHandler.Delete)
			r.Put("/payments/{id}/status", paymentHandler.UpdateStatus)
			r.Get("/payments/{id}", paymentHandler.GetByID)
			r.Delete("/reviews/{id}", reviewHandler.Delete)
		})

		// Только admin
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin"))
			r.Get("/users", userHandler.GetAll)
			r.Get("/users/{id}", userHandler.GetByID)
			r.Post("/providers", providerHandler.Create)
			r.Get("/providers", providerHandler.GetAll)
			r.Put("/providers/{id}/active", providerHandler.ToggleActive)
			r.Get("/provider-applications", providerAppHandler.GetAll)
			r.Put("/provider-applications/{id}/accept", providerAppHandler.Accept)
			r.Put("/provider-applications/{id}/reject", providerAppHandler.Reject)
		})
	})

	return r
}
