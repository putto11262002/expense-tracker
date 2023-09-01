package domains

import (
	"time"

	"github.com/google/uuid"
)

// user struct
type User struct {
	ID uuid.UUID 
	FirstName string
	LastName string
	Username string
	Email string
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TODO - hash password
// constructor for users
func NewUser(firstName, lastName, username, email, password string) (*User, error) {
	return &User{
		FirstName: firstName,
		LastName: lastName,
		Username: username,
		Password: password,
		Email: email,
	}, nil 
}