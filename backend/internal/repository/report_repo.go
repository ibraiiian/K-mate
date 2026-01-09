package repository

import (
	"time"

	"k-mate-backend/config"

	"github.com/shopspring/decimal"
)

// ReportRepository handles database operations for reports and analytics
type ReportRepository struct{}

// NewReportRepository creates a new ReportRepository instance
func NewReportRepository() *ReportRepository {
	return &ReportRepository{}
}

// DailySummary represents daily financial summary from database
type DailySummary struct {
	TotalRevenue     decimal.Decimal
	TotalHPP         decimal.Decimal
	TotalExpense     decimal.Decimal
	TransactionCount int
}

// GetDailySummary calculates financial summary for a specific date
func (r *ReportRepository) GetDailySummary(date time.Time) (*DailySummary, error) {
	dateStr := date.Format("2006-01-02")
	summary := &DailySummary{}

	// 1. Get Total Revenue and Transaction Count
	var revenueResult struct {
		Total decimal.Decimal
		Count int
	}
	err := config.DB.Raw(`
		SELECT 
			COALESCE(SUM(total_amount), 0) as total,
			COUNT(*) as count
		FROM transactions 
		WHERE DATE(transaction_date) = ?
	`, dateStr).Scan(&revenueResult).Error
	if err != nil {
		return nil, err
	}
	summary.TotalRevenue = revenueResult.Total
	summary.TransactionCount = revenueResult.Count

	// 2. Get Total HPP (Cost of Goods Sold)
	var hppResult struct {
		Total decimal.Decimal
	}
	err = config.DB.Raw(`
		SELECT COALESCE(SUM(td.cost_at_moment * td.qty), 0) as total
		FROM transaction_details td
		INNER JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.transaction_date) = ?
	`, dateStr).Scan(&hppResult).Error
	if err != nil {
		return nil, err
	}
	summary.TotalHPP = hppResult.Total

	// 3. Get Total Expense
	var expenseResult struct {
		Total decimal.Decimal
	}
	err = config.DB.Raw(`
		SELECT COALESCE(SUM(amount), 0) as total
		FROM expenses 
		WHERE expense_date = ?
	`, dateStr).Scan(&expenseResult).Error
	if err != nil {
		return nil, err
	}
	summary.TotalExpense = expenseResult.Total

	return summary, nil
}

// SalesChartData represents sales data for a single day
type SalesChartData struct {
	Date    string
	Revenue decimal.Decimal
	HPP     decimal.Decimal
	Count   int
}

