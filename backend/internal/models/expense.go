package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// ExpenseCategory represents the allowed expense categories
type ExpenseCategory string

const (
	ExpenseCategoryOperational ExpenseCategory = "operational"
	ExpenseCategorySalary      ExpenseCategory = "salary"
	ExpenseCategoryRestock     ExpenseCategory = "restock"
	ExpenseCategoryMaintenance ExpenseCategory = "maintenance"
	ExpenseCategoryRefund      ExpenseCategory = "refund"
	ExpenseCategoryOther       ExpenseCategory = "other"
)

// Expense represents the expenses table (Operational Cost)
type Expense struct {
	ID          int             `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int             `gorm:"column:user_id;not null;index" json:"user_id"`
	Title       string          `gorm:"column:title;type:varchar(100);not null" json:"title"`
	Category    ExpenseCategory `gorm:"column:category;type:enum('operational','salary','restock','maintenance','refund','other');not null" json:"category"`
	Amount      decimal.Decimal `gorm:"column:amount;type:decimal(15,2);not null" json:"amount"`
	ExpenseDate time.Time       `gorm:"column:expense_date;type:date;not null" json:"expense_date"`
	Description *string         `gorm:"column:description;type:text" json:"description"`
	CreatedAt   time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Expense model
func (Expense) TableName() string {
	return "expenses"
}
