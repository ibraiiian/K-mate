package services

import (
	"errors"
	"time"

	"k-mate-backend/internal/dto"
	"k-mate-backend/internal/models"
	"k-mate-backend/internal/repository"

	"github.com/shopspring/decimal"
)

// POSService handles POS business logic
type POSService struct {
	transactionRepo *repository.TransactionRepository
	productRepo     *repository.ProductRepository
}

// NewPOSService creates a new POSService instance
func NewPOSService(
	transactionRepo *repository.TransactionRepository,
	productRepo *repository.ProductRepository,
) *POSService {
	return &POSService{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

// Checkout processes a checkout request and creates a transaction
func (s *POSService) Checkout(userID int, cashierName string, req *dto.CheckoutRequest) (*dto.CheckoutResponse, error) {
	// 1. Validate request
	if len(req.Items) == 0 {
		return nil, errors.New("keranjang kosong, tidak ada item untuk checkout")
	}

	// Validate payment method
	validPaymentMethods := map[string]bool{
		"cash": true, "transfer": true, "qris": true, "credit": true,
	}
	if !validPaymentMethods[req.PaymentMethod] {
		return nil, errors.New("metode pembayaran tidak valid (cash/transfer/qris/credit)")
	}

	// 2. Prepare transaction items with validation
	var transactionItems []repository.TransactionItem
	var responseItems []dto.TransactionItemResponse
	totalAmount := decimal.Zero
	totalCost := decimal.Zero

	for _, item := range req.Items {
		if item.Qty <= 0 {
			return nil, errors.New("qty harus lebih dari 0")
		}

		// Get product unit with product info
		unit, err := s.productRepo.FindUnitByID(item.UnitID)
		if err != nil {
			return nil, errors.New("unit produk tidak ditemukan")
		}

		// Verify product_id matches
		if unit.ProductID != item.ProductID {
			return nil, errors.New("product_id tidak cocok dengan unit_id")
		}

		// Get product for stock check
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, errors.New("produk tidak ditemukan")
		}

		// Calculate stock to deduct (in base unit)
		stockToDeduct := item.Qty * unit.ConversionFactor

		// Check if stock is sufficient
		if product.StockQty < stockToDeduct {
			return nil, errors.New("stok tidak cukup untuk " + product.Name +
				". Tersedia: " + string(rune(product.StockQty)) + " " + product.BaseUnit)
		}

		// Calculate subtotal
		subtotal := unit.SellPrice.Mul(decimal.NewFromInt(int64(item.Qty)))
		costTotal := unit.BuyPrice.Mul(decimal.NewFromInt(int64(item.Qty)))

		totalAmount = totalAmount.Add(subtotal)
		totalCost = totalCost.Add(costTotal)

		// Prepare transaction item
		transactionItems = append(transactionItems, repository.TransactionItem{
			ProductID:        item.ProductID,
			ProductName:      product.Name,
			ProductUnitID:    item.UnitID,
			UnitName:         unit.UnitName,
			Qty:              item.Qty,
			ConversionFactor: unit.ConversionFactor,
			SellPrice:        unit.SellPrice, // Price snapshot
			BuyPrice:         unit.BuyPrice,  // HPP snapshot
			Subtotal:         subtotal,
			StockToDeduct:    stockToDeduct,
		})

		// Prepare response item
		responseItems = append(responseItems, dto.TransactionItemResponse{
			ProductID:        item.ProductID,
			ProductName:      product.Name,
			UnitID:           item.UnitID,
			UnitName:         unit.UnitName,
			Qty:              item.Qty,
			PricePerUnit:     unit.SellPrice,
			CostPerUnit:      unit.BuyPrice,
			Subtotal:         subtotal,
			ConversionFactor: unit.ConversionFactor,
			StockDeducted:    stockToDeduct,
		})
	}

	// 3. Create transaction model
	transaction := &models.Transaction{
		UserID:          userID,
		TotalAmount:     totalAmount,
		PaymentMethod:   models.PaymentMethod(req.PaymentMethod),
		TransactionDate: time.Now(),
	}

	// 4. Save transaction with details and update stock (atomic)
	if err := s.transactionRepo.CreateTransaction(transaction, transactionItems); err != nil {
		return nil, err
	}

	// 5. Build response
	grossProfit := totalAmount.Sub(totalCost)

	response := &dto.CheckoutResponse{
		InvoiceNo:       transaction.InvoiceNo,
		TransactionID:   transaction.ID,
		CashierName:     cashierName,
		CustomerName:    req.CustomerName,
		PaymentMethod:   req.PaymentMethod,
		Items:           responseItems,
		TotalAmount:     totalAmount,
		TotalCost:       totalCost,
		GrossProfit:     grossProfit,
		TransactionDate: transaction.TransactionDate.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

// GetTransactionByID retrieves a transaction by ID
func (s *POSService) GetTransactionByID(id int) (*models.Transaction, error) {
	return s.transactionRepo.FindByID(id)
}

// GetTransactionByInvoice retrieves a transaction by invoice number
func (s *POSService) GetTransactionByInvoice(invoiceNo string) (*models.Transaction, error) {
	return s.transactionRepo.FindByInvoiceNo(invoiceNo)
}

// GetAllTransactions retrieves all transactions
func (s *POSService) GetAllTransactions() ([]models.Transaction, error) {
	return s.transactionRepo.FindAll()
}

// GetTransactionsPaginated retrieves transactions with pagination
func (s *POSService) GetTransactionsPaginated(page, limit int) ([]models.Transaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.transactionRepo.FindAllPaginated(page, limit)
}

// GetTodayTransactions retrieves today's transactions
func (s *POSService) GetTodayTransactions() ([]models.Transaction, error) {
	return s.transactionRepo.GetTodayTransactions()
}

// GetCashierTransactions retrieves transactions for a specific cashier
func (s *POSService) GetCashierTransactions(userID int) ([]models.Transaction, error) {
	return s.transactionRepo.FindByUserID(userID)
}

// GetTodaySalesTotal returns total sales amount for today
func (s *POSService) GetTodaySalesTotal() (decimal.Decimal, error) {
	return s.transactionRepo.GetTotalSalesToday()
}
