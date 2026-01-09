package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// ConditionStatus represents the condition of returned items
type ConditionStatus string

const (
	ConditionStatusGood    ConditionStatus = "good"    // Item will be restocked
	ConditionStatusDamaged ConditionStatus = "damaged" // Item will be discarded
)

// Return represents the returns table (Pengembalian Barang Header)
type Return struct {
	ID            int             `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID int             `gorm:"column:transaction_id;not null;index" json:"transaction_id"`
	UserID        int             `gorm:"column:user_id;not null;index" json:"user_id"`
	ReturnNo      string          `gorm:"column:return_no;type:varchar(50);uniqueIndex;not null" json:"return_no"`
	TotalRefund   decimal.Decimal `gorm:"column:total_refund;type:decimal(15,2);not null" json:"total_refund"`
	Reason        *string         `gorm:"column:reason;type:text" json:"reason"`
	ReturnDate    time.Time       `gorm:"column:return_date;autoCreateTime" json:"return_date"`

	// Relations
	Transaction   *Transaction   `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ReturnDetails []ReturnDetail `gorm:"foreignKey:ReturnID" json:"return_details,omitempty"`
}

// TableName specifies the table name for Return model
func (Return) TableName() string {
	return "returns"
}

// ReturnDetail represents the return_details table (Item Retur)
type ReturnDetail struct {
	ID                 int             `gorm:"primaryKey;autoIncrement" json:"id"`
	ReturnID           int             `gorm:"column:return_id;not null;index" json:"return_id"`
	ProductUnitID      int             `gorm:"column:product_unit_id;not null;index" json:"product_unit_id"`
	Qty                int             `gorm:"column:qty;not null" json:"qty"`
	RefundPricePerUnit decimal.Decimal `gorm:"column:refund_price_per_unit;type:decimal(15,2);not null" json:"refund_price_per_unit"`
	ConditionStatus    ConditionStatus `gorm:"column:condition_status;type:enum('good','damaged');not null" json:"condition_status"`
	SubtotalRefund     decimal.Decimal `gorm:"column:subtotal_refund;type:decimal(15,2);not null" json:"subtotal_refund"`

	// Relations
	Return      *Return      `gorm:"foreignKey:ReturnID" json:"return,omitempty"`
	ProductUnit *ProductUnit `gorm:"foreignKey:ProductUnitID" json:"product_unit,omitempty"`
}

// TableName specifies the table name for ReturnDetail model
func (ReturnDetail) TableName() string {
	return "return_details"
}
