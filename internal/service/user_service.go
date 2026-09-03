package service

import (
	"context"
	"fmt"

	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateRating(ctx context.Context, id uuid.UUID, delta int) error
	GetTopByRating(ctx context.Context, limit int) ([]*domain.User, error)
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
	existing, err := ur.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	//Проверка: Есть ли полученный username в БД. Если "Да", то регистрация прекращается. Если "Нет", то идем дальше
	if existing != nil {
		return nil, fmt.Errorf("user %s is already taken", existing.Username)
	}
	// Создание нового юзера
	user := domain.NewUser(name, username)

	//Сохранение в репозиторий
	if err := ur.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user %s is error: %w", user.Username, err)
	}
	return user, nil
}

func (ur *UserService) UpdateRating(ctx context.Context, id uuid.UUID, delta int) error {
	return ur.repo.UpdateRating(ctx, id, delta)
}

func (ur *UserService) GetTopUsers(ctx context.Context, limit int) ([]*domain.User, error) {
	topUsers, err := ur.repo.GetTopByRating(ctx, limit)
	if err != nil {
		return nil, err
	}

	return topUsers, nil
}
func (ur *UserService) GetPlayerByID(ctx context.Context, playerID uuid.UUID) (*domain.User, error) {
	return ur.GetPlayerByID(ctx, playerID)
}
