package service

import (
	"context"
	"fmt"

	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
	"github.com/MikleMkhlv/QuizPlatform/internal/handler"
	"github.com/MikleMkhlv/QuizPlatform/internal/websocket"
	"github.com/google/uuid"
)

type QuestionRepository interface {
	Create(ctx context.Context, quizz *domain.Quizzes) error
	GetQuizzByID(ctx context.Context, quizzID uuid.UUID) (*domain.Quizzes, error)
	AddQuestionsWithAnswers(ctx context.Context, question *domain.Questions, options ...*domain.Options) error
	GetCountQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) (int, error)
}

type QusetionService struct {
	questionRepo QuestionRepository
	roomServs    websocket.RoomServiceInterface
	userServs    websocket.UserServiceInterface
}

func NewQusetionService(questionsRepo QuestionRepository, roomServs websocket.RoomServiceInterface,
	userServs websocket.UserServiceInterface) *QusetionService {
	return &QusetionService{
		questionRepo: questionsRepo,
		roomServs:    roomServs,
		userServs:    userServs,
	}
}

func (q *QusetionService) CreateQuizz(ctx context.Context, playerID uuid.UUID, title string) (*domain.Quizzes, error) {
	foundPlayer, err := q.userServs.GetPlayerByID(ctx, playerID)
	if err != nil {
		return nil, err
	}
	if foundPlayer == nil {
		return nil, fmt.Errorf("player {%v} is not found", playerID)
	}

	newQuizz := domain.NewQuiz(title, playerID)

	if err := q.questionRepo.Create(ctx, newQuizz); err != nil {
		return nil, err
	}
	return newQuizz, nil
}

func (q *QusetionService) GetQuizzByID(ctx context.Context, quizzID uuid.UUID) (*domain.Quizzes, error) {
	quizz, err := q.questionRepo.GetQuizzByID(ctx, quizzID)
	if err != nil {
		return nil, err
	}
	if quizz == nil {
		return nil, fmt.Errorf("quizz by id {%v} is not found", quizzID)
	}
	return quizz, nil
}

func (q *QusetionService) AddNewQuestionsWithAnsvers(ctx context.Context, quizID uuid.UUID, text string, optionsReq []handler.Options) error {
	foundQuiz, err := q.GetQuizzByID(ctx, quizID)
	if err != nil {
		return err
	}
	if foundQuiz == nil {
		return fmt.Errorf("quiz by quizID {%v} is not found", quizID)
	}
	if len(optionsReq) <= 1 {
		return fmt.Errorf("count options there nust be more than 2")
	}
	newQuestions := domain.NewQuestion(quizID, text, 0)
	currentCountQuestions, err := q.questionRepo.GetCountQuestionsByQuizID(ctx, foundQuiz.ID)
	if err != nil {
		return err
	}
	if currentCountQuestions != -1 {
		newQuestions = domain.NewQuestion(quizID, text, currentCountQuestions)
	}

	var options []*domain.Options
	for _, option := range optionsReq {
		op := domain.NewOption(newQuestions.ID, option.Text, option.IsCorrect)
		options = append(options, op)
	}

	if err := q.questionRepo.AddQuestionsWithAnswers(ctx, newQuestions, options...); err != nil {
		return err
	}

	return nil
}
