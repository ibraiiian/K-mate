package handlers

import (
	"strconv"

	"k-mate-backend/internal/dto"
	"k-mate-backend/internal/middleware"
	"k-mate-backend/internal/models"
	"k-mate-backend/internal/services"
	"k-mate-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

// POSHandler handles POS HTTP requests
type POSHandler struct {
	posService *services.POSService
}

// NewPOSHandler creates a new POSHandler instance
func NewPOSHandler(posService *services.POSService) *POSHandler {
	return &POSHandler{
		posService: posService,
	}
}

// Checkout handles POST /pos/checkout
// @Summary Process checkout and create transaction
// @Tags POS
// @Accept json
// @Produce json
// @Param body body dto.CheckoutRequest true "Checkout Data"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Router /pos/checkout [post]
func (h *POSHandler) Checkout(c *fiber.Ctx) error {
	var req dto.CheckoutRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Format JSON tidak valid")
	}

	// Validation
	if len(req.Items) == 0 {
		return utils.BadRequestResponse(c, "Keranjang kosong")
	}
	if req.PaymentMethod == "" {
		return utils.BadRequestResponse(c, "Metode pembayaran wajib diisi")
	}

	// Get cashier info from JWT context
	userID := middleware.GetUserID(c)
	cashierName := middleware.GetUsername(c)

	if userID == 0 {
		return utils.UnauthorizedResponse(c, "User tidak terautentikasi")
	}

	// Process checkout
	response, err := h.posService.Checkout(userID, cashierName, &req)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	// Log checkout activity
	utils.LogActivityAsync(utils.AuditLogParams{
		UserID:      userID,
		Action:      models.ActionCheckout,
		TargetTable: "transactions",
		TargetID:    response.TransactionID,
		NewData: map[string]interface{}{
			"invoice_no":     response.InvoiceNo,
			"total_amount":   response.TotalAmount,
			"payment_method": response.PaymentMethod,
			"items_count":    len(response.Items),
		},
		Context: c,
	})

	return utils.CreatedResponse(c, "Transaksi berhasil", response)
}

// GetAllTransactions handles GET /pos/transactions
// @Summary Get all transactions
// @Tags POS
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} utils.APIResponse
// @Router /pos/transactions [get]
func (h *POSHandler) GetAllTransactions(c *fiber.Ctx) error {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr != "" || limitStr != "" {
		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)

		transactions, total, err := h.posService.GetTransactionsPaginated(page, limit)
		if err != nil {
			return utils.InternalErrorResponse(c, "Gagal mengambil data transaksi")
		}

		return utils.SuccessResponse(c, "Daftar transaksi", fiber.Map{
			"transactions": toTransactionListResponse(transactions),
			"total":        total,
			"page":         page,
			"limit":        limit,
		})
	}

	transactions, err := h.posService.GetAllTransactions()
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data transaksi")
	}

	return utils.SuccessResponse(c, "Daftar transaksi", toTransactionListResponse(transactions))
}

// GetTransactionByID handles GET /pos/transactions/:id
// @Summary Get transaction by ID
// @Tags POS
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /pos/transactions/{id} [get]
func (h *POSHandler) GetTransactionByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "ID transaksi tidak valid")
	}

	transaction, err := h.posService.GetTransactionByID(id)
	if err != nil {
		return utils.NotFoundResponse(c, "Transaksi tidak ditemukan")
	}

	return utils.SuccessResponse(c, "Detail transaksi", toTransactionDetailResponse(transaction))
}

// GetTransactionByInvoice handles GET /pos/transactions/invoice/:invoiceNo
// @Summary Get transaction by invoice number
// @Tags POS
// @Produce json
// @Param invoiceNo path string true "Invoice Number"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /pos/transactions/invoice/{invoiceNo} [get]
func (h *POSHandler) GetTransactionByInvoice(c *fiber.Ctx) error {
	invoiceNo := c.Params("invoiceNo")
	if invoiceNo == "" {
		return utils.BadRequestResponse(c, "Nomor invoice wajib diisi")
	}

	transaction, err := h.posService.GetTransactionByInvoice(invoiceNo)
	if err != nil {
		return utils.NotFoundResponse(c, "Transaksi tidak ditemukan")
	}

	return utils.SuccessResponse(c, "Detail transaksi", toTransactionDetailResponse(transaction))
}

