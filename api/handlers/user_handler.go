package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/configs"
	"github.com/putto11262002/expense-tracker/api/domains"
	"github.com/putto11262002/expense-tracker/api/services"
	"github.com/putto11262002/expense-tracker/api/utils"
	"log"
	"math"
	"net/http"
	"time"
)

type UserHandler struct {
	service services.IUserService
}

func NewUserHandler(service services.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type RegisterRequestBody struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserLoginRequest struct {
	Key    string `json:"key" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

type UpdateUserRequestBody struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func NewLoginResponse(result *services.UserLoginResult) *LoginResponse {
	return &LoginResponse{
		Token: result.Token,
		User:  NewUserResponse(result.User),
	}
}

func NewUserResponse(user *domains.User) *UserResponse {
	if user == nil {
		return &UserResponse{}
	}

	return &UserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (h *UserHandler) HandleRegister(ctx *gin.Context) {
	var userReqBody RegisterRequestBody

	if err := ctx.ShouldBindJSON(&userReqBody); err != nil {
		utils.AbortWithError(ctx, fmt.Errorf("parsing register request body: %w", err))
		return
	}

	user, err := h.service.Register(services.NewUserRegisterInput(
		userReqBody.FirstName,
		userReqBody.LastName,
		userReqBody.Username,
		userReqBody.Email,
		userReqBody.Password,
	))

	if err != nil {
		utils.AbortWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, NewUserResponse(user))
}

func (h *UserHandler) HandleLogin(ctx *gin.Context) {
	var loginReqBody UserLoginRequest
	if err := ctx.ShouldBindJSON(&loginReqBody); err != nil {
		utils.AbortWithError(ctx, fmt.Errorf("parsing login request body: %w", err))
		return
	}

	result, err := h.service.Login(services.NewUserLoginInput(
		loginReqBody.Key,
		loginReqBody.Secret,
	))

	if err != nil {
		utils.AbortWithError(ctx, err)
		return
	}

	domain, err := configs.GetStringEnv("DOMAIN")
	if err != nil {
		log.Println("DOMAIN is not configured")
	}

	var secure bool
	if configs.GetGoEnv() == "production" {
		secure = true
	} else {
		secure = false
	}

	ctx.SetCookie("token", result.Token, int(math.Ceil(result.MaxAge.Seconds())), "/", domain, secure, true)

	ctx.JSON(200, NewLoginResponse(result))
}

func (h *UserHandler) HandleGetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.AbortWithError(ctx, &utils.InvalidArgumentError{
			Message: "invalid id",
		})
		return
	}
	user, err := h.service.GetUserByID(id)
	if err != nil {
		utils.AbortWithError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, NewUserResponse(user))
}
