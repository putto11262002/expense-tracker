package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	FirstName string
	LastName  string
	Username  string `gorm:"uniqueIndex;size:100"`
	Email     string `gorm:"uniqueIndex;size:100"`
	Password  string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewUser(firstName, lastName, username, email, password string) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Password:  password,
		Email:     email,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}
