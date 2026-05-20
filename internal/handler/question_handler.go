package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MikleMkhlv/QuizPlatform/internal/ports"
	"github.com/google/uuid"
)

type QuestionHandler struct {
	qestionServs ports.QuestionService
}

type QuizCreateRequest struct {
	PlayerID uuid.UUID `json:"player_id"`
	Title    string    `json:"title"`
}

type QuizCreateResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

func NewQuestionHandler(qestionServs ports.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		qestionServs: qestionServs,
	}
}

func (q *QuestionHandler) CreateQuizz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	var req QuizCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decode request body", http.StatusBadRequest)
		return
	}

	quizz, err := q.qestionServs.CreateQuizz(ctx, req.PlayerID, req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := QuizCreateResponse{
		ID:    quizz.ID,
		Title: quizz.Title,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("error marshal QuestionCreateResponse: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

type FindedQuizResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedBy uuid.UUID `json:"created_by"`
}

func (q *QuestionHandler) GetQuizzByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	quizzID, err := uuid.Parse(r.URL.Query().Get("quizId"))
	if err != nil {
		http.Error(w, "error parse quizID param", http.StatusInternalServerError)
		return
	}
	if quizzID == uuid.Nil {
		http.Error(w, "quizID is required", http.StatusBadRequest)
		return
	}

	foundQuiz, err := q.qestionServs.GetQuizzByID(ctx, quizzID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := FindedQuizResponse{
		ID:        foundQuiz.ID,
		Title:     foundQuiz.Title,
		CreatedBy: foundQuiz.CreatedBy,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("error marshal FindedQuizResponse: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

type QuestionWithOptionsReruest struct {
	Text    string `json:"text"`
	Options []ports.OptionRequest
}

func (q *QuestionHandler) AddNewQuestionswithOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	quizzID, err := uuid.Parse(r.URL.Query().Get("quizId"))
	if err != nil {
		http.Error(w, "error parse quizID param", http.StatusInternalServerError)
		return
	}
	if quizzID == uuid.Nil {
		http.Error(w, "quizID is required", http.StatusBadRequest)
		return
	}

	var reqQuestion QuestionWithOptionsReruest
	if err := json.NewDecoder(r.Body).Decode(&reqQuestion); err != nil {
		http.Error(w, "error decode request body", http.StatusBadRequest)
		return
	}

	if err := q.qestionServs.AddNewQuestionsWithAnswers(ctx, quizzID, reqQuestion.Text, reqQuestion.Options); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "write questions is sucessful"}); err != nil {
		log.Println("error send response")
		return
	}
}
