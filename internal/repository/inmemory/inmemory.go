package inmemory

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
)

type InMemoryUserRepository struct {
	repo map[int]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{repo: make(map[int]*domain.User)}
}

func (p *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (p *InMemoryUserRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (p *InMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return nil, nil
}

func (p *InMemoryUserRepository) UpdateReating(ctx context.Context, id uuid.UUID, delta int) error {
	return nil
}
