package ports

import (
	"context"

	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, name, username string) (*domain.User, error)
	UpdateRating(ctx context.Context, id uuid.UUID, delta int) error
	GetPlayerByID(ctx context.Context, playerID uuid.UUID) (*domain.User, error)
}

type RoomService interface {
	UpdateRoomStatus(ctx context.Context, roomId uuid.UUID, status domain.RoomStatus) error
	GetRoomByRoomCode(ctx context.Context, roomCode string) (*domain.Room, error)
}

type GameService interface {
	CreateGame(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error)
	StartGame(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error)
	SubmitAnswer(ctx context.Context, roomID, questionID, userID uuid.UUID, answer int) error
	FinishGame(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error)
}

type QuestionService interface {
	CreateQuizz(ctx context.Context, playerID uuid.UUID, title string) (*domain.Quizzes, error)
	GetQuizzByID(ctx context.Context, quizzID uuid.UUID) (*domain.Quizzes, error)
	AddNewQuestionsWithAnswers(ctx context.Context, quizID uuid.UUID, text string, options []OptionRequest) error
}

type OptionRequest struct {
	Text      string
	IsCorrect bool
}
