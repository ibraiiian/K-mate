package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// PaymentMethod represents the allowed payment methods
type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodTransfer PaymentMethod = "transfer"
	PaymentMethodQRIS     PaymentMethod = "qris"
	PaymentMethodCredit   PaymentMethod = "credit" // Kasbon
)

// Transaction represents the transactions table (Sales Header)
type Transaction struct {
	ID              int             `gorm:"primaryKey;autoIncrement" json:"id"`
	InvoiceNo       string          `gorm:"column:invoice_no;type:varchar(50);uniqueIndex;not null" json:"invoice_no"`
	UserID          int             `gorm:"column:user_id;not null;index" json:"user_id"`
	TotalAmount     decimal.Decimal `gorm:"column:total_amount;type:decimal(15,2);not null" json:"total_amount"`
	PaymentMethod   PaymentMethod   `gorm:"column:payment_method;type:enum('cash','transfer','qris','credit');not null" json:"payment_method"`
	TransactionDate time.Time       `gorm:"column:transaction_date;autoCreateTime" json:"transaction_date"`

	// Relations
	User               *User               `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TransactionDetails []TransactionDetail `gorm:"foreignKey:TransactionID" json:"transaction_details,omitempty"`
	Returns            []Return            `gorm:"foreignKey:TransactionID" json:"returns,omitempty"`
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}

// TransactionDetail represents the transaction_details table (Sales Items & COGS Snapshot)
type TransactionDetail struct {
	ID            int             `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID int             `gorm:"column:transaction_id;not null;index" json:"transaction_id"`
	ProductUnitID int             `gorm:"column:product_unit_id;not null;index" json:"product_unit_id"`
	Qty           int             `gorm:"column:qty;not null" json:"qty"`
	PriceAtMoment decimal.Decimal `gorm:"column:price_at_moment;type:decimal(15,2);not null" json:"price_at_moment"`
	CostAtMoment  decimal.Decimal `gorm:"column:cost_at_moment;type:decimal(15,2);not null" json:"cost_at_moment"`
	Subtotal      decimal.Decimal `gorm:"column:subtotal;type:decimal(15,2);not null" json:"subtotal"`

	// Relations
	Transaction *Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	ProductUnit *ProductUnit `gorm:"foreignKey:ProductUnitID" json:"product_unit,omitempty"`
}

// TableName specifies the table name for TransactionDetail model
func (TransactionDetail) TableName() string {
	return "transaction_details"
}
