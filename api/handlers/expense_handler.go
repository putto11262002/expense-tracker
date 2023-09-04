package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/repositories"
	"github.com/putto11262002/expense-tracker/api/services"
	"github.com/putto11262002/expense-tracker/api/utils"
	"net/http"
	"time"
)

type ExpenseHandler struct {
	expenseService services.IExpenseService
}

func NewExpenseHandler(expenseService services.IExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

type ExpenseResponse struct {
	GroupID     uuid.UUID       `json:"groupID"`
	ID          uuid.UUID       `json:"id"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Date        time.Time       `json:"date"`
	PaidBy      uuid.UUID       `json:"paidBy"`
	Amount      float64         `json:"amount"`
	Splits      []SplitResponse `json:"splits"`
	CreateAt    time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func NewExpenseResponse(groupID uuid.UUID, id uuid.UUID, description string, category string, date time.Time, paidBy uuid.UUID, amount float64, splits []SplitResponse, createdAt time.Time, updatedAt time.Time) *ExpenseResponse {
	return &ExpenseResponse{
		GroupID:     groupID,
		ID:          id,
		Description: description,
		Category:    category,
		Date:        date,
		PaidBy:      paidBy,
		Amount:      amount,
		Splits:      splits,
		CreateAt:    createdAt,
		UpdatedAt:   updatedAt,
	}
}

type SplitResponse struct {
	ExpenseID uuid.UUID `json:"expenseID"`
	Value     float64   `json:"value"`
	UserID    uuid.UUID `json:"userID"`
}

func NewSplitResponse(expenseID uuid.UUID, value float64, userID uuid.UUID) *SplitResponse {

	return &SplitResponse{
		ExpenseID: expenseID,
		Value:     value,
		UserID:    userID,
	}
}

type SplitRequest struct {
	UserID uuid.UUID `json:"userID" binding:"required"`
	Value  float64   `json:"value" binding:"required"`
}

type CreateExpenseRequest struct {
	GroupID     uuid.UUID `json:"groupID" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	PaidBy      uuid.UUID
	Amount      float64        `json:"amount" binding:"required"`
	Splits      []SplitRequest `json:"splits"  binding:"required"`
}

func (h *ExpenseHandler) HandleCreateExpense(ctx *gin.Context) {
	var expenseReq CreateExpenseRequest
	if err := ctx.ShouldBindJSON(&expenseReq); err != nil {
		utils.AbortWithError(ctx, fmt.Errorf("parsing create expense request: %w", err))
		return
	}

	var splitsInput []services.SplitInput
	for _, split := range expenseReq.Splits {
		splitsInput = append(splitsInput, *services.NewSplitInput(split.UserID, split.Value))
	}

	expenseID, err := h.expenseService.CreateExpense(services.NewExpenseInput(expenseReq.GroupID,
		expenseReq.Description,
		expenseReq.Category,
		expenseReq.Date,
		expenseReq.PaidBy,
		expenseReq.Amount,
		splitsInput))

	if err != nil {
		utils.AbortWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": expenseID})
}

func (h *ExpenseHandler) HandleGetExpenseByID(ctx *gin.Context) {
	IDStr := ctx.Param("id")
	ID, err := uuid.Parse(IDStr)
	if err != nil {
		utils.AbortWithError(ctx, &utils.InvalidArgumentError{Message: "invalid id"})
		return
	}

	expense, err := h.expenseService.GetExpenseByID(ID)
	if err != nil {
		utils.AbortWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, expense)

}

type GetExpenseQuery struct {
	GroupID uuid.UUID `json:"groupID"`
	UserID  uuid.UUID `json:"userID"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
}

func (h *ExpenseHandler) HandleGetExpense(ctx *gin.Context) {
	userIDStr := ctx.Query("userID")
	var userID uuid.UUID
	var err error

	if userIDStr != "" {
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			utils.AbortWithError(ctx, &utils.InvalidArgumentError{Message: "invalid user id"})
			return
		}

	}

	groupIDStr := ctx.Query("groupID")
	var groupID uuid.UUID
	if groupIDStr != "" {
		groupID, err = uuid.Parse(groupIDStr)
		if err != nil {
			utils.AbortWithError(ctx, &utils.InvalidArgumentError{Message: "invalid group id"})
			return
		}
	}

	fromStr := ctx.Query("from")
	var from time.Time
	if fromStr != "" {
		from, err = utils.ParseTimeFromISO8601(fromStr)
		if err != nil {
			utils.AbortWithError(ctx, &utils.InvalidArgumentError{Message: "invalid time format for from"})
			return
		}
	}

	toStr := ctx.Query("to")
	var to time.Time
	if toStr != "" {
		to, err = utils.ParseTimeFromISO8601(toStr)
		if err != nil {
			utils.AbortWithError(ctx, &utils.InvalidArgumentError{Message: "invalid time format for to"})
			return
		}
	}

	query := GetExpenseQuery{
		UserID:  userID,
		GroupID: groupID,
		From:    from,
		To:      to,
	}

	expenses, err := h.expenseService.GetExpenses(*repositories.NewGetExpenseFilter(
		[]uuid.UUID{query.GroupID},
		[]uuid.UUID{query.UserID},
		query.From,
		query.To,
	))
	if err != nil {
		utils.AbortWithError(ctx, err)
		return
	}

	expenseRespose := []ExpenseResponse{}
	for _, expense := range *expenses {
		splitResponse := []SplitResponse{}

		for _, split := range expense.Splits {
			splitResponse = append(splitResponse, *NewSplitResponse(
				split.ExpenseID,
				split.Value,
				split.UserID,
			))
		}
		expenseRespose = append(expenseRespose, *NewExpenseResponse(
			expense.GroupID,
			expense.ID,
			expense.Description,
			expense.Category,
			expense.Date,
			expense.PaidBy,
			expense.Amount,
			splitResponse,
			expense.CreateAt,
			expense.UpdatedAt,
		))
	}

	ctx.JSON(http.StatusOK, expenseRespose)

}
