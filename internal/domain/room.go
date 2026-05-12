package domain

import (
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
)

type RoomStatus string

const (
	RoomStatusWaiting  RoomStatus = "waiting"
	RoomStatusActive   RoomStatus = "active"
	RoomStatusFinished RoomStatus = "finished"
)

type Room struct {
	ID        uuid.UUID
	HostID    uuid.UUID
	Code      string
	Status    RoomStatus
	MaxPlayer int
	CreatedAt time.Time
}

type RoomPlayer struct {
	RoomID   uuid.UUID
	PlayerID uuid.UUID
	JoinedAt time.Time
}

func NewRoom(hostId uuid.UUID, countPlayer int) *Room {
	return &Room{
		ID:        uuid.New(),
		HostID:    hostId,
		Code:      generateCodeFromRoom(),
		Status:    RoomStatusWaiting,
		MaxPlayer: countPlayer,
		CreatedAt: time.Now(),
	}
}

func NewRoomPlayer(roomId, plauerId uuid.UUID) *RoomPlayer {
	return &RoomPlayer{
		RoomID:   roomId,
		PlayerID: plauerId,
		JoinedAt: time.Now(),
	}
}

func generateCodeFromRoom() string {
	const charSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charSet[rand.IntN(len(charSet))] // берём символ по индексу
	}
	return string(code)
}
