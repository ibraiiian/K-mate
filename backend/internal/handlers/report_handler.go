package handlers

import (
	"time"

	"k-mate-backend/internal/dto"
	"k-mate-backend/internal/repository"
	"k-mate-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ReportHandler handles report/dashboard HTTP requests
type ReportHandler struct {
	reportRepo *repository.ReportRepository
}

// NewReportHandler creates a new ReportHandler instance
func NewReportHandler(reportRepo *repository.ReportRepository) *ReportHandler {
	return &ReportHandler{
		reportRepo: reportRepo,
	}
}

// GetDailySummary handles GET /dashboard/summary
// Returns today's financial summary with Net Profit calculation
func (h *ReportHandler) GetDailySummary(c *fiber.Ctx) error {
	// Check for custom date parameter
	dateStr := c.Query("date")
	var targetDate time.Time

	if dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return utils.BadRequestResponse(c, "Format tanggal tidak valid (gunakan YYYY-MM-DD)")
		}
		targetDate = parsed
	} else {
		targetDate = time.Now()
	}

	summary, err := h.reportRepo.GetDailySummary(targetDate)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil summary")
	}

	// Calculate profits
	grossProfit := summary.TotalRevenue.Sub(summary.TotalHPP)
	netProfit := grossProfit.Sub(summary.TotalExpense)

	return utils.SuccessResponse(c, "Summary harian", dto.DailySummaryResponse{
		Date:             targetDate.Format("2006-01-02"),
		TotalRevenue:     summary.TotalRevenue,
		TotalHPP:         summary.TotalHPP,
		GrossProfit:      grossProfit,
		TotalExpense:     summary.TotalExpense,
		NetProfit:        netProfit,
		TransactionCount: summary.TransactionCount,
	})
}

// GetSalesChart handles GET /dashboard/chart
// Returns sales data for the last 7 days
func (h *ReportHandler) GetSalesChart(c *fiber.Ctx) error {
	// Default 7 days, can be overridden
	days := 7
	daysParam := c.Query("days")
	if daysParam != "" {
		if d, err := parseInt(daysParam); err == nil && d > 0 && d <= 30 {
			days = d
		}
	}

	chartData, err := h.reportRepo.GetSalesChart(days)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data chart")
	}

	// Convert to DTO
	result := make([]dto.SalesChartDataPoint, len(chartData))
	for i, data := range chartData {
		profit := data.Revenue.Sub(data.HPP)
		result[i] = dto.SalesChartDataPoint{
			Date:    data.Date,
			Revenue: data.Revenue,
			HPP:     data.HPP,
			Profit:  profit,
			Count:   data.Count,
		}
	}

	return utils.SuccessResponse(c, "Data chart penjualan", fiber.Map{
		"days": days,
		"data": result,
	})
}

// GetTopProducts handles GET /dashboard/top-products
// Returns top 5 best-selling products
func (h *ReportHandler) GetTopProducts(c *fiber.Ctx) error {
	// Default: this month
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := now

	// Check for custom date range
	startStr := c.Query("start_date")
	endStr := c.Query("end_date")

	if startStr != "" {
		if parsed, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = parsed
		}
	}
	if endStr != "" {
		if parsed, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = parsed.Add(24*time.Hour - time.Second) // End of day
		}
	}

	// Default limit 5
	limit := 5
	limitParam := c.Query("limit")
	if limitParam != "" {
		if l, err := parseInt(limitParam); err == nil && l > 0 && l <= 20 {
			limit = l
		}
	}

	topProducts, err := h.reportRepo.GetTopProducts(limit, startDate, endDate)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data produk terlaris")
	}

	// Convert to DTO
	result := make([]dto.TopProductResponse, len(topProducts))
	for i, p := range topProducts {
		result[i] = dto.TopProductResponse{
			ProductID:   p.ProductID,
			ProductName: p.ProductName,
			TotalQty:    p.TotalQty,
			TotalSales:  p.TotalSales,
		}
	}

	return utils.SuccessResponse(c, "Produk terlaris", fiber.Map{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"products":   result,
	})
}

// GetLowStockAlert handles GET /dashboard/low-stock
// Returns products with stock below minimum alert
func (h *ReportHandler) GetLowStockAlert(c *fiber.Ctx) error {
	lowStock, err := h.reportRepo.GetLowStockProducts()
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data stok")
	}

	return utils.SuccessResponse(c, "Peringatan stok rendah", fiber.Map{
		"count":    len(lowStock),
		"products": lowStock,
	})
}

// parseInt helper function
func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fiber.ErrBadRequest
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

// GetFinancialReport handles GET /reports
// Returns period-based financial report with expense breakdown
func (h *ReportHandler) GetFinancialReport(c *fiber.Ctx) error {
	// Get date parameters (default: first day of current month to today)
	now := time.Now()
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Default values
	if startDate == "" {
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = now.Format("2006-01-02")
	}

	// Validate date format
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return utils.BadRequestResponse(c, "Format start_date tidak valid (gunakan YYYY-MM-DD)")
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return utils.BadRequestResponse(c, "Format end_date tidak valid (gunakan YYYY-MM-DD)")
	}

	// Get financial report data
	reportData, err := h.reportRepo.GetFinancialReportByPeriod(startDate, endDate)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data laporan")
	}

	// Get expense breakdown by category
	expenseBreakdown, err := h.reportRepo.GetExpenseBreakdownByCategory(startDate, endDate)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data kategori pengeluaran")
	}

	// Calculate profits
	grossProfit := reportData.TotalRevenue.Sub(reportData.TotalHPP)
	netProfit := grossProfit.Sub(reportData.TotalExpense)

	// Convert expense breakdown to DTO
	expenseBreakdownDTO := make([]dto.ExpenseCategoryBreakdownDTO, len(expenseBreakdown))
	for i, eb := range expenseBreakdown {
		expenseBreakdownDTO[i] = dto.ExpenseCategoryBreakdownDTO{
			Category: eb.Category,
			Total:    eb.Total,
			Count:    eb.Count,
		}
	}

	return utils.SuccessResponse(c, "Laporan Keuangan", dto.FinancialReportResponse{
		StartDate:        startDate,
		EndDate:          endDate,
		TotalRevenue:     reportData.TotalRevenue,
		TotalHPP:         reportData.TotalHPP,
		GrossProfit:      grossProfit,
		TotalExpense:     reportData.TotalExpense,
		NetProfit:        netProfit,
		TransactionCount: reportData.TransactionCount,
		ExpenseBreakdown: expenseBreakdownDTO,
	})
}
