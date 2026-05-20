package domain

import "github.com/google/uuid"

type Quiz struct {
	ID        uuid.UUID
	Title     string
	CreatedBy uuid.UUID
}
type Question struct {
	ID              uuid.UUID
	QuizID          uuid.UUID
	Text            string
	OrderNum        int
	CorrectOptionID uuid.UUID
}
type Option struct {
	Index      int
	ID         uuid.UUID
	QuestionID uuid.UUID
	Text       string
	IsCorrect  bool
}

func NewQuiz(title string, createdPlayerID uuid.UUID) *Quiz {
	return &Quiz{
		ID:        uuid.New(),
		Title:     title,
		CreatedBy: createdPlayerID,
	}
}

func NewQuestion(quizID uuid.UUID, text string, orderNum int, correctOptionID uuid.UUID) *Question {
	return &Question{
		ID:              uuid.New(),
		QuizID:          quizID,
		Text:            text,
		OrderNum:        orderNum,
		CorrectOptionID: correctOptionID,
	}
}

func NewOption(index int, questionID uuid.UUID, text string, isCorrect bool) *Option {
	return &Option{
		Index:      index,
		ID:         uuid.New(),
		QuestionID: questionID,
		Text:       text,
		IsCorrect:  isCorrect,
	}
}

func (qz *Quiz) UpdateTitle(text string) {
	qz.Title = text
}

func (quest *Question) UpdateTextQuestion(text string) {
	quest.Text = text
}
