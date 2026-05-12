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

func TestCreateNewRoom_Success(t *testing.T) {
	userRepo := inmemory.NewInMemoryUserRepository()
	userServs := service.NewUserService(userRepo)
	ctx := context.Background()

	testCases := []struct {
		name     string
		username string
		rating   int
	}{
		{"Alice", "alice", 500},
	}

	var expectUser *domain.User
	var err error
	for _, tc := range testCases {
		expectUser, err = userServs.Register(ctx, tc.name, tc.username)
		if err != nil {
			t.Fatalf("register failed: %v", err)
		}
	}

	roomRepo := inmemory.NewInMemoryRoomRepository(userRepo)
	roomsServs := service.NewRoomService(userRepo, roomRepo)

	room, err := roomsServs.CreateRoom(ctx, expectUser.ID, 4)
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}

	if room.HostID != expectUser.ID {
		t.Errorf("The created room {%s} does not belong to the host {%s:%s}", room.ID, expectUser.Username, expectUser.ID)
	}
}

// TODO реализовать тесты
// TestCreateRoom_Success          — happy path +
// TestCreateRoom_UserNotFound     — хост не существует +
// TestJoinRoom_Success            — игрок заходит в комнату +
// TestJoinRoom_RoomNotFound       — неверный код комнаты +
// TestJoinRoom_RoomFull           — комната заполнена +- (Пропускаем, пока не реализованы методы в inmemory некоторые методы)
// TestJoinRoom_WrongStatus        — игра уже началась пока нет функции в room service, которая меняет статус.

func TestCreateRoom_UserNotFound(t *testing.T) {
	userRepo := inmemory.NewInMemoryUserRepository()
	roomRepo := inmemory.NewInMemoryRoomRepository(userRepo)
	roomsServs := service.NewRoomService(userRepo, roomRepo)

	ctx := context.Background()
	notExistPlayerId := uuid.New()

	room, err := roomsServs.CreateRoom(ctx, notExistPlayerId, 4)

	// 1. Ошибка должна быть — если её нет, тест провален
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := fmt.Sprintf("user with id: %s not found", notExistPlayerId.String())
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message:\n got:  %q\n want: %q", err.Error(), expectedErr)
	}

	// 3. Комната не должна была создаться
	if room != nil {
		t.Errorf("expected room to be nil, got: %+v", room)
	}
}

func TestJoinRoom_Success(t *testing.T) {
	userRepo := inmemory.NewInMemoryUserRepository()
	userServs := service.NewUserService(userRepo)
	ctx := context.Background()

	testCases := []struct {
		name     string
		username string
		rating   int
	}{
		{"Alice", "alice", 500},
		{"Bob", "bobrito-bondito", 1000000},
	}

	var expectUser []domain.User
	for _, tc := range testCases {
		user, err := userServs.Register(ctx, tc.name, tc.username)
		if err != nil {
			t.Fatalf("register failed: %v", err)
		}
		expectUser = append(expectUser, *user)
	}

	roomRepo := inmemory.NewInMemoryRoomRepository(userRepo)
	roomsServs := service.NewRoomService(userRepo, roomRepo)

	room, err := roomsServs.CreateRoom(ctx, expectUser[0].ID, 4)
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}

	joinedRoom, err := roomsServs.Join(ctx, room.Code, expectUser[1].ID)
	if err != nil {
		t.Fatalf("join failed: %v", err)
	}
	if joinedRoom.ID != room.ID {
		t.Errorf("expected room %s, got %s", room.ID, joinedRoom.ID)
	}
	if joinedRoom.Code != room.Code {
		t.Errorf("expected code %s, got %s", room.Code, joinedRoom.Code)
	}
}

func TestJoinRoom_RoomNotFound(t *testing.T) {
	userRepo := inmemory.NewInMemoryUserRepository()
	roomRepo := inmemory.NewInMemoryRoomRepository(userRepo)
	roomsServs := service.NewRoomService(userRepo, roomRepo)
	ctx := context.Background()

	failedCodeRoom := "test_fail-code-12312412"
	expectedErr := fmt.Sprintf("room with code %s not found", failedCodeRoom)
	_, err := roomsServs.Join(ctx, failedCodeRoom, uuid.New())
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message:\n got:  %q\n want: %q", err.Error(), expectedErr)
	}
}

func TestJoinRoom_RoomFull(t *testing.T) {
	userRepo := inmemory.NewInMemoryUserRepository()
	userServs := service.NewUserService(userRepo)
	ctx := context.Background()

	testCases := []struct {
		name     string
		username string
		rating   int
	}{
		{"Alice", "alice", 500},
		{"Alice-2", "alice-2", 200},
	}

	var expectUser []domain.User
	for _, tc := range testCases {
		user, err := userServs.Register(ctx, tc.name, tc.username)
		if err != nil {
			t.Fatalf("join in room failed: %v", err)
		}
		expectUser = append(expectUser, *user)

	}

	roomRepo := inmemory.NewInMemoryRoomRepository(userRepo)
	roomsServs := service.NewRoomService(userRepo, roomRepo)

	countUserInRoom := 2
	room, err := roomsServs.CreateRoom(ctx, expectUser[0].ID, countUserInRoom)
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}

	user2 := expectUser[1]
	_, err = roomsServs.Join(ctx, room.Code, user2.ID)
	if err != nil {
		t.Fatalf("join in room failed: %v", err)
	}

	testCasesForNewUser := struct {
		name     string
		username string
	}{
		"new-test_user", "boblared",
	}

	newUser, err := userServs.Register(ctx, testCasesForNewUser.name, testCasesForNewUser.username)
	if err != nil {
		t.Fatalf("join in room failed: %v", err)
	}

	expectedErr := fmt.Sprintf("room %s is full (%d/%d)",
		room.ID,
		countUserInRoom,
		room.MaxPlayer)
	_, err = roomsServs.Join(ctx, room.Code, newUser.ID)
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message:\n got:  %q\n want: %q", err.Error(), expectedErr)
	}
}

func TestJoinRoom_WrongStatus(t *testing.T) {
	userRepo := inmemory.NewInMemoryUserRepository()
	userServs := service.NewUserService(userRepo)
	ctx := context.Background()

	testCases := []struct {
		name     string
		username string
		rating   int
	}{
		{"Alice", "alice", 500},
		{"Alice-2", "alice-2", 200},
	}

	var expectUser []domain.User
	for _, tc := range testCases {
		user, err := userServs.Register(ctx, tc.name, tc.username)
		if err != nil {
			t.Fatalf("join in room failed: %v", err)
		}
		expectUser = append(expectUser, *user)

	}

	roomRepo := inmemory.NewInMemoryRoomRepository(userRepo)
	roomsServs := service.NewRoomService(userRepo, roomRepo)

	countUserInRoom := 2
	room, err := roomsServs.CreateRoom(ctx, expectUser[0].ID, countUserInRoom)
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}

	err = roomsServs.UpdateRoomStatus(ctx, room.ID, domain.RoomStatusActive)
	updatedRoom, err := roomsServs.GetRoomByID(ctx, room.ID)
	if err != nil {
		t.Fatalf("get room failed: %v", err)
	}

	user2 := expectUser[1]
	_, err = roomsServs.Join(ctx, room.Code, user2.ID)
	if err == nil {
		t.Fatalf("join in room failed: %v", err)
	}
	expectedErr := fmt.Sprintf("room %s is not accepting players, status: %s", updatedRoom.ID, updatedRoom.Status)
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message:\n got:  %q\n want: %q", err.Error(), expectedErr)
	}

}
