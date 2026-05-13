package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/yourname/quiz-platform/internal/domain"
	"github.com/yourname/quiz-platform/internal/repository/inmemory"
	"github.com/yourname/quiz-platform/internal/service"
)

// TODO реализовать unit-тесты для сервиса GameService
// TestCreateGame_Success +
// TestCreateGame_RoomNotFound +
// TestStartGame_Success +
// TestStartGame_AlreadyStarted +
// TestSubmitAnswer_Success +
// TestSubmitAnswer_GameNotActive +
// TestFinishGame_Success +
// TestFinishGame_GameNotActive +

// ─────────────────────────────────────────────
// Хелпер
// ─────────────────────────────────────────────

func setupGameWithPlayers(t *testing.T, ctx context.Context) (
	*service.GameService,
	*service.RoomService,
	*domain.GameState,
	[]domain.User,
) {
	t.Helper()

	userRepo := inmemory.NewInMemoryUserRepository()
	userServs := service.NewUserService(userRepo)

	names := []struct{ name, username string }{
		{"Alice", "alice"}, {"Fara", "fara"},
		{"Marta", "marta"}, {"Sara", "sara"}, {"Oleg", "oleg"},
	}

	var users []domain.User
	for _, n := range names {
		user, err := userServs.Register(ctx, n.name, n.username)
		if err != nil {
			t.Fatalf("register failed: %v", err)
		}
		users = append(users, *user)
	}

	roomRepo := inmemory.NewInMemoryRoomRepository(userRepo)
	roomServs := service.NewRoomService(userRepo, roomRepo)
	gameRepo := inmemory.NewInMemoryGameRepository()
	gameServs := service.NewGameService(gameRepo, roomRepo, userRepo)

	room, err := roomServs.CreateRoom(ctx, users[0].ID, 6)
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}
	for i := 1; i < len(users); i++ {
		_, err = roomServs.Join(ctx, room.Code, users[i].ID)
		if err != nil {
			t.Fatalf("join failed: %v", err)
		}
	}

	state, err := gameServs.CreateGame(ctx, room.ID)
	if err != nil {
		t.Fatalf("create game failed: %v", err)
	}

	return gameServs, roomServs, state, users
}

// ─────────────────────────────────────────────
// Вспомогательные ответы — общие для нескольких тестов
// ─────────────────────────────────────────────

var testAnswers = []struct {
	answer    int
	isCorrect bool
}{
	{1, true},
	{2, false},
	{3, false},
	{4, false},
	{5, false},
}

// ─────────────────────────────────────────────
// Тесты
// ─────────────────────────────────────────────

func TestCreateGame_Success(t *testing.T) {
	ctx := context.Background()
	_, _, gameState, users := setupGameWithPlayers(t, ctx)

	t.Log(gameState.String())

	if gameState.Status != domain.GameStatusWaiting {
		t.Errorf("expected status %s, got %s", domain.GameStatusWaiting, gameState.Status)
	}
	if len(gameState.Players) != len(users) {
		t.Errorf("expected %d players, got %d", len(users), len(gameState.Players))
	}
}

