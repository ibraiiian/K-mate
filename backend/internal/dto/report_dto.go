package dto

import "github.com/shopspring/decimal"

// DailySummaryResponse represents daily financial summary
type DailySummaryResponse struct {
	Date             string          `json:"date"`
	TotalRevenue     decimal.Decimal `json:"total_revenue"`
	TotalHPP         decimal.Decimal `json:"total_hpp"`
	GrossProfit      decimal.Decimal `json:"gross_profit"`
	TotalExpense     decimal.Decimal `json:"total_expense"`
	NetProfit        decimal.Decimal `json:"net_profit"`
	TransactionCount int             `json:"transaction_count"`
}

// SalesChartDataPoint represents a single data point in sales chart
type SalesChartDataPoint struct {
	Date    string          `json:"date"`
	Revenue decimal.Decimal `json:"revenue"`
	HPP     decimal.Decimal `json:"hpp"`
	Profit  decimal.Decimal `json:"profit"`
	Count   int             `json:"count"`
}

// TopProductResponse represents a top selling product
type TopProductResponse struct {
	ProductID   int             `json:"product_id"`
	ProductName string          `json:"product_name"`
	TotalQty    int             `json:"total_qty"`
	TotalSales  decimal.Decimal `json:"total_sales"`
}

// FinancialReportResponse represents period-based financial report
type FinancialReportResponse struct {
	StartDate        string                        `json:"start_date"`
	EndDate          string                        `json:"end_date"`
	TotalRevenue     decimal.Decimal               `json:"total_revenue"`
	TotalHPP         decimal.Decimal               `json:"total_hpp"`
	GrossProfit      decimal.Decimal               `json:"gross_profit"`
	TotalExpense     decimal.Decimal               `json:"total_expense"`
	NetProfit        decimal.Decimal               `json:"net_profit"`
	TransactionCount int                           `json:"transaction_count"`
	ExpenseBreakdown []ExpenseCategoryBreakdownDTO `json:"expense_breakdown"`
}

// ExpenseCategoryBreakdownDTO represents expense breakdown by category
type ExpenseCategoryBreakdownDTO struct {
	Category string          `json:"category"`
	Total    decimal.Decimal `json:"total"`
	Count    int             `json:"count"`
}
