package handlers

import (
	"database/sql"
	"net/http"
	"task-manager-go/internal/response"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthStatus struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:   "healthy",
		Services: make(map[string]string),
	}

	statusCode := http.StatusOK

	// Check Database health
	if h.db != nil {
		if err := h.db.PingContext(r.Context()); err != nil {
			status.Status = "unhealthy"
			status.Services["database"] = "down"
			statusCode = http.StatusServiceUnavailable
		} else {
			status.Services["database"] = "up"
		}
	} else {
		status.Status = "unhealthy"
		status.Services["database"] = "not initialized"
		statusCode = http.StatusInternalServerError
	}

	response.JSON(w, statusCode, status)
}
