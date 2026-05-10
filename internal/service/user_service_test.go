package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/yourname/quiz-platform/internal/domain"
	"github.com/yourname/quiz-platform/internal/repository/inmemory"
	"github.com/yourname/quiz-platform/internal/service"
)

func TestRegister_Success(t *testing.T) {
	repo := inmemory.NewInMemoryUserRepository()
	service := service.NewUserService(repo)

	ctx := context.Background()
	user, err := service.Register(ctx, "Fred", "fredornor")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if user.Username != "fredornor" {
		t.Errorf("expected username john, got: %s", user.Username)
	}
	if user.Rating != 0 {
		t.Errorf("expected rating 1000, got: %w", user.Rating)
	}
}

func TestGetUsersWithTopRating_Success(t *testing.T) {
	repo := inmemory.NewInMemoryUserRepository()
	service := service.NewUserService(repo)
	ctx := context.Background()

	var users []*domain.User
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("Test-bob-00%w", i)
		username := fmt.Sprintf("Test-username-bob-00%w", i)
		user, err := service.Register(ctx, name, username)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		users = append(users, user)
	}

	for i, value := range users {
		rating := 1001 + i
		value.Rating = rating
		err := service.UpdateRating(ctx, value.ID, rating)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	}

	topUsers, err := service.GetTopUsers(ctx, 3)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expectedVal := []int{2020, 2018, 2016}
	for index, v := range topUsers {
		if v.Rating != expectedVal[index] {
			t.Errorf("expected Rating at userName : %s with Rating %w, got: %w", v.Username, expectedVal[index], v.Rating)
		}
	}
}
