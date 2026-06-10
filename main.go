package main

import (
	"log"
	"net/http"
	"task-manager-go/handlers"
	"task-manager-go/repositories"
	"task-manager-go/services"
)

func main() {
	repo := repositories.NewInMemoryTaskRepo()
	service := services.NewTaskService(repo)
	handler := handlers.NewTaskHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", handler.GetTasks)
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
	mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
