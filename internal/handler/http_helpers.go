package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func ParseUserID(r *http.Request) (uuid.UUID, error) {
	userID, err := uuid.Parse(r.Header.Get("X-User-ID"))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parse X-User-ID")

	}
	if userID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("userID is empty. userID is required")
	}

	return userID, nil
}
