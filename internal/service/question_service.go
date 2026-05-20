package service

import (
	"context"
	"fmt"

	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
	"github.com/MikleMkhlv/QuizPlatform/internal/handler"
	"github.com/google/uuid"
)

type QuestionRepository interface {
	Create(ctx context.Context, quizz *domain.Quizzes) error
	GetQuizzByID(ctx context.Context, quizzID uuid.UUID) (*domain.Quizzes, error)
	AddQuestionsWithAnswers(ctx context.Context, question *domain.Questions, options ...*domain.Options) error
	GetCountQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) (int, error)
	GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*domain.Questions, error)
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

	var options []*domain.Options
	for index, option := range optionsReq {
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

	newQuestion := domain.NewQuestion(quizID, text, orderNum, correctIndex)

	// Обновляем questionID в опциях
	for _, op := range options {
		op.Question_id = newQuestion.ID
	}

	return q.questionRepo.AddQuestionsWithAnswers(ctx, newQuestion, options...)
}

func checkSingleCorrectAnswer(options []*domain.Options) error {
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

func getCorrectOptionIndex(options []*domain.Options) int {
	for index, op := range options {
		if op.IsCorrect {
			return index
		}
	}
	return -1
}

func (q *QusetionService) GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*domain.Questions, error) {
	return q.questionRepo.GetQuestionByID(ctx, questionID)
}
