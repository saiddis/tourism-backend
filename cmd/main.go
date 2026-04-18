package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"tourism-backend/config"
	"tourism-backend/internal/email"
	"tourism-backend/internal/handler"
	middleware "tourism-backend/internal/middlewares"
	"tourism-backend/internal/repository/postgres"
	"tourism-backend/internal/service"
	"tourism-backend/internal/ticker"
	"tourism-backend/storage"
)

func main() {
	// 1. Загружаем конфиг
	cfg := config.Load()

	middleware.InitSecrets(cfg.AccessTokenSecret, cfg.RefreshTokenSecret)

	// Set upload directory for avatars
	handler.SetUploadDir(cfg.UploadDir)

	// 2. Подключаемся к БД
	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	fmt.Println("Connected to database")

	// 3. Repository
	userRepo := postgres.NewUserRepository(db)
	tourRepo := postgres.NewTourRepositoryPostgres(db)
	bookingRepo := postgres.NewBookingRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	reviewRepo := postgres.NewReviewRepository(db)
	destinationRepo := postgres.NewDestinationRepository(db)
	providerRepo := postgres.NewProviderRepository(db)
	providerAppRepo := postgres.NewProviderApplicationRepository(db)
	tourHighlightRepo := postgres.NewTourHighlightRepositoryPostgres(db)

	// 4. Service
	userService := service.NewUserService(userRepo)
	tourService := service.NewTourService(tourRepo)
	emailService := email.NewEmailService()
	bookingService := service.NewBookingService(bookingRepo, userService, tourService, emailService)
	paymentService := service.NewPaymentService(paymentRepo)
	reviewService := service.NewReviewService(reviewRepo)
	destinationService := service.NewDestinationService(destinationRepo)
	providerService := service.NewProviderService(providerRepo, userRepo)
	providerAppService := service.NewProviderApplicationService(providerAppRepo, providerRepo, userRepo, emailService)
	tourHighlightService := service.NewTourHighlightService(tourHighlightRepo)

	// Start booking confirmation ticker
	bookingTicker := ticker.NewBookingTicker(bookingService, time.Hour)
	bookingTicker.Start()

	// 5. Handler
	userHandler := handler.NewUserHandler(userService, cfg.CookieDomain)
	tourHandler := handler.NewTourHandler(tourService, providerService)
	bookingHandler := handler.NewBookingHandler(bookingService, tourService, userService, paymentService)
	paymentHandler := handler.NewPaymentHandler(paymentService, bookingService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	destinationHandler := handler.NewDestinationHandler(destinationService)
	providerHandler := handler.NewProviderHandler(providerService)
	providerAppHandler := handler.NewProviderApplicationHandler(providerAppService)
	tourHighlightHandler := handler.NewTourHighlightHandler(tourHighlightService)

	// 6. Router
	router := handler.NewRouter(
		userHandler,
		tourHandler,
		bookingHandler,
		paymentHandler,
		reviewHandler,
		destinationHandler,
		providerHandler,
		providerAppHandler,
		tourHighlightHandler,
	)

	// Serve static files for uploads
	router.Mount("/uploads", http.StripPrefix("/uploads", http.FileServer(http.Dir(cfg.UploadDir))))

	// 7. Запуск сервера
	fmt.Println("Server starting on port", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
