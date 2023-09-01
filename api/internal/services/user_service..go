package services

import (
	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/internal/domains"
	"github.com/putto11262002/expense-tracker/api/internal/repositories"
)




type UserService struct {
	 repositories *repositories.IUserRepsotory
}


func NewUserService(repositories *repositories.IUserRepsotory) *UserService {
	return &UserService{
		repositories: repositories,
	}
}

type UserSignUpInput struct {

}

type UserLoginInput struct {

}

type UserUpdateInput struct {

}


type UserLoginResult struct {

}


type IUserService interface {
	Register(*UserSignUpInput) (*domains.User, error)
	Login(*UserLoginInput) (*UserLoginResult, error)
	UpdateUser(*UserUpdateInput) (*domains.User, error)
	GetUserByEmail(string) (*domains.User, error)
	GetUserByID(uuid.UUID) (*domains.User, error)
	GetUserByUsername(string)(*domains.User, error)
}


func (s *UserService) Register(input *UserSignUpInput) (*domains.User, error) {
	return nil, nil

}

func (s *UserService) Login(input *UserLoginInput) (*UserLoginResult, error) {
	return nil, nil 
}

func (s *UserService) UpdateUser(input *UserUpdateInput) (*domains.User, error) {
	return nil, nil
}

func (s *UserService) GetUserByEmail(email string) (*domains.User, error) {
	return nil, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*domains.User, error) {
	return nil, nil
}

func (s *UserService) GetUserByUsername(username string) (*domains.User, error) {
	return nil, nil
}





