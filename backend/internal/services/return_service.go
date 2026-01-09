package services

import (
	"errors"
	"time"

	"k-mate-backend/internal/dto"
	"k-mate-backend/internal/models"
	"k-mate-backend/internal/repository"

	"github.com/shopspring/decimal"
)

// ReturnService handles return business logic
type ReturnService struct {
	returnRepo      *repository.ReturnRepository
	transactionRepo *repository.TransactionRepository
	productRepo     *repository.ProductRepository
}

// NewReturnService creates a new ReturnService instance
func NewReturnService(
	returnRepo *repository.ReturnRepository,
	transactionRepo *repository.TransactionRepository,
	productRepo *repository.ProductRepository,
) *ReturnService {
	return &ReturnService{
		returnRepo:      returnRepo,
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

// ProcessReturn processes a return request
func (s *ReturnService) ProcessReturn(userID int, req *dto.CreateReturnRequest) (*dto.ReturnResponse, error) {
	// 1. Validate request
	if len(req.Items) == 0 {
		return nil, errors.New("tidak ada item untuk diretur")
	}

	// 2. Get the original transaction
	transaction, err := s.transactionRepo.FindByID(req.TransactionID)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	// 3. Build a map of transaction details for validation
	detailMap := make(map[int]*models.TransactionDetail)
	for i := range transaction.TransactionDetails {
		detail := &transaction.TransactionDetails[i]
		detailMap[detail.ID] = detail
	}

	// 4. Prepare return items with validation
	var returnItems []repository.ReturnItem
	var responseItems []dto.ReturnItemResponse
	totalRefund := decimal.Zero

	for _, item := range req.Items {
		// Validate transaction detail exists
		detail, exists := detailMap[item.TransactionDetailID]
		if !exists {
			return nil, errors.New("detail transaksi tidak valid")
		}

		// Validate qty doesn't exceed original qty
		if item.Qty > detail.Qty {
			return nil, errors.New("qty retur melebihi qty pembelian")
		}

		// Validate condition
		if item.Condition != "good" && item.Condition != "damaged" {
			return nil, errors.New("kondisi harus 'good' atau 'damaged'")
		}

		// Get product unit for stock calculation
		unit, err := s.productRepo.FindUnitByID(detail.ProductUnitID)
		if err != nil {
			return nil, errors.New("unit produk tidak ditemukan")
		}

		// Calculate refund
		refundPerUnit := detail.PriceAtMoment
		subtotalRefund := refundPerUnit.Mul(decimal.NewFromInt(int64(item.Qty)))
		totalRefund = totalRefund.Add(subtotalRefund)

		// Calculate stock to restore (in base unit)
		stockToRestore := 0
		if item.Condition == "good" {
			stockToRestore = item.Qty * unit.ConversionFactor
		}

		// Get product name
		productName := "Unknown"
		unitName := unit.UnitName
		if unit.Product != nil {
			productName = unit.Product.Name
		}

		returnItems = append(returnItems, repository.ReturnItem{
			TransactionDetailID: item.TransactionDetailID,
			ProductUnitID:       detail.ProductUnitID,
			ProductID:           unit.ProductID,
			ProductName:         productName,
			UnitName:            unitName,
			Qty:                 item.Qty,
			RefundPricePerUnit:  refundPerUnit,
			SubtotalRefund:      subtotalRefund,
			Condition:           item.Condition,
			ConversionFactor:    unit.ConversionFactor,
			StockToRestore:      stockToRestore,
		})

		responseItems = append(responseItems, dto.ReturnItemResponse{
			ProductName:    productName,
			UnitName:       unitName,
			Qty:            item.Qty,
			RefundPerUnit:  refundPerUnit,
			SubtotalRefund: subtotalRefund,
			Condition:      item.Condition,
			StockRestored:  item.Condition == "good",
		})
	}

	// 5. Create return record
	returnRecord := &models.Return{
		TransactionID: req.TransactionID,
		UserID:        userID,
		TotalRefund:   totalRefund,
		ReturnDate:    time.Now(),
	}

	if req.Reason != "" {
		returnRecord.Reason = &req.Reason
	}

	// 6. Save to database
	if err := s.returnRepo.CreateReturn(returnRecord, returnItems); err != nil {
		return nil, err
	}

	// 7. Build response
	response := &dto.ReturnResponse{
		ID:          returnRecord.ID,
		ReturnNo:    returnRecord.ReturnNo,
		TotalRefund: totalRefund,
		ReturnDate:  returnRecord.ReturnDate.Format("2006-01-02 15:04:05"),
		Items:       responseItems,
	}

	return response, nil
}

// GetAllReturns retrieves all returns
func (s *ReturnService) GetAllReturns() ([]models.Return, error) {
	return s.returnRepo.FindAll()
}

// GetReturnByID retrieves a return by ID
func (s *ReturnService) GetReturnByID(id int) (*models.Return, error) {
	return s.returnRepo.FindByID(id)
}
