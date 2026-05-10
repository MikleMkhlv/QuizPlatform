package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
)

type inMemoryUserRepository struct {
	repo  map[int]*domain.User
	mutex sync.RWMutex
}

var Instanse *inMemoryUserRepository

func NewInMemoryUserRepository() *inMemoryUserRepository {
	if Instanse != nil {
		return Instanse
	}
	Instanse = &inMemoryUserRepository{repo: make(map[int]*domain.User)}
	return Instanse
}

func (p *inMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.repo[len(p.repo)+1] = user
	return nil
}

func (p *inMemoryUserRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for _, value := range p.repo {
		if value.Id == id {
			return value, nil
		}
	}
	return nil, fmt.Errorf("such user {%s} does not exist\n", id.String())
}

func (p *inMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for _, value := range p.repo {
		if value.Username == username {
			return value, nil
		}
	}
	return nil, fmt.Errorf("such user username {%s} does not exist\n", username)
}

func (p *inMemoryUserRepository) UpdateReating(ctx context.Context, id uuid.UUID, delta int) error {
	existing, err := p.GetById(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("such user {%s} does not exist\n", id.String())
	}

	newReating := existing.Reating + delta
	if newReating < 0 {
		newReating = 0
	}

	existing.Reating = newReating

	p.mutex.Lock()
	defer p.mutex.Unlock()
	for k, v := range p.repo {
		if v.Id == existing.Id {
			p.repo[k] = existing
			fmt.Printf("id: %s - reating updated\n", existing.Id)
		}
	}
	return nil
}
