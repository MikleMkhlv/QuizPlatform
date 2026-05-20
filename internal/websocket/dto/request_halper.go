package dto

import "github.com/google/uuid"

type BaseRequest struct {
	Type string `json:"type"`
}
type AnsverRequsest struct {
	QuizID uuid.UUID `json:"quizID"`
	Answer int       `json:"answer"`
}
