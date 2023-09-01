package repositories

import (
	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/internal/domains"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: DB,
	}
}

type IUserRepository interface {
	CreateUser(*domains.User) (*domains.User, error)
	DeleteUserByID(uuid.UUID) error
	GetUserByID(uuid.UUID) (*domains.User, error)
	GetUserByUsername(string) (*domains.User, error)
	GetUserByEmail(string) (*domains.User, error)
	UpdateUser(*domains.User) (*domains.User, error)
	ExistByUsername(string)(bool, error)
	ExistByEmail(string)(bool, error)
}

func (r *UserRepository) CreateUser(user *domains.User) (*domains.User, error) {
	result := r.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *UserRepository) DeleteUserByID(id uuid.UUID) error {
	return nil
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*domains.User, error) {
	return nil, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*domains.User, error) {
	return nil, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*domains.User, error) {
	return nil, nil
}

func (r *UserRepository) UpdateUser(user *domains.User) (*domains.User, error) {
	return nil, nil
}

func (r *UserRepository) ExistByUsername(username string)(bool, error){
	var count int64
	if err := r.DB.Model(&domains.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r * UserRepository) ExistByEmail(email string)(bool, error){
	var count int64
	if err := r.DB.Model(&domains.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}