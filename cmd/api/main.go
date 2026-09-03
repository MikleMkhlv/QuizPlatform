package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MikleMkhlv/QuizPlatform/internal/handler"
	"github.com/MikleMkhlv/QuizPlatform/internal/repository/postgres"
	"github.com/MikleMkhlv/QuizPlatform/internal/repository/rediss"
	"github.com/MikleMkhlv/QuizPlatform/internal/service"
	"github.com/MikleMkhlv/QuizPlatform/internal/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {

	connectionString := "postgres://postgres:postgres@localhost:5555/postgres"

	// Создаем пул соединений
	dbPool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось подключиться к базе данных: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	err = dbPool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось пропинговать базу данных: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Успешное подключение к PostgreSQL!")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес Redis сервера
		Password: "",               // Пароль (пустой, если не задан)
		DB:       0,                // Номер базы данных по умолчанию
	})
	defer rdb.Close()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	fmt.Printf("Успешное подключение к Redis! Ответ сервера: %s\n", pong)

	// dataBasePool := initDB()
	// redis := initRedis()
	userRepo := postgres.NewPostgresUserRepository(dbPool)
	roomRepo := postgres.NewPostgresRoomRepository(dbPool)
	gameRepo := rediss.NewRedisRepository(rdb)
	quizRepo := postgres.NewPostgresQuestionsRepository(dbPool)

	userServ := service.NewUserService(userRepo)
	roomServ := service.NewRoomService(userRepo, roomRepo)
	gameServ := service.NewGameService(gameRepo, roomRepo, userRepo, quizRepo)
	quizServ := service.NewQusetionService(quizRepo, userServ)

	userHandl := handler.NewUserHandler(userServ)
	roomHandl := handler.NewRoomHandler(roomServ)

	wsHub := websocket.NewWebsocketHub()
	wsHandl := websocket.NewWSHandler(userServ, roomServ, gameServ, quizServ, wsHub)
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

func initDB() *pgxpool.Pool {

	connectionString := "postgres://postgres:postgres@localhost:5555/postgres"

	// Создаем пул соединений
	dbPool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось подключиться к базе данных: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	err = dbPool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось пропинговать базу данных: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Успешное подключение к PostgreSQL!")

	return dbPool
}

func initRedis() *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес Redis сервера
		Password: "",               // Пароль (пустой, если не задан)
		DB:       0,                // Номер базы данных по умолчанию
	})
	defer rdb.Close()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	fmt.Printf("Успешное подключение к Redis! Ответ сервера: %s\n", pong)

	return rdb
}
