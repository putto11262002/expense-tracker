package repositories

import (
	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/internal/domains"
	"gorm.io/gorm"
)

type UserRepsotory struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *UserRepsotory {
	return &UserRepsotory{
		DB: DB,
	}
}


type IUserRepsotory interface {
	CreateUser(*domains.User) (*domains.User, error)
	DeleteUserByID(uuid.UUID) (error)
	GetUserByID(uuid.UUID) (*domains.User, error)
	GetUserByUsername(uuid.UUID) (*domains.User, error)
	GetUserByEmail(uuid.UUID) (*domains.User, error)
	UpdateUser(*domains.User) (*domains.User, error)
}


func (r *UserRepsotory) CreateUser(user *domains.User) (*domains.User, error){
	return nil, nil
} 


func (r *UserRepsotory) DeleteUserByID(id uuid.UUID) (error) {
	return nil
}

func (r *UserRepsotory) GetUserByID(id uuid.UUID) (*domains.User, error){
	return nil, nil
}

func (r *UserRepsotory) GetUserByUsername(id string) (*domains.User, error){
	return nil, nil
}

func (r *UserRepsotory) GetUserByEmail(id string) ( *domains.User, error){
	return nil, nil
}

func (r *UserRepsotory) UpdateUser(*domains.User) (*domains.User, error){
	return nil, nil
}











