package main

import (
	"fmt"
	"net/http"

	"github.com/yourname/quiz-platform/internal/handler"
	"github.com/yourname/quiz-platform/internal/repository/inmemory"
	"github.com/yourname/quiz-platform/internal/service"
)

func main() {
	userRepo := inmemory.NewInMemoryUserRepository()
	roomRepo := inmemory.NewInMemoryRoomRepository(userRepo)
	// gameRepo := inmemory.NewInMemoryGameRepository()

	userServ := service.NewUserService(userRepo)
	roomServ := service.NewRoomService(userRepo, roomRepo)

	userHandl := handler.NewUserHandler(userServ)
	roomHandl := handler.NewRoomHandler(roomServ)
	registerRoutes(userHandl, roomHandl)

	fmt.Println("Started server with port 8080.")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func registerRoutes(userHandl *handler.UserHandler, roomHandl *handler.RoomHandler) {

	http.HandleFunc("POST /api/users/register", userHandl.Register)

	http.HandleFunc("POST /api/rooms", roomHandl.CreateRoom)
	http.HandleFunc("POST /api/rooms/join", roomHandl.JoinInRoom)
}
