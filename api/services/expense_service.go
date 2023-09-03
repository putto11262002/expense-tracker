package services

import (
	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/domains"
	"github.com/putto11262002/expense-tracker/api/repositories"
	"github.com/putto11262002/expense-tracker/api/utils"
	"time"
)

type IExpenseService interface {
	// CreateExpense validates that Split.Value add up to Amount and create a new expense
	CreateExpense(input *CreateExpenseInput) (uuid.UUID, error)
	GetExpenseByID(id uuid.UUID) (*domains.Expense, error)
	GetExpenses(filter GetExpenseFilter) (*[]domains.Expense, error)
}

type ExpenseService struct {
	expenseRepository repositories.IExpenseRepository
}

func NewExpenseService(expenseRepository repositories.IExpenseRepository) *ExpenseService {
	return &ExpenseService{
		expenseRepository: expenseRepository,
	}

}

type SplitInput struct {
	UserID uuid.UUID
	Value  float64
}

type CreateExpenseInput struct {
	GroupID     uuid.UUID
	Description string
	Category    string
	Date        time.Time
	PaidBy      uuid.UUID
	Amount      float64
	Splits      []SplitInput
}

func NewSplitInput(userID uuid.UUID, value float64) *SplitInput {
	return &SplitInput{
		UserID: userID,
		Value:  value,
	}
}

func NewExpenseInput(groupID uuid.UUID, description string, category string, date time.Time, paidBy uuid.UUID, amount float64, splits []SplitInput) *CreateExpenseInput {
	return &CreateExpenseInput{
		GroupID:     groupID,
		Description: description,
		Category:    category,
		Date:        date,
		PaidBy:      paidBy,
		Amount:      amount,
		Splits:      splits,
	}
}

type GetExpenseFilter struct {
}

func (e ExpenseService) CreateExpense(input *CreateExpenseInput) (uuid.UUID, error) {
	// validate that the all Split.Value in input.Splits add up to input.Amount
	sum := int64(0)

	for _, split := range input.Splits {
		sum += utils.FloatToIntCurrency(split.Value)
	}

	if (sum - utils.FloatToIntCurrency(input.Amount)) != 0 {
		return uuid.Nil, &utils.InvalidArgumentError{
			Message: "invalid split values",
		}
	}

	var split []domains.Split
	for _, splitInput := range input.Splits {
		split = append(split, *domains.NewSplits(splitInput.Value, splitInput.UserID))
	}

	expenseID, err := e.expenseRepository.CreateExpense(
		domains.NewExpense(
			input.GroupID,
			input.Description,
			input.Category,
			input.Date,
			input.PaidBy,
			input.Amount,
			split,
		))
	if err != nil {
		return uuid.Nil, err
	}

	return expenseID, nil
}

func (e ExpenseService) GetExpenseByID(id uuid.UUID) (*domains.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e ExpenseService) GetExpenses(filter GetExpenseFilter) (*[]domains.Expense, error) {
	//TODO implement me
	panic("implement me")
}
