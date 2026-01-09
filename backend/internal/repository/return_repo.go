package repository

import (
	"fmt"
	"time"

	"k-mate-backend/config"
	"k-mate-backend/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ReturnRepository handles database operations for Return model
type ReturnRepository struct{}

// NewReturnRepository creates a new ReturnRepository instance
func NewReturnRepository() *ReturnRepository {
	return &ReturnRepository{}
}

// ReturnItem represents a prepared return item with all necessary data
type ReturnItem struct {
	TransactionDetailID int
	ProductUnitID       int
	ProductID           int
	ProductName         string
	UnitName            string
	Qty                 int
	RefundPricePerUnit  decimal.Decimal
	SubtotalRefund      decimal.Decimal
	Condition           string // "good" or "damaged"
	ConversionFactor    int
	StockToRestore      int // In base unit (only if condition = good)
}

// CreateReturn creates a return record with details and updates stock in a single transaction
func (r *ReturnRepository) CreateReturn(
	returnHeader *models.Return,
	items []ReturnItem,
) error {
	// Start database transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Defer rollback in case of panic
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 1. Generate return number
	returnHeader.ReturnNo = generateReturnNo()

	// 2. Calculate total refund
	totalRefund := decimal.Zero
	for _, item := range items {
		totalRefund = totalRefund.Add(item.SubtotalRefund)
	}
	returnHeader.TotalRefund = totalRefund

	// 3. Create return header
	if err := tx.Create(returnHeader).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menyimpan data return: %v", err)
	}

	// 4. Create return details and update stock
	for _, item := range items {
		// Create detail
		detail := models.ReturnDetail{
			ReturnID:           returnHeader.ID,
			ProductUnitID:      item.ProductUnitID,
			Qty:                item.Qty,
			RefundPricePerUnit: item.RefundPricePerUnit,
			ConditionStatus:    models.ConditionStatus(item.Condition),
			SubtotalRefund:     item.SubtotalRefund,
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal menyimpan detail return: %v", err)
		}

		// Update stock only if condition is "good" (layak jual)
		if item.Condition == "good" && item.StockToRestore > 0 {
			result := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock_qty", gorm.Expr("stock_qty + ?", item.StockToRestore))

			if result.Error != nil {
				tx.Rollback()
				return fmt.Errorf("gagal update stok: %v", result.Error)
			}
		}
	}

	// 5. Create expense entry for the refund
	refundDescription := fmt.Sprintf("Refund untuk %s - %d item diretur", returnHeader.ReturnNo, len(items))
	expense := models.Expense{
		UserID:      returnHeader.UserID,
		Title:       fmt.Sprintf("Refund %s", returnHeader.ReturnNo),
		Category:    models.ExpenseCategoryRefund,
		Amount:      totalRefund,
		ExpenseDate: time.Now(),
		Description: &refundDescription,
	}

	if err := tx.Create(&expense).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menyimpan expense refund: %v", err)
	}

	// 6. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi return: %v", err)
	}

	return nil
}

// FindAll retrieves all returns with details
func (r *ReturnRepository) FindAll() ([]models.Return, error) {
	var returns []models.Return
	err := config.DB.Preload("User").
		Preload("Transaction").
		Preload("ReturnDetails").
		Preload("ReturnDetails.ProductUnit").
		Preload("ReturnDetails.ProductUnit.Product").
		Order("id DESC").
		Find(&returns).Error
	return returns, err
}

// FindByID retrieves a return by ID with details
func (r *ReturnRepository) FindByID(id int) (*models.Return, error) {
	var returnRecord models.Return
	err := config.DB.Preload("User").
		Preload("Transaction").
		Preload("ReturnDetails").
		Preload("ReturnDetails.ProductUnit").
		Preload("ReturnDetails.ProductUnit.Product").
		First(&returnRecord, id).Error
	if err != nil {
		return nil, err
	}
	return &returnRecord, nil
}

// FindByTransactionID retrieves returns for a specific transaction
func (r *ReturnRepository) FindByTransactionID(transactionID int) ([]models.Return, error) {
	var returns []models.Return
	err := config.DB.Preload("ReturnDetails").
		Where("transaction_id = ?", transactionID).
		Order("id DESC").
		Find(&returns).Error
	return returns, err
}

// generateReturnNo generates a unique return number
// Format: RTN-YYYYMMDD-XXXX
func generateReturnNo() string {
	now := time.Now()
	dateStr := now.Format("20060102")

	// Get count of returns today for sequence
	var count int64
	config.DB.Model(&models.Return{}).
		Where("DATE(return_date) = ?", now.Format("2006-01-02")).
		Count(&count)

	return fmt.Sprintf("RTN-%s-%04d", dateStr, count+1)
}
