package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourname/quiz-platform/internal/domain"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (p *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, name, username, reating, crteatedAt) VALUES ($1, $2, $3, $4, $5)`
	_, err := p.pool.Exec(ctx, query, user.Id, user.Name, user.Username, user.Reating, user.CrteatedAt)
	if err != nil {
		return fmt.Errorf("postgres create user: %w", err)
	}
	return nil
}

func (p *PostgresUserRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, name, username, reating, crteatedAt FROM users WHERE id = $1`
	user := domain.User{}
	err := p.pool.QueryRow(ctx, query, id.String()).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Reating,
		&user.CrteatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, name, username, reating, crteatedAt FROM users WHERE username = $1`
	user := domain.User{}
	err := p.pool.QueryRow(ctx, query, username).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Reating,
		&user.CrteatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresUserRepository) UpdateReating(ctx context.Context, id uuid.UUID, delta int) error {
	existing, err := p.GetById(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return fmt.Errorf("such user {%s} does not exist", id.String())
	}
	newReating := existing.Reating + delta
	if newReating < 0 {
		newReating = 0
	}

	query := `UPDATE users SET reating = $1 WHERE id = $2`
	_, err = p.pool.Exec(ctx, query, delta, id.String())
	if err != nil {
		return fmt.Errorf("postgres update rating: %w", err)
	}
	return nil
}
