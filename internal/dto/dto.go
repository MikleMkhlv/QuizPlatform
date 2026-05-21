package dto

import "github.com/MikleMkhlv/QuizPlatform/internal/domain"

type QuestionWithOptions struct {
	Question domain.Question
	Options  []domain.Option
}

func NewQuestionWithOptions(question domain.Question, options ...domain.Option) *QuestionWithOptions {
	return &QuestionWithOptions{
		Question: question,
		Options:  options,
	}
}
