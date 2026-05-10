package service_test

import (
	"context"
	"testing"

	"github.com/yourname/quiz-platform/internal/repository/inmemory"
	"github.com/yourname/quiz-platform/internal/service"
)

func TestGetUsersWithTopRating_Success(t *testing.T) {
	repo := inmemory.NewInMemoryUserRepository()
	svc := service.NewUserService(repo)
	ctx := context.Background()

	testCases := []struct {
		name     string
		username string
		rating   int
	}{
		{"Alice", "alice", 500},
		{"Bob", "bob", 1500},
		{"Charlie", "charlie", 1000},
		{"Dave", "dave", 2000},
		{"Eve", "eve", 750},
	}

	// Register сразу возвращает *User с ID — используем его!
	for _, tc := range testCases {
		user, err := svc.Register(ctx, tc.name, tc.username)
		if err != nil {
			t.Fatalf("register failed: %v", err)
		}

		// Сразу устанавливаем рейтинг — никакого лишнего запроса
		err = svc.UpdateRating(ctx, user.ID, tc.rating)
		if err != nil {
			t.Fatalf("update rating failed: %v", err)
		}
	}

	topUsers, err := svc.GetTopUsers(ctx, 3)
	if err != nil {
		t.Fatalf("get top users failed: %v", err)
	}

	if len(topUsers) != 3 {
		t.Fatalf("expected 3 users, got %d", len(topUsers))
	}

	expected := []struct {
		username string
		rating   int
	}{
		{"dave", 2000},
		{"bob", 1500},
		{"charlie", 1000},
	}

	for i, exp := range expected {
		if topUsers[i].Username != exp.username {
			t.Errorf("position %d: expected username %s, got %s",
				i, exp.username, topUsers[i].Username)
		}
		if topUsers[i].Rating != exp.rating {
			t.Errorf("position %d: expected rating %d, got %d",
				i, exp.rating, topUsers[i].Rating)
		}
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	repo := inmemory.NewInMemoryUserRepository()
	svc := service.NewUserService(repo)
	ctx := context.Background()

	_, err := svc.Register(ctx, "Fred", "fredornor")
	if err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	_, err = svc.Register(ctx, "Another Fred", "fredornor")
	if err == nil {
		t.Fatal("expected error for duplicate username, got nil")
	}
}
