package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
	"github.com/MikleMkhlv/QuizPlatform/internal/ports"
	"github.com/MikleMkhlv/QuizPlatform/internal/websocket/dto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub       *Hub
	userServ  ports.UserService
	roomServ  ports.RoomService
	gameServ  ports.GameServiceI
	questServ ports.QuestionService
	wsUpg     *websocket.Upgrader
	// mu        sync.Mutex
}

func NewWSHandler(userServ ports.UserService, roomServ ports.RoomService, gameServ ports.GameServiceI, questServ ports.QuestionService, hub *Hub) *WSHandler {
	return &WSHandler{
		userServ:  userServ,
		roomServ:  roomServ,
		gameServ:  gameServ,
		questServ: questServ,
		hub:       hub,
		wsUpg: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (wsh *WSHandler) Connect(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	ctx := r.Context()

	roomCode := r.URL.Query().Get("room_code")
	if roomCode == "" {
		log.Println("room_code param is requaired")
		http.Error(w, "roomID param is requaired", http.StatusBadRequest)
		return
	}
	playerID, err := uuid.Parse(r.URL.Query().Get("playerID"))
	if err != nil {
		log.Printf("error parse playerID param: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if playerID == uuid.Nil {
		log.Printf("playerID param is requaired")
		http.Error(w, "playerID param is requaired", http.StatusBadRequest)
		return
	}

	room, err := wsh.roomServ.GetRoomByRoomCode(ctx, roomCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if room == nil {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
	if room.Status == domain.RoomStatusFinished {
		errMsg := fmt.Errorf("game in room {%s} is finished", room.ID)
		http.Error(w, errMsg.Error(), http.StatusConflict)
		return
	}
	conn, err := wsh.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	if err := wsh.hub.AddPlayer(room.ID, playerID, conn); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	defer wsh.hub.RemovePlayer(room.ID, playerID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var reqMsg dto.BaseRequest
		if err := json.Unmarshal(message, &reqMsg); err != nil {
			log.Printf("error unmarshal request: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("request body for baseType is not correct"))
			// http.Error(w, "request body for baseType is not correct", http.StatusBadRequest)
			continue
		}

		room, err := wsh.roomServ.GetRoomByRoomCode(ctx, roomCode)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("room not found"))
			break
		}

		switch reqMsg.Type {
		case "answer":
			if room.Status != domain.RoomStatusActive {
				conn.WriteMessage(websocket.TextMessage, []byte("game in not active"))
				continue
			}
			var answerReq dto.AnswerRequsest
			if err := json.Unmarshal(message, &answerReq); err != nil {
				log.Printf("error unmarshal request: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte("request body for answer is not correct"))
				continue
			}
			log.Printf("answer received: questionID=%s optionID=%s", answerReq.QuestionID, answerReq.OptionID)

			if err := wsh.gameServ.SubmitAnswer(ctx, room.ID, answerReq.QuestionID, playerID, answerReq.OptionID); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				continue
			}
		case "create":
			if room.Status == domain.RoomStatusWaiting && room.HostID == playerID {
				_, err = wsh.gameServ.CreateGame(ctx, room.ID)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
					continue
				}
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte("attempt to create a game without being the host"))
				continue
			}
			conn.WriteMessage(websocket.TextMessage, []byte("game is created"))
		case "start":
			if room.Status != domain.RoomStatusWaiting {
				conn.WriteMessage(websocket.TextMessage, []byte("game already started"))
				continue
			}
			if room.HostID == playerID {
				_, err = wsh.gameServ.StartGame(ctx, room.ID)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
					continue
				}
				if err := wsh.roomServ.UpdateRoomStatus(ctx, room.ID, domain.RoomStatusActive); err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
					continue
				}
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte("attempt to start game without being the host"))
				continue
			}

		default:
			conn.WriteMessage(websocket.TextMessage, []byte("type request is not correct"))
			// http.Error(w, "type request is not correct", http.StatusBadRequest)
			continue
		}
	}
}

// func (wsh *WSHandler) handler(conn *websocket.Conn, ctx context.Context) {
// }