func TestCreateGame_RoomNotFound(t *testing.T) {
	ctx := context.Background()
	gameServs, _, _, _ := setupGameWithPlayers(t, ctx)

	invalidRoomID := uuid.New()
	_, err := gameServs.CreateGame(ctx, invalidRoomID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := fmt.Sprintf("not found room with id: %s", invalidRoomID)
	if err.Error() != expectedErr {
		t.Errorf("\nunexpected error message:\n got:  %q\n want: %q", err.Error(), expectedErr)
	}
}

func TestStartGame_Success(t *testing.T) {
	ctx := context.Background()
	gameServs, _, gameState, _ := setupGameWithPlayers(t, ctx)

	t.Log(gameState.String())

	updatedState, err := gameServs.StartGame(ctx, gameState.RoomID)
	if err != nil {
		t.Fatalf("error changing status to 'active': %v", err)
	}
	if updatedState.Status != domain.GameStatusActive {
		t.Errorf("expected status %s, got %s", domain.GameStatusActive, updatedState.Status)
	}
}

func TestStartGame_AlreadyStarted(t *testing.T) {
	ctx := context.Background()
	gameServs, _, gameState, _ := setupGameWithPlayers(t, ctx)

	t.Log(gameState.String())

	// Первый старт — должен пройти успешно
	updatedState, err := gameServs.StartGame(ctx, gameState.RoomID)
	if err != nil {
		t.Fatalf("error changing status to 'active': %v", err)
	}
	if updatedState.Status != domain.GameStatusActive {
		t.Errorf("expected status %s, got %s", domain.GameStatusActive, updatedState.Status)
	}

	// Второй старт — должна прийти ошибка
	_, err = gameServs.StartGame(ctx, gameState.RoomID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := fmt.Sprintf(
		"game in room %s already started or finished, status: %s",
		updatedState.RoomID,
		updatedState.Status,
	)
	if err.Error() != expectedErr {
		t.Errorf("\nunexpected error message:\n got:  %q\n want: %q", err.Error(), expectedErr)
	}
}

func TestSubmitAnswer_Success(t *testing.T) {
	ctx := context.Background()
	gameServs, _, gameState, users := setupGameWithPlayers(t, ctx)

	t.Log(gameState.String())

	updatedState, err := gameServs.StartGame(ctx, gameState.RoomID)
	if err != nil {
		t.Fatalf("error changing status to 'active': %v", err)
	}
	if updatedState.Status != domain.GameStatusActive {
		t.Errorf("expected status %s, got %s", domain.GameStatusActive, updatedState.Status)
	}

	for i, player := range users {
		err := gameServs.SubmitAnswer(
			ctx,
			updatedState.RoomID,
			player.ID,
			testAnswers[i].answer,
			testAnswers[i].isCorrect,
		)
		if err != nil {
			t.Fatalf("unexpected error submitting answer for player %s: %v", player.ID, err)
		}
	}
}

func TestSubmitAnswer_GameNotActive(t *testing.T) {
	ctx := context.Background()
	gameServs, _, gameState, users := setupGameWithPlayers(t, ctx)

	t.Log(gameState.String())

	// Намеренно НЕ стартуем игру — статус остаётся Waiting
	expectedErrNotFound := fmt.Sprintf(
		"game state from room %s is not found in redis",
		gameState.RoomID,
	)
	expectedErrNotActive := fmt.Sprintf(
		"game in room %s is not active, status: %s",
		gameState.RoomID,
		gameState.Status,
	)

	for i, player := range users {
		err := gameServs.SubmitAnswer(
			ctx,
			gameState.RoomID,
			player.ID,
			testAnswers[i].answer,
			testAnswers[i].isCorrect,
		)
		if err == nil {
			t.Fatalf("expected error from SubmitAnswer, got nil")
		}

		got := err.Error()
		if got != expectedErrNotFound && got != expectedErrNotActive {
			t.Errorf(
				"\nunexpected error message:\n got:  %q\n want one of:\n  %q\n  %q",
				got,
				expectedErrNotFound,
				expectedErrNotActive,
			)
		}
	}
}

func TestFinishGame_Success(t *testing.T) {
	ctx := context.Background()
	gameServs, _, gameState, users := setupGameWithPlayers(t, ctx)

	t.Log(gameState.String())

	updatedState, err := gameServs.StartGame(ctx, gameState.RoomID)
	if err != nil {
		t.Fatalf("error changing status to 'active': %v", err)
	}
	if updatedState.Status != domain.GameStatusActive {
		t.Errorf("expected status %s, got %s", domain.GameStatusActive, updatedState.Status)
	}

	for i, player := range users {
		err := gameServs.SubmitAnswer(
			ctx,
			updatedState.RoomID,
			player.ID,
			testAnswers[i].answer,
			testAnswers[i].isCorrect,
		)
		if err != nil {
			t.Fatalf("error submitting answer for player %s: %v", player.ID, err)
		}
	}

	finishedState, err := gameServs.FinishGame(ctx, updatedState.RoomID)
	if err != nil {
		t.Fatalf("error finishing game: %v", err)
	}
	if finishedState.Status != domain.GameStatusFinished {
		t.Errorf("expected status %s, got %s", domain.GameStatusFinished, finishedState.Status)
	}
	if finishedState.RoomID != gameState.RoomID {
		t.Errorf("expected room %s, got %s", gameState.RoomID, finishedState.RoomID)
	}
}

func TestFinishGame_GameNotActive(t *testing.T) {
	ctx := context.Background()
	gameServs, _, gameState, _ := setupGameWithPlayers(t, ctx)

	t.Log(gameState.String())

	// Намеренно НЕ стартуем игру — статус остаётся Waiting
	_, err := gameServs.FinishGame(ctx, gameState.RoomID)
	if err == nil {
		t.Fatal("expected error from FinishGame, got nil")
	}

	expectedErr := fmt.Sprintf(
		"game in room %s is not active, status: %s",
		gameState.RoomID,
		gameState.Status,
	)
	if err.Error() != expectedErr {
		t.Errorf("\nunexpected error message:\n got:  %q\n want: %q", err.Error(), expectedErr)
	}
}
