package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
	"github.com/MikleMkhlv/QuizPlatform/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresQuestionsRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresQuestionsRepository(pool *pgxpool.Pool) *PostgresQuestionsRepository {
	return &PostgresQuestionsRepository{
		pool: pool,
	}
}

func (pq *PostgresQuestionsRepository) Create(ctx context.Context, quizz *domain.Quiz) error {
	query := `
			INSERT INTO quizzes (id, title, created_by)
			VALUES ($1, $2, $3)
`
	_, err := pq.pool.Exec(ctx, query, quizz.ID, quizz.Title, quizz.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (pq *PostgresQuestionsRepository) GetQuizzByID(ctx context.Context, quizzID uuid.UUID) (*domain.Quiz, error) {
	query := `
			SELECT id, title, created_by
			FROM quizzes
			WHERE id = $1
	`
	quizz := &domain.Quiz{}
	err := pq.pool.QueryRow(ctx, query, quizzID).Scan(&quizz.ID, &quizz.Title, &quizz.CreatedBy)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("quiz %s not found", quizzID)
	}
	if err != nil {
		return nil, err
	}
	return quizz, nil
}

func (pq *PostgresQuestionsRepository) AddQuestionsWithAnswers(ctx context.Context, question *domain.Question, options ...*domain.Option) error {
	tx, err := pq.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	queryQuestion := `
					INSERT INTO questions (id, quiz_id, text, order_num, correct_option_id)
					VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, queryQuestion, &question.ID, &question.QuizID, &question.Text, &question.OrderNum, &question.CorrectOptionID)
	if err != nil {
		return err
	}
	queryOption := `
				INSERT INTO options (id, question_id, text, is_correct)
				VALUES ($1, $2, $3, $4)
	`
	for _, option := range options {
		_, err = tx.Exec(ctx, queryOption, &option.ID, &option.QuestionID, &option.Text, &option.IsCorrect)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (pq *PostgresQuestionsRepository) GetCountQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) (int, error) {
	query := `
			SELECT COUNT(*) FROM questions WHERE quiz_id = $1
`
	var countQuestion int
	if err := pq.pool.QueryRow(ctx, query, quizID).Scan(&countQuestion); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, nil
		}
		return 0, err
	}
	return countQuestion, nil
}

func (pq *PostgresQuestionsRepository) GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*domain.Question, error) {
	query := `
		SELECT id, quiz_id, text, order_num, correct_option_id FROM questions WHERE id = $1
	`

	var question domain.Question
	if err := pq.pool.QueryRow(ctx, query, questionID).Scan(&question.ID, &question.QuizID, &question.Text, &question.OrderNum, &question.CorrectOptionID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("question {%v} not found", questionID)
		}
		return nil, err
	}
	return &question, nil
}

func (pq *PostgresQuestionsRepository) GetQuestionWithOptions(ctx context.Context, questionID uuid.UUID) (*dto.QuestionWithOptions, error) {
	query := `
			SELECT q.id, q.text,o.id, o.text, o.is_correct
			FROM questions q
			JOIN options o ON o.question_id = q.id
			WHERE q.id = $1
	`
	var result dto.QuestionWithOptions
	var options []domain.Option

	rows, err := pq.pool.Query(ctx, query, questionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("not found question: %v", err)
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var question domain.Question
		var option domain.Option

		if err := rows.Scan(
			&question.ID,
			&question.QuizID,
			&question.Text,
			&question.OrderNum,
			&question.CorrectOptionID,
			&option.ID,
			&option.Text,
			&option.IsCorrect,
		); err != nil {
			return nil, fmt.Errorf("error scan row: %v", err)
		}

		option.QuestionID = question.ID
		result.Question = question
		options = append(options, option)
	}

	result.Options = options
	return &result, nil
}
