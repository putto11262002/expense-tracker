package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Group struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	Name      string
	Members   []User `gorm:"many2many:user_groups;"`
	OwnerID   uuid.UUID
	Owner     User      `gorm:"foreignKey:OwnerID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewGroup(name string, owner User) *Group {
	return &Group{
		Name:    name,
		OwnerID: owner.ID,
		Owner:   owner,
		Members: []User{owner},
	}
}

func (g *Group) BeforeCreate(tx *gorm.DB) (err error) {
	g.ID = uuid.New()
	return
}
