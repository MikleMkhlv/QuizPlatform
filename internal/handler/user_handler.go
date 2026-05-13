package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
)

type UserServiceInterface interface {
	Register(ctx context.Context, name, username string) (*domain.User, error)
}
type UserHandler struct {
	userService UserServiceInterface
}

func NewUserHandler(userService UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type RegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	var userReq RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "error decode request body", http.StatusBadRequest)
		return
	}
	if userReq.Name == "" {
		http.Error(w, "error request body. Name required", http.StatusBadRequest)
		return
	}
	if userReq.Username == "" {
		http.Error(w, "error request body. Username required", http.StatusBadRequest)
		return
	}

	registeredUser, err := h.userService.Register(ctx, userReq.Name, userReq.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := RegisterResponse{
		ID:        registeredUser.ID,
		Name:      registeredUser.Name,
		Username:  registeredUser.Username,
		Rating:    registeredUser.Rating,
		CreatedAt: registeredUser.CreatedAt,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("error marshal UserResponse: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
