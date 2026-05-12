package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
)

type RoomRepository interface {
	Create(ctx context.Context, room *domain.Room) error
	GetRoomById(ctx context.Context, roomId uuid.UUID) (*domain.Room, error)
	GetRoomByCode(ctx context.Context, roomCode string) (*domain.Room, error)
	AddPlayer(ctx context.Context, newPlauer *domain.RoomPlayer) error
	GetPlayersFromRoom(ctx context.Context, roomId uuid.UUID) ([]*domain.User, error)
	UpdateRoomStatus(ctx context.Context, roomId uuid.UUID, status domain.RoomStatus) error
}

type RoomService struct {
	roomRepo RoomRepository
	userRepo UserRepository
}

func NewRoomService(userRepo UserRepository, roomRepo RoomRepository) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
		userRepo: userRepo,
	}
}

func (r *RoomService) CreateRoom(ctx context.Context, userHostId uuid.UUID, countUserInRoom int) (*domain.Room, error) {
	existingUser, err := r.userRepo.GetByID(ctx, userHostId)
	if existingUser == nil && err == nil {
		return nil, fmt.Errorf("user with id: %s not found", userHostId.String())
	}
	if err != nil {
		return nil, err
	}

	newRoom := domain.NewRoom(userHostId, countUserInRoom)
	err = r.roomRepo.Create(ctx, newRoom)
	if err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (r *RoomService) Join(ctx context.Context, code string, userId uuid.UUID) (*domain.Room, error) {
	existingRoom, err := r.roomRepo.GetRoomByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get room by code: %w", err)
	}
	if existingRoom == nil {
		return nil, fmt.Errorf("room with code %s not found", code)
	}
	currentCountPlauerInRoom, err := r.roomRepo.GetPlayersFromRoom(ctx, existingRoom.ID)
	if err != nil {
		return nil, err
	}
	if existingRoom.Status != domain.RoomStatusWaiting {
		return nil, fmt.Errorf("room %s is not accepting players, status: %s", existingRoom.ID, existingRoom.Status)
	}
	if len(currentCountPlauerInRoom) >= existingRoom.MaxPlayer {
		return nil, fmt.Errorf("room %s is full (%d/%d)", existingRoom.ID, len(currentCountPlauerInRoom), existingRoom.MaxPlayer)
	}

	addedPlayer := domain.NewRoomPlayer(existingRoom.ID, userId)

	err = r.roomRepo.AddPlayer(ctx, addedPlayer)
	if err != nil {
		return nil, err
	}
	return existingRoom, nil
}

func (r *RoomService) UpdateRoomStatus(ctx context.Context, roomId uuid.UUID, status domain.RoomStatus) error {
	return r.roomRepo.UpdateRoomStatus(ctx, roomId, status)
}

func (r *RoomService) GetRoomByID(ctx context.Context, roomId uuid.UUID) (*domain.Room, error) {
	return r.roomRepo.GetRoomById(ctx, roomId)
}
