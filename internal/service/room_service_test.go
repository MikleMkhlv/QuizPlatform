package service_test

import (
	"context"
	"testing"

	"github.com/yourname/quiz-platform/internal/domain"
	"github.com/yourname/quiz-platform/internal/repository/inmemory"
	"github.com/yourname/quiz-platform/internal/service"
	_ "github.com/yourname/quiz-platform/internal/service"
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

	roomRepo := inmemory.NewInMemoryRoomRepository()
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
// TestCreateRoom_Success          — happy path
// TestCreateRoom_UserNotFound     — хост не существует
// TestJoinRoom_Success            — игрок заходит в комнату
// TestJoinRoom_RoomNotFound       — неверный код комнаты
// TestJoinRoom_RoomFull           — комната заполнена
// TestJoinRoom_WrongStatus        — игра уже началась

func TestAddNewPlayerInRoomByRoomID_Success(t *testing.T) {

}
