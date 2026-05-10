package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateReating(ctx context.Context, id uuid.UUID, delta int) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (ur *UserService) Register(ctx context.Context, name, username string) (*domain.User, error) {
	return nil, nil
}