// GetTodayTransactions handles GET /pos/transactions/today
// @Summary Get today's transactions
// @Tags POS
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Router /pos/transactions/today [get]
func (h *POSHandler) GetTodayTransactions(c *fiber.Ctx) error {
	transactions, err := h.posService.GetTodayTransactions()
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data transaksi")
	}

	return utils.SuccessResponse(c, "Transaksi hari ini", toTransactionListResponse(transactions))
}

// GetMyTransactions handles GET /pos/transactions/my
// @Summary Get current user's transactions
// @Tags POS
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Router /pos/transactions/my [get]
func (h *POSHandler) GetMyTransactions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.UnauthorizedResponse(c, "User tidak terautentikasi")
	}

	transactions, err := h.posService.GetCashierTransactions(userID)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data transaksi")
	}

	return utils.SuccessResponse(c, "Transaksi saya", toTransactionListResponse(transactions))
}

// GetTodaySummary handles GET /pos/summary/today
// @Summary Get today's sales summary
// @Tags POS
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Router /pos/summary/today [get]
func (h *POSHandler) GetTodaySummary(c *fiber.Ctx) error {
	total, err := h.posService.GetTodaySalesTotal()
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil summary")
	}

	transactions, _ := h.posService.GetTodayTransactions()

	return utils.SuccessResponse(c, "Summary hari ini", fiber.Map{
		"total_sales":       total,
		"transaction_count": len(transactions),
	})
}

// ================== HELPER FUNCTIONS ==================

func toTransactionListResponse(transactions []models.Transaction) []dto.TransactionListResponse {
	result := make([]dto.TransactionListResponse, len(transactions))
	for i, t := range transactions {
		cashierName := ""
		if t.User != nil {
			cashierName = t.User.FullName
		}
		result[i] = dto.TransactionListResponse{
			ID:              t.ID,
			InvoiceNo:       t.InvoiceNo,
			CashierName:     cashierName,
			TotalAmount:     t.TotalAmount,
			PaymentMethod:   string(t.PaymentMethod),
			TransactionDate: t.TransactionDate.Format("2006-01-02 15:04:05"),
			ItemCount:       len(t.TransactionDetails),
		}
	}
	return result
}

func toTransactionDetailResponse(t *models.Transaction) fiber.Map {
	cashierName := ""
	if t.User != nil {
		cashierName = t.User.FullName
	}

	items := make([]fiber.Map, len(t.TransactionDetails))
	totalCost := decimal.Zero

	for i, d := range t.TransactionDetails {
		productName := ""
		unitName := ""
		conversionFactor := 0

		if d.ProductUnit != nil {
			unitName = d.ProductUnit.UnitName
			conversionFactor = d.ProductUnit.ConversionFactor
			if d.ProductUnit.Product != nil {
				productName = d.ProductUnit.Product.Name
			}
		}

		costTotal := d.CostAtMoment.Mul(decimal.NewFromInt(int64(d.Qty)))
		totalCost = totalCost.Add(costTotal)

		items[i] = fiber.Map{
			"id":                d.ID,
			"product_unit_id":   d.ProductUnitID,
			"product_name":      productName,
			"unit_name":         unitName,
			"qty":               d.Qty,
			"price_at_moment":   d.PriceAtMoment,
			"cost_at_moment":    d.CostAtMoment,
			"subtotal":          d.Subtotal,
			"conversion_factor": conversionFactor,
		}
	}

	grossProfit := t.TotalAmount.Sub(totalCost)

	return fiber.Map{
		"id":               t.ID,
		"invoice_no":       t.InvoiceNo,
		"cashier_name":     cashierName,
		"payment_method":   string(t.PaymentMethod),
		"total_amount":     t.TotalAmount,
		"total_cost":       totalCost,
		"gross_profit":     grossProfit,
		"transaction_date": t.TransactionDate.Format("2006-01-02 15:04:05"),
		"items":            items,
	}
}
