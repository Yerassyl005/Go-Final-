package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"myapp-auth/internal/database"
	"myapp-auth/internal/handler"
	"myapp-auth/internal/middleware"
	"myapp-auth/internal/repository"
	"myapp-auth/internal/service"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	db, err := database.NewPostgresPool(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/me", middleware.AuthMiddleware(authService, authHandler.Me))

	log.Println("server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
