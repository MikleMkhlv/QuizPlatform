package main

import (
	"fmt"
	"net/http"

	"github.com/MikleMkhlv/QuizPlatform/internal/handler"
	"github.com/MikleMkhlv/QuizPlatform/internal/repository/inmemory"
	"github.com/MikleMkhlv/QuizPlatform/internal/service"
	"github.com/MikleMkhlv/QuizPlatform/internal/websocket"
)

func main() {
	userRepo := inmemory.NewInMemoryUserRepository()
	roomRepo := inmemory.NewInMemoryRoomRepository(userRepo)
	gameRepo := inmemory.NewInMemoryGameRepository()

	userServ := service.NewUserService(userRepo)
	roomServ := service.NewRoomService(userRepo, roomRepo)
	gameServ := service.NewGameService(gameRepo, roomRepo, userRepo)

	userHandl := handler.NewUserHandler(userServ)
	roomHandl := handler.NewRoomHandler(roomServ)

	wsHub := websocket.NewWebsocketHub()
	wsHandl := websocket.NewWSHandler(userServ, roomServ, gameServ, wsHub)
	registerRoutes(userHandl, roomHandl, wsHandl)

	fmt.Println("Started server with port 8080.")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func registerRoutes(userHandl *handler.UserHandler, roomHandl *handler.RoomHandler, wsHandl *websocket.WSHandler) {

	http.HandleFunc("POST /api/users/register", userHandl.Register)

	http.HandleFunc("POST /api/rooms", roomHandl.CreateRoom)
	http.HandleFunc("POST /api/rooms/join", roomHandl.JoinInRoom)

	http.HandleFunc("GET /ws", wsHandl.Connect)
}
