package websocket

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	rwMutex sync.RWMutex
	rooms   map[uuid.UUID]map[uuid.UUID]*websocket.Conn
}

func NewWebsocketHub() *Hub {
	return &Hub{
		rooms: make(map[uuid.UUID]map[uuid.UUID]*websocket.Conn),
	}
}

func (h *Hub) AddPlayer(roomID uuid.UUID, playerID uuid.UUID, conn *websocket.Conn) error {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[uuid.UUID]*websocket.Conn)
	}
	h.rooms[roomID][playerID] = conn
	return nil
}

func (h *Hub) RemovePlayer(roomID uuid.UUID, playerID uuid.UUID) error {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()
	if room, ok := h.rooms[roomID]; ok {
		if conn, ok := room[playerID]; ok {
			conn.Close()
			delete(room, playerID)
		} else {
			return fmt.Errorf("not found player in hub")
		}
	} else {
		return fmt.Errorf("not found room in hub")
	}
	return nil

}

func (h *Hub) SendMessage(roomID uuid.UUID, message []byte) error {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()
	isError := false
	for id, conn := range h.rooms[roomID] {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("error write message for playerID %s: %v", id, err)
			isError = true
		}
	}
	if isError {
		return fmt.Errorf("there were errors sending the message")
	}
	return nil
}
