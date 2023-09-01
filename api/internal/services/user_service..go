package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/internal/domains"
	"github.com/putto11262002/expense-tracker/api/internal/repositories"
	"github.com/putto11262002/expense-tracker/api/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IUserRepository
}

func NewUserService(repository repositories.IUserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

type userRegisterInput struct {
	Username string
	Email string
	Password string
	FirstName string
	LastName string
}


func NewUserRegisterInput(firstName, lastName, username, email, password string) *userRegisterInput{
	return &userRegisterInput{
		FirstName: firstName,
		LastName: lastName,
		Username: username,
		Email: email,
		Password: password,
	}
}



type UserLoginInput struct {
}

type UserUpdateInput struct {
}

type UserLoginResult struct {
}

type IUserService interface {
	Register(*userRegisterInput) (*domains.User, error)
	Login(*UserLoginInput) (*UserLoginResult, error)
	UpdateUser(*UserUpdateInput) (*domains.User, error)
	GetUserByEmail(string) (*domains.User, error)
	GetUserByID(uuid.UUID) (*domains.User, error)
	GetUserByUsername(string) (*domains.User, error)
}

func (s *UserService) Register(input *userRegisterInput) (*domains.User, error) {
	// check if user with this username or email exist



	emailExistRes := make(chan bool)
	emailExistErr := make(chan error)
	usernameExistRes := make(chan bool)
	usernameExistErr := make(chan error)

	go func ()  {
		exist, err := s.repository.ExistByEmail(input.Email)
		if err != nil {
			emailExistRes <- false
			emailExistErr <- err
			return
		}
		emailExistRes <- exist
		emailExistErr <- nil
	}()

	go func ()  {
		exist, err := s.repository.ExistByUsername(input.Username)
		if err != nil {
			usernameExistRes <- false
			usernameExistErr <- err
			return
		}
		usernameExistRes <- exist
		usernameExistErr <- nil
	}()


	if exist, err := <-emailExistRes, <- emailExistErr; err != nil {
		return nil , err
	}else if exist {
		return nil, &utils.DataIntegrityError{
			Message: "user with this email already exist",
		}
	}

	if exist, err := <- usernameExistRes, <- usernameExistErr; err != nil {
		return nil , err
	}else if exist {
		return nil, &utils.DataIntegrityError{
			Message: "user with this username already exist",
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing user password: %w", err)
	}


	user, err := s.repository.CreateUser(domains.NewUser(
		input.FirstName, 
		input.LastName, 
		input.Username, 
		input.Email, 
		string(hashed)))

	if err != nil {
		return nil , err
	}
	
	return user, nil

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