// GetSalesChart returns sales data for the last N days
func (r *ReportRepository) GetSalesChart(days int) ([]SalesChartData, error) {
	var results []SalesChartData

	err := config.DB.Raw(`
		SELECT 
			DATE(t.transaction_date) as date,
			COALESCE(SUM(t.total_amount), 0) as revenue,
			COALESCE(SUM(td.cost_at_moment * td.qty), 0) as hpp,
			COUNT(DISTINCT t.id) as count
		FROM transactions t
		LEFT JOIN transaction_details td ON t.id = td.transaction_id
		WHERE t.transaction_date >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY DATE(t.transaction_date)
		ORDER BY date ASC
	`, days).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

// TopProduct represents a top selling product
type TopProduct struct {
	ProductID   int
	ProductName string
	TotalQty    int
	TotalSales  decimal.Decimal
}

// GetTopProducts returns top selling products
func (r *ReportRepository) GetTopProducts(limit int, startDate, endDate time.Time) ([]TopProduct, error) {
	var results []TopProduct

	err := config.DB.Raw(`
		SELECT 
			p.id as product_id,
			p.name as product_name,
			SUM(td.qty * pu.conversion_factor) as total_qty,
			SUM(td.subtotal) as total_sales
		FROM transaction_details td
		INNER JOIN transactions t ON td.transaction_id = t.id
		INNER JOIN product_units pu ON td.product_unit_id = pu.id
		INNER JOIN products p ON pu.product_id = p.id
		WHERE t.transaction_date BETWEEN ? AND ?
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT ?
	`, startDate, endDate, limit).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetMonthlyRevenue returns total revenue for a specific month
func (r *ReportRepository) GetMonthlyRevenue(year, month int) (decimal.Decimal, error) {
	var result struct {
		Total decimal.Decimal
	}
	err := config.DB.Raw(`
		SELECT COALESCE(SUM(total_amount), 0) as total
		FROM transactions 
		WHERE YEAR(transaction_date) = ? AND MONTH(transaction_date) = ?
	`, year, month).Scan(&result).Error
	return result.Total, err
}

// GetMonthlyExpense returns total expense for a specific month
func (r *ReportRepository) GetMonthlyExpense(year, month int) (decimal.Decimal, error) {
	var result struct {
		Total decimal.Decimal
	}
	err := config.DB.Raw(`
		SELECT COALESCE(SUM(amount), 0) as total
		FROM expenses 
		WHERE YEAR(expense_date) = ? AND MONTH(expense_date) = ?
	`, year, month).Scan(&result).Error
	return result.Total, err
}

// LowStockProduct represents a product with low stock
type LowStockProduct struct {
	ID            int
	Name          string
	StockQty      int
	MinStockAlert int
	BaseUnit      string
}

// GetLowStockProducts returns products with stock below minimum alert
func (r *ReportRepository) GetLowStockProducts() ([]LowStockProduct, error) {
	var results []LowStockProduct
	err := config.DB.Raw(`
		SELECT 
			id, name, stock_qty, min_stock_alert, base_unit
		FROM products 
		WHERE stock_qty <= min_stock_alert
		ORDER BY stock_qty ASC
	`).Scan(&results).Error
	return results, err
}

// FinancialReportData represents aggregated financial data for a period
type FinancialReportData struct {
	TotalRevenue     decimal.Decimal
	TotalHPP         decimal.Decimal
	TotalExpense     decimal.Decimal
	TransactionCount int
}

// GetFinancialReportByPeriod calculates financial summary for a date range
func (r *ReportRepository) GetFinancialReportByPeriod(startDate, endDate string) (*FinancialReportData, error) {
	report := &FinancialReportData{}

	// 1. Get Total Revenue and Transaction Count
	var revenueResult struct {
		Total decimal.Decimal
		Count int
	}
	err := config.DB.Raw(`
		SELECT 
			COALESCE(SUM(total_amount), 0) as total,
			COUNT(*) as count
		FROM transactions 
		WHERE DATE(transaction_date) BETWEEN ? AND ?
	`, startDate, endDate).Scan(&revenueResult).Error
	if err != nil {
		return nil, err
	}
	report.TotalRevenue = revenueResult.Total
	report.TransactionCount = revenueResult.Count

	// 2. Get Total HPP (Cost of Goods Sold)
	var hppResult struct {
		Total decimal.Decimal
	}
	err = config.DB.Raw(`
		SELECT COALESCE(SUM(td.cost_at_moment * td.qty), 0) as total
		FROM transaction_details td
		INNER JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.transaction_date) BETWEEN ? AND ?
	`, startDate, endDate).Scan(&hppResult).Error
	if err != nil {
		return nil, err
	}
	report.TotalHPP = hppResult.Total

	// 3. Get Total Expense
	var expenseResult struct {
		Total decimal.Decimal
	}
	err = config.DB.Raw(`
		SELECT COALESCE(SUM(amount), 0) as total
		FROM expenses 
		WHERE expense_date BETWEEN ? AND ?
	`, startDate, endDate).Scan(&expenseResult).Error
	if err != nil {
		return nil, err
	}
	report.TotalExpense = expenseResult.Total

	return report, nil
}

// ExpenseCategoryBreakdown represents expense breakdown by category
type ExpenseCategoryBreakdown struct {
	Category string          `json:"category"`
	Total    decimal.Decimal `json:"total"`
	Count    int             `json:"count"`
}

// GetExpenseBreakdownByCategory returns expenses grouped by category for a period
func (r *ReportRepository) GetExpenseBreakdownByCategory(startDate, endDate string) ([]ExpenseCategoryBreakdown, error) {
	var results []ExpenseCategoryBreakdown
	err := config.DB.Raw(`
		SELECT 
			category,
			COALESCE(SUM(amount), 0) as total,
			COUNT(*) as count
		FROM expenses 
		WHERE expense_date BETWEEN ? AND ?
		GROUP BY category
		ORDER BY total DESC
	`, startDate, endDate).Scan(&results).Error
	return results, err
}
