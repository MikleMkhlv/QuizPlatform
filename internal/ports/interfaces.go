package ports

import (
	"context"

	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
	"github.com/MikleMkhlv/QuizPlatform/internal/dto"
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

type GameServiceI interface {
	CreateGame(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error)
	StartGame(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error)
	SubmitAnswer(ctx context.Context, roomID, questionID, userID uuid.UUID, optionID uuid.UUID) error
	FinishGame(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error)
}

type QuestionService interface {
	CreateQuizz(ctx context.Context, playerID uuid.UUID, title string) (*domain.Quiz, error)
	GetQuizzByID(ctx context.Context, quizzID uuid.UUID) (*domain.Quiz, error)
	AddNewQuestionsWithAnswers(ctx context.Context, quizID uuid.UUID, text string, options []OptionRequest) error
	GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*domain.Question, error)
	GetQuestionWithOptions(ctx context.Context, questionID uuid.UUID) (*dto.QuestionWithOptions, error)
}

type OptionRequest struct {
	Text      string
	IsCorrect bool
}
