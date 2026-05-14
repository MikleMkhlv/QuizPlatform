package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
)

type RoomServiceInterface interface {
	CreateRoom(ctx context.Context, userHostId uuid.UUID, countUserInRoom int) (*domain.Room, error)
	Join(ctx context.Context, code string, userId uuid.UUID) (*domain.Room, error)
}

type RoomHandler struct {
	roomSrvs RoomServiceInterface
}

func NewRoomHandler(roomSrvs RoomServiceInterface) *RoomHandler {
	return &RoomHandler{
		roomSrvs: roomSrvs,
	}
}

type RoomCreateRequest struct {
	MaxPlayers int `json:"max_players"`
}

type RoomCreateResponse struct {
	Code       string            `json:"room_code"`
	Status     domain.RoomStatus `json:"status"`
	MaxPlayers int               `json:"max_players"`
	CreatedAt  time.Time         `json:"created_at"`
}

func (rh *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	userID, err := ParseUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var roomRequset RoomCreateRequest
	err = json.NewDecoder(r.Body).Decode(&roomRequset)
	if err != nil {
		http.Error(w, "error decode request body", http.StatusBadRequest)
		return
	}
	if roomRequset.MaxPlayers <= 0 {
		http.Error(w, "'playerCount' must be greater than 0", http.StatusBadRequest)
		return
	}

	createdRoom, err := rh.roomSrvs.CreateRoom(ctx, userID, roomRequset.MaxPlayers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := RoomCreateResponse{
		Code:       createdRoom.Code,
		Status:     createdRoom.Status,
		MaxPlayers: createdRoom.MaxPlayer,
		CreatedAt:  createdRoom.CreatedAt,
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

type RoomJoinRequest struct {
	Code string `json:"room_code"`
}

type RoomJoinResponse struct {
	Code      string            `json:"room_code"`
	Status    domain.RoomStatus `json:"status"`
	MaxPlayer int               `json:"max_player"`
}

func (rh *RoomHandler) JoinInRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	userID, err := ParseUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req RoomJoinRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "error decode request body", http.StatusBadRequest)
		return
	}
	if req.Code == "" {
		http.Error(w, "error request body. Code required", http.StatusBadRequest)
		return
	}
	room, err := rh.roomSrvs.Join(ctx, req.Code, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := RoomJoinResponse{
		Code:      room.Code,
		Status:    room.Status,
		MaxPlayer: room.MaxPlayer,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("error marshal UserResponse: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
