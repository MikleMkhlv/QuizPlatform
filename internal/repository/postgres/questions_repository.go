package postgres

import (
	"context"
	"errors"

	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
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

func (pq *PostgresQuestionsRepository) Create(ctx context.Context, quizz *domain.Quizzes) error {
	query := `
			INSERT INTO quizzes (id, title, created_by)
			VALUES ($1, $2, $3)
`
	_, err := pq.pool.Exec(ctx, query, quizz.ID, quizz.Title, quizz.Created_by)
	if err != nil {
		return err
	}

	return nil
}

func (pq *PostgresQuestionsRepository) GetQuizzByID(ctx context.Context, quizzID uuid.UUID) (*domain.Quizzes, error) {
	query := `
			SELECT id, title, created_by
			FROM quizzes
			WHERE id = $1
	`
	quizz := &domain.Quizzes{}
	err := pq.pool.QueryRow(ctx, query, quizzID).Scan(&quizz.ID, &quizz.Title, &quizz.Created_by)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return quizz, nil
}

func (pq *PostgresQuestionsRepository) AddQuestionsWithAnswers(ctx context.Context, question *domain.Questions, options ...*domain.Options) error {
	tx, err := pq.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	queryQuestion := `
					INSERT INTO questions (id, quiz_id, text, order_num)
					VALUES ($1, $2, $3, $4)
	`
	_, err = pq.pool.Exec(ctx, queryQuestion, &question.ID, &question.Quiz_id, &question.Text, &question.Order_num)
	if err != nil {
		return err
	}
	queryOption := `
				INSERT INTO options (id, questions_id, text, is_correct)
				VALUES ($1, $2, $3, $4)
	`
	for _, option := range options {
		_, err = pq.pool.Exec(ctx, queryOption, &option.ID, &option.Question_id, &option.Text, &option.IsCorrect)
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
			SELECT COUNT(quiz_id) FROM questions WHERE quiz_id = $1
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

func (pq *PostgresQuestionsRepository) GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*domain.Questions, error) {
	query := `
		SELECT (id, quiz_id, text, order_num) FROM questions WHERE id = $1
	`

	var question *domain.Questions
	if err := pq.pool.QueryRow(ctx, query, questionID).Scan(&question); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return question, nil
}
