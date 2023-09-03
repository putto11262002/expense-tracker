package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	SplitMode   string         `json:"splitMode" binding:"required"`
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
	group  uuid.UUID
	userID uuid.UUID
	from   time.Time
	to     time.Time
}

func (h *ExpenseHandler) HandleGetExpense() {

}
