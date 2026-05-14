package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
)

type userGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type InMemoryRoomRepository struct {
	mx         sync.RWMutex
	rooms      map[uuid.UUID]*domain.Room
	roomPlayer map[uuid.UUID][]domain.RoomPlayer
	usersRepo  userGetter
}

func NewInMemoryRoomRepository(usersRepo userGetter) *InMemoryRoomRepository {
	return &InMemoryRoomRepository{
		rooms:      make(map[uuid.UUID]*domain.Room),
		roomPlayer: make(map[uuid.UUID][]domain.RoomPlayer),
		usersRepo:  usersRepo,
	}
}

func (rr *InMemoryRoomRepository) Create(ctx context.Context, room *domain.Room) error {
	rr.mx.Lock()
	defer rr.mx.Unlock()
	rr.rooms[room.ID] = room
	var plauers []domain.RoomPlayer
	plauers = append(plauers, *domain.NewRoomPlayer(room.ID, room.HostID))
	rr.roomPlayer[room.ID] = plauers
	return nil
}

func (rr *InMemoryRoomRepository) GetRoomById(ctx context.Context, roomId uuid.UUID) (*domain.Room, error) {
	rr.mx.RLock()
	defer rr.mx.RUnlock()
	if _, ok := rr.rooms[roomId]; !ok {
		return nil, fmt.Errorf("not found room with id: %s", roomId.String())
	}
	existingRoom := rr.rooms[roomId]
	return existingRoom, nil
}

func (rr *InMemoryRoomRepository) GetRoomByCode(ctx context.Context, roomCode string) (*domain.Room, error) {
	rr.mx.RLock()
	defer rr.mx.RUnlock()
	for _, v := range rr.rooms {
		if v.Code == roomCode {
			return v, nil
		}
	}
	return nil, nil
}

func (rr *InMemoryRoomRepository) AddPlayer(ctx context.Context, newPlayer *domain.RoomPlayer) error {
	rr.mx.Lock()
	defer rr.mx.Unlock()
	if _, ok := rr.roomPlayer[newPlayer.RoomID]; !ok {
		return fmt.Errorf("room %s not found", newPlayer.RoomID)
	}
	existingRoom := rr.roomPlayer[newPlayer.RoomID]
	existingRoom = append(existingRoom, *newPlayer)
	rr.roomPlayer[newPlayer.RoomID] = existingRoom
	return nil
}

func (rr *InMemoryRoomRepository) GetPlayersFromRoom(ctx context.Context, roomId uuid.UUID) ([]*domain.User, error) {
	rr.mx.RLock()
	usersInRoom := rr.roomPlayer[roomId]
	rr.mx.RUnlock()
	var u []*domain.User
	for _, v := range usersInRoom {
		if v.RoomID == roomId {
			user, err := rr.usersRepo.GetByID(ctx, v.PlayerID)
			if err != nil {
				return nil, err
			}
			u = append(u, user)
		}
	}
	return u, nil
}

func (rr *InMemoryRoomRepository) UpdateRoomStatus(ctx context.Context, roomId uuid.UUID, status domain.RoomStatus) error {
	rr.mx.Lock()
	defer rr.mx.Unlock()

	room, ok := rr.rooms[roomId]
	if !ok {
		return fmt.Errorf("room %s not found", roomId)
	}
	room.Status = status
	return nil
}
