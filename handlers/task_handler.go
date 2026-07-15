package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"task-manager-go/internal/ctxkey"
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

func (h *TaskHandler) getUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(ctxkey.UserIDKey).(int)
	return userID, ok
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.getUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	task, err := h.service.GetAllTasks(ctx, userID)

	if err != nil {
		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Done:        false,
	}

	ctx := r.Context()
	createdTask, err := h.service.CreateTask(ctx, task)

	if err != nil {
		response.JSON(
			w,
			http.StatusBadRequest,
			response.ErrorResponse{
				Error: err.Error(),
			},
		)
		return
	}

	response.JSON(
		w,
		http.StatusCreated,
		createdTask,
	)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.getUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	task, err := h.service.GetTask(ctx, id, userID)

	if err != nil {

		if errors.Is(
			err,
			repositories.ErrTaskNotFound,
		) {
			response.JSON(
				w,
				http.StatusNotFound,
				response.ErrorResponse{
					Error: "task not found",
				},
			)
			return
		}

		response.JSON(
			w,
			http.StatusInternalServerError,
			response.ErrorResponse{
				Error: "internal server error",
			},
		)
		return
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
	userID, ok := h.getUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
	}

	ctx := r.Context()
	updatedTask, err := h.service.UpdateTask(
		ctx,
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
	userID, ok := h.getUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	ctx := r.Context()
	err = h.service.DeleteTask(ctx, id, userID)

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
