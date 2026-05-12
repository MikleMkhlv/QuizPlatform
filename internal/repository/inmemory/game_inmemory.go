package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
)

type InMemoryGameRepository struct {
	mx       sync.RWMutex
	gameRepo map[uuid.UUID]*domain.GameState
}

func NewInMemoryGameRepository() *InMemoryGameRepository {
	return &InMemoryGameRepository{
		gameRepo: make(map[uuid.UUID]*domain.GameState),
	}
}

func (gr *InMemoryGameRepository) Save(ctx context.Context, gameState *domain.GameState) error {
	gr.mx.Lock()
	defer gr.mx.Unlock()
	gr.gameRepo[gameState.RoomID] = gameState
	return nil
}

func (gr *InMemoryGameRepository) Get(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error) {
	gr.mx.RLock()
	defer gr.mx.RUnlock()
	if v, ok := gr.gameRepo[roomID]; ok {
		return v, nil
	}
	return nil, nil
}

func (gr *InMemoryGameRepository) Delete(ctx context.Context, roomID uuid.UUID) error {
	gr.mx.Lock()
	defer gr.mx.Unlock()

	if _, ok := gr.gameRepo[roomID]; !ok {
		return fmt.Errorf("game state not found for room %s", roomID)
	}
	delete(gr.gameRepo, roomID)
	return nil
}
