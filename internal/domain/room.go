package domain

import (
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
)

type RoomStatus string

const (
	RoomStatusWating   RoomStatus = "wating"
	RoomStatusActive   RoomStatus = "active"
	RoomStatusFinished RoomStatus = "finished"
)

type Room struct {
	ID        uuid.UUID
	HostID    uuid.UUID
	Code      string
	Status    RoomStatus
	MaxPlauer int
	CreatedAt time.Time
}

type RoomPlauer struct {
	RoomId   uuid.UUID
	PlauerId uuid.UUID
	JoinedAt time.Time
}

func NewRoom(hostId uuid.UUID, countPlauers int) *Room {
	return &Room{
		ID:        uuid.New(),
		HostID:    hostId,
		Code:      generateCodeFromRoom(),
		Status:    RoomStatusWating,
		MaxPlauer: countPlauers,
		CreatedAt: time.Now(),
	}
}

func NewRoomPlauer(roomId, plauerId uuid.UUID) *RoomPlauer {
	return &RoomPlauer{
		RoomId:   roomId,
		PlauerId: plauerId,
		JoinedAt: time.Now(),
	}
}

func generateCodeFromRoom() string {
	charSet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+")
	length := 6
	code := make([]rune, length)
	for i := range code {
		num := rand.IntN(len(charSet) + 1)
		code[i] = rune(num)
	}
	return string(code)
}
