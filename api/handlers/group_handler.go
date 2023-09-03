package handlers

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/domains"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/services"
	"github.com/putto11262002/expense-tracker/api/utils"
)

type GroupHandler struct {
	groupService services.IGroupService
	userService  services.IUserService
}

func NewGroupHandler(groupService services.IGroupService, userService services.IUserService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
		userService:  userService,
	}
}

type CreateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type GroupResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Members   []UserResponse `json:"members"`
	Owner     UserResponse   `json:"owner"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func NewGroupResponse(group *domains.Group) *GroupResponse {
	if group == nil {
		return &GroupResponse{}
	}

	members := []UserResponse{}
	for _, member := range group.Members {
		members = append(members, *NewUserResponse(&member))
	}
	owner := *NewUserResponse(&group.Owner)

	return &GroupResponse{
		Name:      group.Name,
		Owner:     owner,
		Members:   members,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}

}

func (h *GroupHandler) HandleCreateGroup(ctx *gin.Context) {
	var createGroupReq CreateGroupRequest
	if err := ctx.ShouldBindJSON(&createGroupReq); err != nil {
		utils.AbortWithError(ctx, fmt.Errorf("parsing create group request: %w", err))
		return
	}

	claims := utils.GetClaimsFromCtx(ctx)
	if claims == nil {
		utils.AbortWithError(ctx, errors.New("cannot retrieve claims from context"))
		return
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		utils.AbortWithError(ctx, fmt.Errorf("parse user id from claims: %w", err))
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return
	}
	if user == nil {
		utils.AbortWithError(ctx, errors.New("cannot find user using claims subject"))
	}

	group, err := h.groupService.CreateGroup(user, createGroupReq.Name)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, NewGroupResponse(group))

}
