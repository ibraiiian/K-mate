package services

import (
	"errors"

	"k-mate-backend/internal/dto"
	"k-mate-backend/internal/models"
	"k-mate-backend/internal/repository"

	"github.com/shopspring/decimal"
)

// ExpenseService handles expense business logic
type ExpenseService struct {
	expenseRepo *repository.ExpenseRepository
}

// NewExpenseService creates a new ExpenseService instance
func NewExpenseService(expenseRepo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{
		expenseRepo: expenseRepo,
	}
}

// CreateExpense creates a new expense record
func (s *ExpenseService) CreateExpense(userID int, req *dto.CreateExpenseRequest) (*models.Expense, error) {
	// Validate amount
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("amount harus lebih dari 0")
	}

	// Validate category
	validCategories := map[string]bool{
		"operational": true, "salary": true, "restock": true, "maintenance": true, "other": true,
	}
	if !validCategories[req.Category] {
		return nil, errors.New("kategori tidak valid (operational/salary/restock/maintenance/other)")
	}

	// Parse expense date
	expenseDate, err := dto.ParseExpenseDate(req.ExpenseDate)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid (gunakan YYYY-MM-DD)")
	}

	expense := &models.Expense{
		UserID:      userID,
		Title:       req.Title,
		Category:    models.ExpenseCategory(req.Category),
		Amount:      req.Amount,
		ExpenseDate: expenseDate,
		Description: req.Description,
	}

	if err := s.expenseRepo.Create(expense); err != nil {
		return nil, errors.New("gagal menyimpan pengeluaran")
	}

	return expense, nil
}

// GetAllExpenses retrieves all expenses
func (s *ExpenseService) GetAllExpenses() ([]models.Expense, error) {
	return s.expenseRepo.FindAll()
}

// GetExpensesPaginated retrieves expenses with pagination
func (s *ExpenseService) GetExpensesPaginated(page, limit int) ([]models.Expense, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.expenseRepo.FindAllPaginated(page, limit)
}

// GetExpenseByID retrieves an expense by ID
func (s *ExpenseService) GetExpenseByID(id int) (*models.Expense, error) {
	expense, err := s.expenseRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("pengeluaran tidak ditemukan")
	}
	return expense, nil
}

// GetExpensesByCategory retrieves expenses by category
func (s *ExpenseService) GetExpensesByCategory(category string) ([]models.Expense, error) {
	return s.expenseRepo.FindByCategory(category)
}

// GetTodayExpenses retrieves today's expenses
func (s *ExpenseService) GetTodayExpenses() ([]models.Expense, error) {
	return s.expenseRepo.GetTodayExpenses()
}

// UpdateExpense updates an expense
func (s *ExpenseService) UpdateExpense(id int, req *dto.UpdateExpenseRequest) (*models.Expense, error) {
	expense, err := s.expenseRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("pengeluaran tidak ditemukan")
	}

	if req.Title != "" {
		expense.Title = req.Title
	}
	if req.Category != "" {
		expense.Category = models.ExpenseCategory(req.Category)
	}
	if !req.Amount.IsZero() && req.Amount.GreaterThan(decimal.Zero) {
		expense.Amount = req.Amount
	}
	if req.ExpenseDate != "" {
		expenseDate, err := dto.ParseExpenseDate(req.ExpenseDate)
		if err == nil {
			expense.ExpenseDate = expenseDate
		}
	}
	if req.Description != nil {
		expense.Description = req.Description
	}

	if err := s.expenseRepo.Update(expense); err != nil {
		return nil, errors.New("gagal mengupdate pengeluaran")
	}

	return expense, nil
}

// DeleteExpense deletes an expense
func (s *ExpenseService) DeleteExpense(id int) error {
	_, err := s.expenseRepo.FindByID(id)
	if err != nil {
		return errors.New("pengeluaran tidak ditemukan")
	}

	return s.expenseRepo.Delete(id)
}
