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
	GroupIDs []uuid.UUID
	UserIDs  []uuid.UUID
	From     time.Time
	To       time.Time
}

func NewGetExpenseFilter(groupIDs []uuid.UUID, userIDs []uuid.UUID, from time.Time, to time.Time) *GetExpenseFilter {
	return &GetExpenseFilter{
		GroupIDs: groupIDs,
		UserIDs:  userIDs,
		From:     from,
		To:       to,
	}
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
	query := e.db

	var expenses []domains.Expense

	var subquery *gorm.DB

	if filter.UserIDs != nil {
		subquery = e.db.Table("splits").Where("user_id in ?", filter.UserIDs)

	}

	if filter.GroupIDs != nil {
		query.Where("group_id in ?", filter.GroupIDs)
	}

	if !filter.To.IsZero() {
		query.Where("date <= ?", filter.To)
	}

	if !filter.From.IsZero() {
		query.Where("date >= ?", filter.From)
	}

	if subquery != nil {
		query.Where("id in (?)", subquery.Select("group_id"))
	}

	if err := query.Preload("Splits").Find(&expenses).Error; err != nil {
		return nil, err
	}

	return &expenses, nil
}
