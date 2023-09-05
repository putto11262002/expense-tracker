package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/domains"

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
		return nil
	}

	var members []UserResponse
	for _, member := range group.Members {
		members = append(members, *UserDomainToResponse(&member))
	}
	owner := *UserDomainToResponse(&group.Owner)

	return &GroupResponse{
		ID:        group.ID.String(),
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

func (h *GroupHandler) HandleGetGroupByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.AbortWithError(ctx, &utils.InvalidArgumentError{
			Message: "invalid group id",
		})
		return
	}
	group, err := h.groupService.GetGroupByID(id)
	if err != nil {
		utils.AbortWithError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, NewGroupResponse(group))
}

func (h *GroupHandler) HandleGetGroupsByUserID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.AbortWithError(ctx, &utils.InvalidArgumentError{
			Message: "invalid group id",
		})
		return
	}
	groups, err := h.groupService.GetGroupsByUserID(id)
	if err != nil {
		return
	}

	var groupResponses []*GroupResponse

	for _, group := range *groups {
		groupResponse := NewGroupResponse(&group)
		groupResponses = append(groupResponses, groupResponse)
	}

	ctx.JSON(http.StatusOK, groupResponses)
}

type AddUserToGroupRequest struct {
	GroupID uuid.UUID `json:"groupID"`
	UserID  uuid.UUID `json:"userID"`
}

func (h *GroupHandler) HandleAddGroupMember(ctx *gin.Context) {
	var request AddUserToGroupRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.AbortWithError(ctx, err)
		return
	}

	user, err := h.userService.GetUserByID(request.UserID)
	if err != nil {
		utils.AbortWithError(ctx, err)
		return
	}

	if user == nil {
		utils.AbortWithError(ctx, &utils.InvalidArgumentError{Message: "invalid user id"})
		return
	}

	err = h.groupService.AddMember(request.GroupID, user)
	if err != nil {
		utils.AbortWithError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)

}
