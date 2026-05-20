package service

import (
	"context"
	"fmt"

	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
	"github.com/MikleMkhlv/QuizPlatform/internal/ports"
	"github.com/google/uuid"
)

type QuestionRepository interface {
	Create(ctx context.Context, quizz *domain.Quiz) error
	GetQuizzByID(ctx context.Context, quizzID uuid.UUID) (*domain.Quiz, error)
	AddQuestionsWithAnswers(ctx context.Context, question *domain.Question, options ...*domain.Option) error
	GetCountQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) (int, error)
	GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*domain.Question, error)
}

type UserServiceInterface interface {
	GetPlayerByID(ctx context.Context, playerID uuid.UUID) (*domain.User, error)
}

type QusetionService struct {
	questionRepo QuestionRepository
	userServs    UserServiceInterface
}

func NewQusetionService(questionsRepo QuestionRepository, userServs UserServiceInterface,
) *QusetionService {
	return &QusetionService{
		questionRepo: questionsRepo,
		userServs:    userServs,
	}
}

func (q *QusetionService) CreateQuizz(ctx context.Context, playerID uuid.UUID, title string) (*domain.Quiz, error) {
	foundPlayer, err := q.userServs.GetPlayerByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("player %v not found: %w", playerID, err)
	}

	newQuizz := domain.NewQuiz(title, foundPlayer.ID)

	if err := q.questionRepo.Create(ctx, newQuizz); err != nil {
		return nil, err
	}
	return newQuizz, nil
}

func (q *QusetionService) GetQuizzByID(ctx context.Context, quizzID uuid.UUID) (*domain.Quiz, error) {
	quizz, err := q.questionRepo.GetQuizzByID(ctx, quizzID)
	if err != nil {
		return nil, err
	}
	if quizz == nil {
		return nil, fmt.Errorf("quizz by id {%v} is not found", quizzID)
	}
	return quizz, nil
}

func (q *QusetionService) AddNewQuestionsWithAnswers(ctx context.Context, quizID uuid.UUID, text string, opt []ports.OptionRequest) error {
	const minOptions = 2

	if _, err := q.GetQuizzByID(ctx, quizID); err != nil {
		return err
	}
	if len(opt) < minOptions {
		return fmt.Errorf("question must have at least %d options, got %d", minOptions, len(opt))
	}

	var options []*domain.Option
	for index, option := range opt {
		op := domain.NewOption(index, uuid.Nil, option.Text, option.IsCorrect)
		options = append(options, op)
	}

	if err := checkSingleCorrectAnswer(options); err != nil {
		return err
	}

	correctIndex := getCorrectOptionIndex(options)
	if correctIndex == -1 {
		return fmt.Errorf("no correct answer provided")
	}

	currentCount, err := q.questionRepo.GetCountQuestionsByQuizID(ctx, quizID)
	if err != nil {
		return err
	}

	orderNum := 0
	if currentCount > 0 {
		orderNum = currentCount
	}

	correctOptionID := options[correctIndex].ID
	newQuestion := domain.NewQuestion(quizID, text, orderNum, correctOptionID)

	for _, op := range options {
		op.QuestionID = newQuestion.ID
	}

	return q.questionRepo.AddQuestionsWithAnswers(ctx, newQuestion, options...)
}

func checkSingleCorrectAnswer(options []*domain.Option) error {
	correctCount := 0
	for _, op := range options {
		if op.IsCorrect {
			correctCount++
		}
	}
	if correctCount == 0 {
		return fmt.Errorf("no correct answer in options")
	}
	if correctCount > 1 {
		return fmt.Errorf("multiple correct answers in options, expected exactly one")
	}
	return nil
}

func getCorrectOptionIndex(options []*domain.Option) int {
	for index, op := range options {
		if op.IsCorrect {
			return index
		}
	}
	return -1
}

func (q *QusetionService) GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*domain.Question, error) {
	return q.questionRepo.GetQuestionByID(ctx, questionID)
}
