package repository

import (
	"fmt"
	"time"

	"k-mate-backend/config"
	"k-mate-backend/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TransactionRepository handles database operations for Transaction model
type TransactionRepository struct{}

// NewTransactionRepository creates a new TransactionRepository instance
func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

// TransactionItem represents a prepared item for checkout with all necessary data
type TransactionItem struct {
	ProductID        int
	ProductName      string
	ProductUnitID    int
	UnitName         string
	Qty              int
	ConversionFactor int
	SellPrice        decimal.Decimal // Price at moment
	BuyPrice         decimal.Decimal // Cost at moment (HPP)
	Subtotal         decimal.Decimal
	StockToDeduct    int // In base unit
}

// CreateTransaction creates a transaction with details and updates stock in a single transaction
func (r *TransactionRepository) CreateTransaction(
	transaction *models.Transaction,
	items []TransactionItem,
) error {
	// Start database transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Defer rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Generate invoice number
	transaction.InvoiceNo = generateInvoiceNo()

	// 2. Calculate total amount
	totalAmount := decimal.Zero
	for _, item := range items {
		totalAmount = totalAmount.Add(item.Subtotal)
	}
	transaction.TotalAmount = totalAmount

	// 3. Create transaction header
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menyimpan transaksi: %v", err)
	}

	// 4. Create transaction details and update stock
	for _, item := range items {
		// Create detail
		detail := models.TransactionDetail{
			TransactionID: transaction.ID,
			ProductUnitID: item.ProductUnitID,
			Qty:           item.Qty,
			PriceAtMoment: item.SellPrice,
			CostAtMoment:  item.BuyPrice,
			Subtotal:      item.Subtotal,
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal menyimpan detail transaksi: %v", err)
		}

		// Update stock (deduct in base unit)
		result := tx.Model(&models.Product{}).
			Where("id = ? AND stock_qty >= ?", item.ProductID, item.StockToDeduct).
			Update("stock_qty", gorm.Expr("stock_qty - ?", item.StockToDeduct))

		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update stok: %v", result.Error)
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("stok tidak cukup untuk produk %s", item.ProductName)
		}
	}

	// 5. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi: %v", err)
	}

	return nil
}

// FindAll retrieves all transactions with user info
func (r *TransactionRepository) FindAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := config.DB.Preload("User").
		Order("id DESC").
		Find(&transactions).Error
	return transactions, err
}

// FindAllPaginated retrieves transactions with pagination
func (r *TransactionRepository) FindAllPaginated(page, limit int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	if err := config.DB.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := config.DB.Preload("User").
		Offset(offset).
		Limit(limit).
		Order("id DESC").
		Find(&transactions).Error

	return transactions, total, err
}

// FindByID retrieves a transaction by ID with details
func (r *TransactionRepository) FindByID(id int) (*models.Transaction, error) {
	var transaction models.Transaction
	err := config.DB.Preload("User").
		Preload("TransactionDetails").
		Preload("TransactionDetails.ProductUnit").
		Preload("TransactionDetails.ProductUnit.Product").
		First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// FindByInvoiceNo retrieves a transaction by invoice number
func (r *TransactionRepository) FindByInvoiceNo(invoiceNo string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := config.DB.Preload("User").
		Preload("TransactionDetails").
		Preload("TransactionDetails.ProductUnit").
		Preload("TransactionDetails.ProductUnit.Product").
		Where("invoice_no = ?", invoiceNo).
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// FindByUserID retrieves transactions for a specific user (cashier)
func (r *TransactionRepository) FindByUserID(userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := config.DB.Preload("User").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&transactions).Error
	return transactions, err
}

// FindByDateRange retrieves transactions within a date range
func (r *TransactionRepository) FindByDateRange(startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := config.DB.Preload("User").
		Preload("TransactionDetails").
		Where("transaction_date BETWEEN ? AND ?", startDate, endDate).
		Order("id DESC").
		Find(&transactions).Error
	return transactions, err
}

// GetTodayTransactions retrieves today's transactions
func (r *TransactionRepository) GetTodayTransactions() ([]models.Transaction, error) {
	today := time.Now().Format("2006-01-02")
	var transactions []models.Transaction
	err := config.DB.Preload("User").
		Where("DATE(transaction_date) = ?", today).
		Order("id DESC").
		Find(&transactions).Error
	return transactions, err
}

// GetTotalSalesToday returns total sales amount for today
func (r *TransactionRepository) GetTotalSalesToday() (decimal.Decimal, error) {
	today := time.Now().Format("2006-01-02")
	var result struct {
		Total decimal.Decimal
	}
	err := config.DB.Model(&models.Transaction{}).
		Select("COALESCE(SUM(total_amount), 0) as total").
		Where("DATE(transaction_date) = ?", today).
		Scan(&result).Error
	return result.Total, err
}

// generateInvoiceNo generates a unique invoice number
// Format: INV-YYYYMMDD-XXXX
func generateInvoiceNo() string {
	now := time.Now()
	dateStr := now.Format("20060102")

	// Get count of transactions today for sequence
	var count int64
	config.DB.Model(&models.Transaction{}).
		Where("DATE(transaction_date) = ?", now.Format("2006-01-02")).
		Count(&count)

	return fmt.Sprintf("INV-%s-%04d", dateStr, count+1)
}
