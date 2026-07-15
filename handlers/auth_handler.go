package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"task-manager-go/internal/response"
	"task-manager-go/repositories"
	"task-manager-go/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorResponse{Error: "invalid request body"})
		return
	}

	user, err := h.authService.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, repositories.ErrUsernameTaken) ||
			errors.Is(err, services.ErrUsernameRequired) ||
			errors.Is(err, services.ErrPasswordRequired) {
			response.JSON(w, http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
			return
		}
		response.JSON(w, http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorResponse{Error: "invalid request body"})
		return
	}

	token, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			response.JSON(w, http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
			return
		}
		response.JSON(w, http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		return
	}

	response.JSON(w, http.StatusOK, LoginResponse{Token: token})
}
