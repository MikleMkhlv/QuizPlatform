package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (p *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, name, username, rating, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := p.pool.Exec(ctx, query, user.ID, user.Name, user.Username, user.Rating, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("postgres create user: %w", err)
	}
	return nil
}

func (p *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, name, username, rating, created_at FROM users WHERE id = $1`
	user := domain.User{}
	err := p.pool.QueryRow(ctx, query, id.String()).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Rating,
		&user.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, name, username, rating, created_at FROM users WHERE username = $1`
	user := domain.User{}
	err := p.pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Rating,
		&user.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresUserRepository) UpdateRating(ctx context.Context, id uuid.UUID, delta int) error {
	existing, err := p.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return fmt.Errorf("such user {%s} does not exist", id.String())
	}
	newRating := existing.Rating + delta
	if newRating < 0 {
		newRating = 0
	}

	query := `UPDATE users SET rating = $1 WHERE id = $2`
	_, err = p.pool.Exec(ctx, query, newRating, id.String())
	if err != nil {
		return fmt.Errorf("postgres update rating: %w", err)
	}
	return nil
}

func (p *PostgresUserRepository) GetTopByRating(ctx context.Context, limit int) ([]*domain.User, error) {
	query := `SELECT id, name, username, rating, created_at FROM users ORDER BY rating DESC LIMIT $1`

	rows, err := p.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topUsers []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Rating, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		topUsers = append(topUsers, user)
	}

	return topUsers, nil
}
