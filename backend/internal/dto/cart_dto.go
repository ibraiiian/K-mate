package dto

import "github.com/shopspring/decimal"

// ================== REQUEST DTOs ==================

// CartItemRequest represents a single item in the cart
type CartItemRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	UnitID    int `json:"unit_id" validate:"required"`
	Qty       int `json:"qty" validate:"required,gt=0"`
}

// CheckoutRequest represents the JSON body for checkout
type CheckoutRequest struct {
	CustomerName  *string           `json:"customer_name"`
	PaymentMethod string            `json:"payment_method" validate:"required,oneof=cash transfer qris credit"`
	Items         []CartItemRequest `json:"items" validate:"required,min=1,dive"`
}

// ================== RESPONSE DTOs ==================

// TransactionItemResponse represents a single item in transaction response
type TransactionItemResponse struct {
	ProductID        int             `json:"product_id"`
	ProductName      string          `json:"product_name"`
	UnitID           int             `json:"unit_id"`
	UnitName         string          `json:"unit_name"`
	Qty              int             `json:"qty"`
	PricePerUnit     decimal.Decimal `json:"price_per_unit"`
	CostPerUnit      decimal.Decimal `json:"cost_per_unit"`
	Subtotal         decimal.Decimal `json:"subtotal"`
	ConversionFactor int             `json:"conversion_factor"`
	StockDeducted    int             `json:"stock_deducted"` // In base unit
}

// CheckoutResponse represents the response after successful checkout
type CheckoutResponse struct {
	InvoiceNo       string                    `json:"invoice_no"`
	TransactionID   int                       `json:"transaction_id"`
	CashierName     string                    `json:"cashier_name"`
	CustomerName    *string                   `json:"customer_name"`
	PaymentMethod   string                    `json:"payment_method"`
	Items           []TransactionItemResponse `json:"items"`
	TotalAmount     decimal.Decimal           `json:"total_amount"`
	TotalCost       decimal.Decimal           `json:"total_cost"`   // HPP Total
	GrossProfit     decimal.Decimal           `json:"gross_profit"` // Total - HPP
	TransactionDate string                    `json:"transaction_date"`
}

// TransactionListResponse represents a transaction in list view
type TransactionListResponse struct {
	ID              int             `json:"id"`
	InvoiceNo       string          `json:"invoice_no"`
	CashierName     string          `json:"cashier_name"`
	TotalAmount     decimal.Decimal `json:"total_amount"`
	PaymentMethod   string          `json:"payment_method"`
	TransactionDate string          `json:"transaction_date"`
	ItemCount       int             `json:"item_count"`
}
