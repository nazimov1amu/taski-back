package handlers

import (
	"encoding/json"
	"net/http"
	"taski_backend/internal/models"
	"taski_backend/internal/service"
)

type UsersHandler struct {
	service *service.UsersService
}

func NewUsersHandler(service *service.UsersService) *UsersHandler {
	return &UsersHandler{service: service}
}

func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.Get(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	createdUser, err := h.service.Create(r.Context(), user)
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}
	writeJSON(w, http.StatusOK, createdUser)
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user.ID = userID
	updatedUser, err := h.service.Update(r.Context(), user)
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}
	writeJSON(w, http.StatusOK, updatedUser)
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err := h.service.Delete(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}
	writeJSON(w, http.StatusOK, "User deleted successfully")
}

func (h *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var response models.LoginResponse
	token, refreshToken, err := h.service.Login(r.Context(), login)
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}
	response.Token = token
	response.RefreshToken = refreshToken
	
	writeJSON(w, http.StatusOK, response)
}

func (h *UsersHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshToken models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	token, err := h.service.RefreshToken(r.Context(), refreshToken.RefreshToken)
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}
	var response models.RefreshTokenResponse
	response.Token = token
	writeJSON(w, http.StatusOK, response)
}