package domains

import (
	"github.com/google/uuid"
	"time"
)

type Group struct {
	ID        uuid.UUID
	Name      string
	Members   []Member
	UpdatedAt time.Time
	CreatedAt time.Time
}

type Member struct {
	GroupID   uuid.UUID
	ID        uuid.UUID
	UserID    uuid.UUID
	Balance   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
