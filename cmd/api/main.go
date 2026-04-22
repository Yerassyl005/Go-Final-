package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"smartqueue/internal/database"
	"smartqueue/internal/handler"
	"smartqueue/internal/middleware"
	"smartqueue/internal/repository"
	"smartqueue/internal/service"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	router := mux.NewRouter()

	servicePointRepo := repository.NewServicePointPostgresRepository(db)
	queueRepo := repository.NewQueuePostgresRepository(db)
	ticketRepo := repository.NewTicketPostgresRepository(db)
	userRepo := repository.NewUserPostgresRepository(db)

	servicePointService := service.NewServicePointService(servicePointRepo)
	queueService := service.NewQueueService(queueRepo)
	ticketService := service.NewTicketService(ticketRepo, userRepo)
	authService := service.NewAuthService(userRepo, jwtSecret)
	qrService := service.NewQRService()

	servicePointHandler := handler.NewServicePointHandler(servicePointService)
	queueHandler := handler.NewQueueHandler(queueService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	authHandler := handler.NewAuthHandler(authService)
	qrHandler := handler.NewQRHandler(qrService)

	router.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/qr", qrHandler.GenerateQR).Methods("GET")
	router.HandleFunc("/tickets/{id}/qr", qrHandler.GetTicketQR).Methods("GET")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService))

	protected.HandleFunc("/auth/me", authHandler.Me).Methods("GET")

	protected.HandleFunc("/servicepoints", servicePointHandler.Create).Methods("POST")
	protected.HandleFunc("/servicepoints", servicePointHandler.GetAll).Methods("GET")

	protected.HandleFunc("/queues", queueHandler.Create).Methods("POST")
	protected.HandleFunc("/queues", queueHandler.GetAll).Methods("GET")
	protected.HandleFunc("/servicepoints/{id}/queues", queueHandler.GetByServicePoint).Methods("GET")
	protected.HandleFunc("/queues/{id}/display", queueHandler.GetDisplay).Methods("GET")
	protected.HandleFunc("/queues/{id}/stats", queueHandler.GetStats).Methods("GET")

	protected.HandleFunc("/tickets", ticketHandler.TakeTicket).Methods("POST")
	protected.HandleFunc("/tickets", ticketHandler.GetTickets).Methods("GET")
	protected.HandleFunc("/tickets/{id}/position", ticketHandler.GetPosition).Methods("GET")
	protected.HandleFunc("/tickets/{id}/call-skipped", ticketHandler.CallSkipped).Methods("POST")
	protected.HandleFunc("/queues/{id}/tickets/call-next", ticketHandler.CallNext).Methods("POST")
	protected.HandleFunc("/queues/{id}/tickets/recall-current", ticketHandler.RecallCurrent).Methods("POST")
	protected.HandleFunc("/queues/{id}/tickets/skip-current", ticketHandler.SkipCurrent).Methods("POST")
	protected.HandleFunc("/queues/{id}/tickets/complete-current", ticketHandler.CompleteCurrent).Methods("POST")

	log.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
