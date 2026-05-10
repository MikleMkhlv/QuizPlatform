package postgres

import (
	"context"

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
	return nil
}

func (p *PostgresUserRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (p *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return nil, nil
}

func (p *PostgresUserRepository) UpdateReating(ctx context.Context, id uuid.UUID, delta int) error {
	return nil
}
