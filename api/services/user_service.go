package services

import (
	"errors"
	"fmt"
	"github.com/putto11262002/expense-tracker/api/domains"
	"github.com/putto11262002/expense-tracker/api/repositories"
	"github.com/putto11262002/expense-tracker/api/utils"
	"gorm.io/gorm"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IUserRepository
	jwtSecret string
}

func NewUserService(repository repositories.IUserRepository, jwtSecret string) *UserService {
	return &UserService{
		repository: repository,
		jwtSecret: jwtSecret,
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
	ID        string
	Username  string
	Email     string
	FirstName string
	LastName  string
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

type GetUserFilter struct {
	Q string
	NotInGroupID uuid.UUID
	Email string
}

func NewGetUserFilter(q string, notInGroupID uuid.UUID, email string) *GetUserFilter {
	return &GetUserFilter{
		Q: q,
		NotInGroupID: notInGroupID,
		Email: email,
	}
}

// IUserService is an interface for managing user-related operations in the application.
type IUserService interface {
	// Register registers a new user with the provided registration input.
	Register(input *UserRegisterInput) (*domains.User, error)

	// Login performs user authentication with the provided login input.
	// It returns a result containing user information upon successful login.
	Login(input *UserLoginInput) (*UserLoginResult, error)

	// UpdateUser updates user information based on the provided input.
	UpdateUser(input *UserUpdateInput) (*domains.User, error)

	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(email string) (*domains.User, error)

	// GetUserByID retrieves a user by their unique identifier (ID).
	GetUserByID(userID uuid.UUID) (*domains.User, error)

	// GetUserByUsername retrieves a user by their username.
	GetUserByUsername(username string) (*domains.User, error)

	GetUsers(filter GetUserFilter) (*[]domains.User, error)
}

func (s *UserService) Register(input *UserRegisterInput) (*domains.User, error) {
	// check if user with this username or email exist

	exist, err := s.repository.ExistByUsernameOrEmail(input.Username, input.Email)
	if err != nil {
		return nil, fmt.Errorf("check for existing user before register: %w", err)
	}

	if exist {
		return nil, &utils.DataIntegrityError{
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

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &utils.AuthorizationError{
			Message: "invalid credentials",
		}
	}

	if err != nil {
		return nil, fmt.Errorf("checking if user exist before login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Secret))
	if err != nil {
		if errors.Is(err, bcrypt.ErrHashTooShort) || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, &utils.AuthorizationError{
				Message: "invalid credentials",
			}
		}
		return nil, fmt.Errorf("comparing password for login: %w", err)

	}

	token, maxAge, err := utils.GenerateJWTToken(user, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("signing user token: %w", err)
	}

	return NewLoginResult(user, token, maxAge), nil
}

func (s *UserService) UpdateUser(input *UserUpdateInput) (*domains.User, error) {
	return nil, nil
}

func (s *UserService) GetUserByEmail(email string) (*domains.User, error) {
	 user, err := s.repository.GetUserByEmail(email)
	 if err != nil {
		return nil, err
	 }
	 return user, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*domains.User, error) {
	user, err := s.repository.GetUserByID(id)

	if err != nil {
		
		return nil, fmt.Errorf("getting user by id: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*domains.User, error) {
	return nil, nil
}


func (s *UserService) GetUsers(filter GetUserFilter) (*[]domains.User, error) {
	users, err := s.repository.GetUsers(*repositories.NewUserFilter(filter.Email, filter.NotInGroupID, filter.Q))
	if err != nil {
		return nil, err
	}
	return users, nil
}
