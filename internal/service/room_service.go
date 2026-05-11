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
}

type RoomService struct {
	RoomRepo RoomRepository
	UserRepo UserRepository
}

func NewRoomService(userRepo UserRepository, roomRepo RoomRepository) *RoomService {
	return &RoomService{
		RoomRepo: roomRepo,
		UserRepo: userRepo,
	}
}

func (r *RoomService) CreateRoom(ctx context.Context, userHostId uuid.UUID, countUserInRoom int) (*domain.Room, error) {
	existingUser, err := r.UserRepo.GetByID(ctx, userHostId)
	if existingUser == nil && err == nil {
		return nil, fmt.Errorf("user with id: %s not found", userHostId.String())
	}
	if err != nil {
		return nil, err
	}

	newRoom := *domain.NewRoom(userHostId, countUserInRoom)

	err = r.RoomRepo.Create(ctx, &newRoom)
	if err != nil {
		return nil, err
	}

	return &newRoom, nil
}
