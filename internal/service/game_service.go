package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
)

type GameRepository interface {
	Save(ctx context.Context, gameState *domain.GameState) error
	Get(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error)
	Delete(ctx context.Context, roomID uuid.UUID) error
}

type PlayService struct {
	playRRepo GameRepository
}

func NewGameRepository(playRRepo GameRepository) *PlayService {
	return &PlayService{
		playRRepo: playRRepo,
	}
}
