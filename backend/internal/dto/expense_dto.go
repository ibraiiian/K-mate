package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// ================== REQUEST DTOs ==================

// CreateExpenseRequest represents the JSON body for creating an expense
type CreateExpenseRequest struct {
	Title       string          `json:"title" validate:"required"`
	Category    string          `json:"category" validate:"required,oneof=operational salary restock maintenance other"`
	Amount      decimal.Decimal `json:"amount" validate:"required"`
	ExpenseDate string          `json:"expense_date" validate:"required"` // Format: YYYY-MM-DD
	Description *string         `json:"description"`
}

// UpdateExpenseRequest represents the JSON body for updating an expense
type UpdateExpenseRequest struct {
	Title       string          `json:"title"`
	Category    string          `json:"category"`
	Amount      decimal.Decimal `json:"amount"`
	ExpenseDate string          `json:"expense_date"`
	Description *string         `json:"description"`
}

// ================== RESPONSE DTOs ==================

// ExpenseResponse represents expense data in response
type ExpenseResponse struct {
	ID          int             `json:"id"`
	UserID      int             `json:"user_id"`
	UserName    string          `json:"user_name"`
	Title       string          `json:"title"`
	Category    string          `json:"category"`
	Amount      decimal.Decimal `json:"amount"`
	ExpenseDate string          `json:"expense_date"`
	Description *string         `json:"description"`
	CreatedAt   string          `json:"created_at"`
}

// ExpenseListResponse represents expense list
type ExpenseListResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Total    int64             `json:"total"`
}

// ParseExpenseDate parses expense date string to time.Time
func ParseExpenseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
