package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Username  string
	Rating    int
	CreatedAt time.Time
}

func NewUser(name string, username string) *User {
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Username:  username,
		Rating:    0,
		CreatedAt: time.Now(),
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
