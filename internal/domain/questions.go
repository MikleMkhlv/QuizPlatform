package domain

import "github.com/google/uuid"

type Quizzes struct {
	ID         uuid.UUID
	Title      string
	Created_by uuid.UUID
}

type Questions struct {
	ID        uuid.UUID
	Quiz_id   uuid.UUID
	Text      string
	Order_num int
}

type Options struct {
	ID          uuid.UUID
	Question_id uuid.UUID
	Text        string
	IsCorrect   bool
}

func NewQuiz(title string, createdPlayerID uuid.UUID) *Quizzes {
	return &Quizzes{
		ID:         uuid.New(),
		Title:      title,
		Created_by: createdPlayerID,
	}
}

func NewQuestion(quizID uuid.UUID, text string, orderNum int) *Questions {
	return &Questions{
		ID:        uuid.New(),
		Quiz_id:   quizID,
		Text:      text,
		Order_num: orderNum,
	}
}

func NewOption(questionID uuid.UUID, text string, isCorrect bool) *Options {
	return &Options{
		ID:          uuid.New(),
		Question_id: questionID,
		Text:        text,
		IsCorrect:   isCorrect,
	}
}

func (qz *Quizzes) UpdateTitle(text string) {
	qz.Title = text
}

func (quest *Questions) UpdateTextQuestion(text string) {
	quest.Text = text
}
