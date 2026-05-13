package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
)

type GameRepository interface {
	Save(ctx context.Context, gameState *domain.GameState) error
	Get(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error)
	Delete(ctx context.Context, roomID uuid.UUID) error
}

type GameService struct {
	gameRepo GameRepository
	roomRepo RoomRepository
	userRepo UserRepository
}

func NewGameService(gameRepo GameRepository, roomRepo RoomRepository, userRepo UserRepository) *GameService {
	return &GameService{
		gameRepo: gameRepo,
		roomRepo: roomRepo,
		userRepo: userRepo,
	}
}

func (s *GameService) CreateGame(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error) {
	room, err := s.roomRepo.GetRoomById(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, fmt.Errorf("room %s not found", roomID)
	}
	playersInRoom, err := s.roomRepo.GetPlayersFromRoom(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("get players from room %s: %w", roomID, err)
	}
	state := domain.NewGameState(roomID, playersInRoom)
	if err := s.gameRepo.Save(ctx, state); err != nil {
		return nil, err
	}

	return state, nil
}

func (s *GameService) StartGame(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error) {
	existingState, err := s.gameRepo.Get(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if existingState == nil {
		return nil, fmt.Errorf("game state from room %s is not found in redis", roomID)
	}

	if existingState.Status != domain.GameStatusWaiting {
		return nil, fmt.Errorf("game in room %s already started or finished, status: %s",
			roomID, existingState.Status)
	}

	existingState.Start()

	if err := s.gameRepo.Save(ctx, existingState); err != nil {
		return nil, err
	}

	return existingState, nil
}
func (s *GameService) SubmitAnswer(ctx context.Context, roomID, userID uuid.UUID, answer int, isCorrect bool) error {
	existingState, err := s.gameRepo.Get(ctx, roomID)
	if err != nil {
		return err
	}
	if existingState == nil {
		return fmt.Errorf("game state from room %s is not found in redis", roomID)
	}
	if existingState.Status != domain.GameStatusActive {
		return fmt.Errorf("game in room %s is not active, status: %s", roomID, existingState.Status)
	}
	userAnswer := domain.NewPlayerAnswer(userID, answer, isCorrect)
	existingState.AddAnswer(*userAnswer)

	if err := s.gameRepo.Save(ctx, existingState); err != nil {
		return err
	}

	return nil
}

func (s *GameService) FinishGame(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error) {
	existingState, err := s.gameRepo.Get(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if existingState == nil {
		return nil, fmt.Errorf("game state from room %s is not found in redis", roomID)
	}
	if existingState.Status != domain.GameStatusActive {
		return nil, fmt.Errorf("game in room %s is not active, status: %s",
			roomID, existingState.Status)
	}
	existingState.Finish()

	if err := s.gameRepo.Delete(ctx, existingState.RoomID); err != nil {
		return nil, err
	}

	return existingState, nil
}
