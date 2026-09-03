package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type GameStatus string

const (
	GameStatusWaiting  GameStatus = "waiting"
	GameStatusActive   GameStatus = "active"
	GameStatusFinished GameStatus = "finished"
)

type PlayerAnswer struct {
	UserID    uuid.UUID `json:"user_id"`
	OptionID  uuid.UUID `json:"option_id"`
	IsCorrect bool      `json:"is_correct"`
	AnswerAt  time.Time `json:"answer_at"`
}

type GameState struct {
	RoomID          uuid.UUID      `json:"room_id"`
	Status          GameStatus     `json:"status"`
	CurrentQuestion int            `json:"current_question"`
	Players         []uuid.UUID    `json:"players"`
	Answers         []PlayerAnswer `json:"answers"`
	CreatedAt       time.Time      `json:"created_at"`
	StartedAt       time.Time      `json:"started_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func NewGameState(roomID uuid.UUID, players []*User) *GameState {
	ids := make([]uuid.UUID, 0, len(players))
	for _, item := range players {
		ids = append(ids, item.ID)
	}
	return &GameState{
		RoomID:          roomID,
		Status:          GameStatusWaiting,
		CurrentQuestion: 0,
		Players:         ids,
		Answers:         make([]PlayerAnswer, 0),
		CreatedAt:       time.Now(),
	}
}

func NewPlayerAnswer(userID, optionID uuid.UUID, isCorrect bool) *PlayerAnswer {
	return &PlayerAnswer{
		UserID:    userID,
		OptionID:  optionID,
		IsCorrect: isCorrect,
		AnswerAt:  time.Now(),
	}
}

func (g *GameState) Start() {
	g.Status = GameStatusActive
	g.StartedAt = time.Now()
	g.UpdatedAt = time.Now()
}

func (g *GameState) Finish() {
	g.Status = GameStatusFinished
	g.UpdatedAt = time.Now()
}

func (g *GameState) ConvertToByte() ([]byte, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (g *GameState) AddAnswer(answer *PlayerAnswer) {
	g.Answers = append(g.Answers, *answer)
	g.UpdatedAt = time.Now()
}

func (g *GameState) String() string {
	data, err := json.Marshal(g)
	if err != nil {
		return fmt.Sprintf("error marshaling GameState: %v", err)
	}
	return string(data)
}

func GameStateFromBytes(data []byte) (*GameState, error) {
	var state GameState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}
