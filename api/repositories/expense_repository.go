package repositories

import (
	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/domains"
	"gorm.io/gorm"
	"time"
)

type IExpenseRepository interface {
	CreateExpense(expense *domains.Expense) (uuid.UUID, error)
	GetExpenseByID(id uuid.UUID) (*domains.Expense, error)
	GetExpenses(filter GetExpenseFilter) (*[]domains.Expense, error)
}

type GetExpenseFilter struct {
	From     time.Time
	To       time.Time
	GroupIDs []uuid.UUID
	userIDs  []uuid.UUID
}

type ExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{
		db: db,
	}
}

func (e ExpenseRepository) CreateExpense(expense *domains.Expense) (uuid.UUID, error) {
	if err := e.db.Create(expense).Error; err != nil {
		return uuid.Nil, err
	}
	return expense.ID, nil
}

func (e ExpenseRepository) GetExpenseByID(id uuid.UUID) (*domains.Expense, error) {
	var expense domains.Expense
	if err := e.db.Preload("Splits").First(&expense, id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (e ExpenseRepository) GetExpenses(filter GetExpenseFilter) (*[]domains.Expense, error) {
	//TODO implement me
	panic("implement me")
}
