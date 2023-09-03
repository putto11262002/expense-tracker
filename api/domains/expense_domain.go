package domains

import (
	"github.com/google/uuid"
	"time"
)

type Expense struct {
	GroupID     uuid.UUID `gorm:"type:char(36)"`
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	Description string    // optional: description of the expense
	Category    string    // the category of expense, it will be a predefined options
	Date        time.Time // required: the date in which this expense was incurred
	PaidBy      uuid.UUID // required: the user id that paid for the expense
	Amount      float64   // required: the amount of the expense
	Splits      []Split   `gorm:"foreignKey:ExpenseID"` // required, Split.Value must add up to Amount; a splice of splits associated to this expense
	CreateAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type Split struct {
	ExpenseID uuid.UUID `gorm:"type:char(36)"` // the ID of the expense that the split is associated with
	Value     float64
	UserID    uuid.UUID // the id of the group member associated with the split
}

func NewSplits(
	value float64,
	userID uuid.UUID,
) *Split {
	return &Split{
		UserID: userID,
		Value:  value,
	}
}
func NewExpense(
	groupID uuid.UUID,
	description string,
	category string,
	date time.Time,
	paidBy uuid.UUID,
	amount float64,
	splits []Split,
) *Expense {
	return &Expense{
		ID:          uuid.New(),
		GroupID:     groupID,
		Description: description,
		Category:    category,
		Date:        date,
		PaidBy:      paidBy,
		Amount:      amount,
		Splits:      splits,
	}
}
