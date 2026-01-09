
package repository

import (
	"time"

	"k-mate-backend/config"
	"k-mate-backend/internal/models"
)

// ExpenseRepository handles database operations for Expense model
type ExpenseRepository struct{}

// NewExpenseRepository creates a new ExpenseRepository instance
func NewExpenseRepository() *ExpenseRepository {
	return &ExpenseRepository{}
}

// Create inserts a new expense into the database
func (r *ExpenseRepository) Create(expense *models.Expense) error {
	return config.DB.Create(expense).Error
}

// FindAll retrieves all expenses with user info
func (r *ExpenseRepository) FindAll() ([]models.Expense, error) {
	var expenses []models.Expense
	err := config.DB.Preload("User").
		Order("expense_date DESC, id DESC").
		Find(&expenses).Error
	return expenses, err
}

// FindAllPaginated retrieves expenses with pagination
func (r *ExpenseRepository) FindAllPaginated(page, limit int) ([]models.Expense, int64, error) {
	var expenses []models.Expense
	var total int64

	if err := config.DB.Model(&models.Expense{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := config.DB.Preload("User").
		Offset(offset).
		Limit(limit).
		Order("expense_date DESC, id DESC").
		Find(&expenses).Error

	return expenses, total, err
}

// FindByID retrieves an expense by ID
func (r *ExpenseRepository) FindByID(id int) (*models.Expense, error) {
	var expense models.Expense
	err := config.DB.Preload("User").First(&expense, id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

// FindByDateRange retrieves expenses within a date range
func (r *ExpenseRepository) FindByDateRange(startDate, endDate time.Time) ([]models.Expense, error) {
	var expenses []models.Expense
	err := config.DB.Preload("User").
		Where("expense_date BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("expense_date DESC").
		Find(&expenses).Error
	return expenses, err
}

// FindByCategory retrieves expenses by category
func (r *ExpenseRepository) FindByCategory(category string) ([]models.Expense, error) {
	var expenses []models.Expense
	err := config.DB.Preload("User").
		Where("category = ?", category).
		Order("expense_date DESC").
		Find(&expenses).Error
	return expenses, err
}

// GetTodayExpenses retrieves today's expenses
func (r *ExpenseRepository) GetTodayExpenses() ([]models.Expense, error) {
	today := time.Now().Format("2006-01-02")
	var expenses []models.Expense
	err := config.DB.Preload("User").
		Where("expense_date = ?", today).
		Find(&expenses).Error
	return expenses, err
}

// Update updates an expense
func (r *ExpenseRepository) Update(expense *models.Expense) error {
	return config.DB.Save(expense).Error
}

// Delete deletes an expense
func (r *ExpenseRepository) Delete(id int) error {
	return config.DB.Delete(&models.Expense{}, id).Error
}

// CountAll returns total number of expenses
func (r *ExpenseRepository) CountAll() (int64, error) {
	var count int64
	err := config.DB.Model(&models.Expense{}).Count(&count).Error
	return count, err
}
