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
	if user.Reating != 0 {
		t.Errorf("expected rating 1000, got: %d", user.Reating)
	}
}

func TestGetUsersWithTopReating_Success(t *testing.T) {
	repo := inmemory.NewInMemoryUserRepository()
	service := service.NewUserService(repo)
	ctx := context.Background()

	var users []*domain.User
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("Test-bob-00%d", i)
		username := fmt.Sprintf("Test-username-bob-00%d", i)
		user, err := service.Register(ctx, name, username)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		users = append(users, user)
	}

	for i, value := range users {
		rating := 1001 + i
		value.Reating = rating
		err := service.UpdateReating(ctx, value.Id, rating)
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
		if v.Reating != expectedVal[index] {
			t.Errorf("expected reating at userName : %s with reating %d, got: %d", v.Username, expectedVal[index], v.Reating)
		}
	}
}
