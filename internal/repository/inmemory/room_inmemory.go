package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
)

type InMemoryRoomRepository struct {
	mx         sync.RWMutex
	rooms      map[uuid.UUID]*domain.Room
	roomPlauer map[uuid.UUID][]domain.RoomPlauer
}

func NewInMemoryRoomRepository() *InMemoryRoomRepository {
	return &InMemoryRoomRepository{
		rooms:      make(map[uuid.UUID]*domain.Room),
		roomPlauer: make(map[uuid.UUID][]domain.RoomPlauer),
	}
}

func (rr *InMemoryRoomRepository) Create(ctx context.Context, room *domain.Room) error {
	rr.mx.Lock()
	defer rr.mx.Unlock()
	rr.rooms[room.ID] = room
	var plauers []domain.RoomPlauer
	plauers = append(plauers, *domain.NewRoomPlauer(room.ID, room.HostID))
	rr.roomPlauer[room.ID] = plauers
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

func (rr *InMemoryRoomRepository) AddPlauer(ctx context.Context, newPlauer *domain.RoomPlauer) error {
	rr.mx.Lock()
	defer rr.mx.Unlock()
	if _, ok := rr.roomPlauer[newPlauer.RoomId]; !ok {
		return nil
	}
	existingRoom := rr.roomPlauer[newPlauer.RoomId]
	existingRoom = append(existingRoom, *newPlauer)
	rr.roomPlauer[newPlauer.RoomId] = existingRoom
	return nil
}

func (rr *InMemoryRoomRepository) GetPlauersFromRoom(ctx context.Context, roomId uuid.UUID) ([]*domain.User, error) {
	rr.mx.RLock()
	defer rr.mx.RUnlock()
	return nil, nil
}

func (rr *InMemoryRoomRepository) UpdateRoomStatus(ctx context.Context, roomId uuid.UUID, status domain.RoomStatus) error {
	rr.mx.Lock()
	defer rr.mx.Unlock()
	if v, ok := rr.rooms[roomId]; ok {
		v.Status = status
		rr.rooms[roomId] = v
	}
	return nil
}
