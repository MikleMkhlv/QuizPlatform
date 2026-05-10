package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID
	Name       string
	Username   string
	Reating    int
	CrteatedAt time.Time
}

func NewUser(name string, username string) *User {
	return &User{
		Id:         uuid.New(),
		Name:       name,
		Username:   username,
		Reating:    0,
		CrteatedAt: time.Now(),
	}
}

// func (u *User) updateReating(count int) {
// 	newReating := u.Reating + count
// 	if newReating < 0 {
// 		u.Reating = 0
// 	} else {
// 		u.Reating = newReating
// 	}
// }
