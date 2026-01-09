package dto

import "github.com/shopspring/decimal"

// CreateReturnRequest represents the request payload for creating a return
type CreateReturnRequest struct {
	TransactionID int                 `json:"transaction_id" validate:"required"`
	Items         []ReturnItemRequest `json:"items" validate:"required,min=1"`
	TotalRefund   decimal.Decimal     `json:"total_refund"`
	Reason        string              `json:"reason"`
}

// ReturnItemRequest represents a single item in the return request
type ReturnItemRequest struct {
	TransactionDetailID int    `json:"transaction_detail_id" validate:"required"`
	Qty                 int    `json:"qty" validate:"required,min=1"`
	Condition           string `json:"condition" validate:"required,oneof=good damaged"`
}

// ReturnResponse represents the response after a successful return
type ReturnResponse struct {
	ID          int                  `json:"id"`
	ReturnNo    string               `json:"return_no"`
	TotalRefund decimal.Decimal      `json:"total_refund"`
	ReturnDate  string               `json:"return_date"`
	Items       []ReturnItemResponse `json:"items"`
}

// ReturnItemResponse represents a single item in the return response
type ReturnItemResponse struct {
	ProductName    string          `json:"product_name"`
	UnitName       string          `json:"unit_name"`
	Qty            int             `json:"qty"`
	RefundPerUnit  decimal.Decimal `json:"refund_per_unit"`
	SubtotalRefund decimal.Decimal `json:"subtotal_refund"`
	Condition      string          `json:"condition"`
	StockRestored  bool            `json:"stock_restored"`
}
