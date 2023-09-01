package services

import (
	"errors"
	"fmt"
	"github.com/putto11262002/expense-tracker/api/domains"
	"github.com/putto11262002/expense-tracker/api/repositories"
	utils2 "github.com/putto11262002/expense-tracker/api/utils"
	"time"

	"github.com/google/uuid"
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

type UserRegisterInput struct {
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func NewUserRegisterInput(firstName, lastName, username, email, password string) *UserRegisterInput {
	return &UserRegisterInput{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
		Password:  password,
	}
}

type UserLoginInput struct {
	Key    string
	Secret string
}

func NewUserLoginInput(key, secret string) *UserLoginInput {
	return &UserLoginInput{
		Key:    key,
		Secret: secret,
	}
}

type UserUpdateInput struct {
}

type UserLoginResult struct {
	User   *domains.User
	Token  string
	MaxAge time.Duration
}

func NewLoginResult(user *domains.User, token string, maxAge time.Duration) *UserLoginResult {
	return &UserLoginResult{
		User:   user,
		Token:  token,
		MaxAge: maxAge,
	}
}

type IUserService interface {
	Register(*UserRegisterInput) (*domains.User, error)
	Login(*UserLoginInput) (*UserLoginResult, error)
	UpdateUser(*UserUpdateInput) (*domains.User, error)
	GetUserByEmail(string) (*domains.User, error)
	GetUserByID(uuid.UUID) (*domains.User, error)
	GetUserByUsername(string) (*domains.User, error)
}

func (s *UserService) Register(input *UserRegisterInput) (*domains.User, error) {
	// check if user with this username or email exist

	exist, err := s.repository.ExistByUsernameOrEmail(input.Username, input.Email)
	if err != nil {
		return nil, fmt.Errorf("check for existing user before register: %w", err)
	}

	if exist {
		return nil, &utils2.DataIntegrityError{
			Message: "user already exist",
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
		return nil, err
	}

	return user, nil

}

func (s *UserService) Login(input *UserLoginInput) (*UserLoginResult, error) {
	user, err := s.repository.GetUserByUsernameOrEmail(input.Key, input.Key)
	if err != nil {
		return nil, fmt.Errorf("checking if user exist before login: %w", err)
	}

	if user == nil {
		return nil, &utils2.AuthorizationError{
			Message: "invalid credentials",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Secret))
	if err != nil {
		if errors.Is(err, bcrypt.ErrHashTooShort) || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, &utils2.AuthorizationError{
				Message: "invalid credentials",
			}
		}
		return nil, fmt.Errorf("comparing password for login: %w", err)

	}

	token, maxAge, err := utils2.GenerateJWTToken(user, utils2.GetJWTSecret())
	if err != nil {
		return nil, fmt.Errorf("signing user token: %w", err)
	}

	return NewLoginResult(user, token, maxAge), nil
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
