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
	roomPlauer map[uuid.UUID]*domain.RoomPlauer
}

func NewInMemoryRoomRepository() *InMemoryRoomRepository {
	return &InMemoryRoomRepository{
		rooms:      make(map[uuid.UUID]*domain.Room),
		roomPlauer: make(map[uuid.UUID]*domain.RoomPlauer),
	}
}

func (rr *InMemoryRoomRepository) Create(ctx context.Context, room *domain.Room) error {
	rr.mx.Lock()
	defer rr.mx.Unlock()
	rr.rooms[room.ID] = room
	rr.roomPlauer[room.ID] = domain.NewRoomPlauer(room.ID, room.HostID)
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
