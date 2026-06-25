package main

import (
	"log"
	"net/http"
	"task-manager-go/config"
	"task-manager-go/database"
	"task-manager-go/handlers"
	"task-manager-go/middleware"
	"task-manager-go/repositories"
	"task-manager-go/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repositories.NewPostgresTaskRepo(db)
	service := services.NewTaskService(repo)
	handler := handlers.NewTaskHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", handler.GetTasks)
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
	mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", middleware.Logging(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
