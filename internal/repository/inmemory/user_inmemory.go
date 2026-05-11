package inmemory

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
)

type InMemoryUserRepository struct {
	repo  map[uuid.UUID]*domain.User
	mutex sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		repo: make(map[uuid.UUID]*domain.User),
	}
}

func (p *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.repo[user.ID] = user
	return nil
}

func (p *InMemoryUserRepository) GetByID(ctx context.Context, ID uuid.UUID) (*domain.User, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for _, value := range p.repo {
		if value.ID == ID {
			return value, nil
		}
	}
	// return nil, fmt.Errorf("such user {%s} does not exist", ID.String())
	return nil, nil
}

func (p *InMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for _, value := range p.repo {
		if value.Username == username {
			return value, nil
		}
	}
	// return nil, fmt.Errorf("such user username {%s} does not exist", username)
	return nil, nil
}

func (p *InMemoryUserRepository) UpdateRating(ctx context.Context, ID uuid.UUID, delta int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	user, ok := p.repo[ID] // O(1) поиск по uuid-ключу, без цикла!
	if !ok {
		return fmt.Errorf("user %s not found", ID)
	}

	newRating := user.Rating + delta
	if newRating < 0 {
		newRating = 0
	}
	user.Rating = newRating
	return nil
}

func (p *InMemoryUserRepository) GetTopByRating(ctx context.Context, limit int) ([]*domain.User, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	var allUsers []*domain.User

	for _, v := range p.repo {
		allUsers = append(allUsers, v)
	}

	sort.Slice(allUsers, func(i, j int) bool {
		return allUsers[i].Rating > allUsers[j].Rating
	})

	var userWithTopReating []*domain.User
	for i := 0; i < limit && i < len(allUsers); i++ {
		userWithTopReating = append(userWithTopReating, allUsers[i])
	}

	return userWithTopReating, nil
}
