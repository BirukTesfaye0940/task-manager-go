package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-manager-go/config"
	"task-manager-go/database"
	"task-manager-go/handlers"
	"task-manager-go/middleware"
	"task-manager-go/repositories"
	"task-manager-go/services"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create a bootstrap context with timeout for connecting to the database
	initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := database.Connect(initCtx, cfg)
	initCancel()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repositories.NewPostgresTaskRepo(db)
	service := services.NewTaskService(repo)
	handler := handlers.NewTaskHandler(service)

	userRepo := repositories.NewPostgresUserRepo(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	authMiddleware := middleware.Auth(cfg.JWTSecret)

	healthHandler := handlers.NewHealthHandler(db)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /health", healthHandler.Health)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	// Protected task routes — wrapped individually with the auth middleware
	mux.Handle("GET /tasks", authMiddleware(http.HandlerFunc(handler.GetTasks)))
	mux.Handle("POST /tasks", authMiddleware(http.HandlerFunc(handler.CreateTask)))
	mux.Handle("GET /tasks/{id}", authMiddleware(http.HandlerFunc(handler.GetTask)))
	mux.Handle("PUT /tasks/{id}", authMiddleware(http.HandlerFunc(handler.UpdateTask)))
	mux.Handle("DELETE /tasks/{id}", authMiddleware(http.HandlerFunc(handler.DeleteTask)))

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: middleware.RequestID(middleware.Logging(mux)),
	}

	// Create a channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run the server in a goroutine
	go func() {
		log.Printf("Starting server on :%s\n", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for an interrupt signal
	<-stop
	log.Println("\nShutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
