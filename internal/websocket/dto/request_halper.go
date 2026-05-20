package dto

import "github.com/google/uuid"

type BaseRequest struct {
	Type string `json:"type"`
}
type AnswerRequsest struct {
	QuestionID uuid.UUID `json:"question_id"`
	OptionID   uuid.UUID `json:"option_id"`
}
