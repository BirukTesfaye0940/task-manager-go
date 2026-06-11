package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"task-manager-go/internal/response"
	"task-manager-go/models"
	"task-manager-go/repositories"
	"task-manager-go/services"
)

type TaskHandler struct {
	service *services.TaskService
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	task := h.service.GetAllTasks()

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Done:        false,
	}

	createdTask := h.service.CreateTask(task)

	response.JSON(
		w,
		http.StatusCreated,
		createdTask,
	)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(
			w,
			"invalid task ID",
			http.StatusBadRequest,
		)
		return
	}

	task, err := h.service.GetTask(id)

	if err != nil {

		if errors.Is(
			err,
			repositories.ErrTaskNotFound,
		) {
			http.Error(
				w,
				"task not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(
	w http.ResponseWriter,
	r *http.Request,
) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(
			w,
			"invalid task id",
			http.StatusBadRequest,
		)
		return
	}

	var req UpdateTaskRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
	}

	updatedTask, err := h.service.UpdateTask(
		id,
		task,
	)

	if err != nil {

		http.Error(
			w,
			"task not found",
			http.StatusNotFound,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(updatedTask)
}

func (h *TaskHandler) DeleteTask(
	w http.ResponseWriter,
	r *http.Request,
) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(
			w,
			"invalid task id",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.DeleteTask(id)

	if err != nil {

		http.Error(
			w,
			"task not found",
			http.StatusNotFound,
		)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
