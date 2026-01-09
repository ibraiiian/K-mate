package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Product represents the products table (Base Unit Master Data)
type Product struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductCode   *string   `gorm:"column:product_code;type:varchar(50);uniqueIndex" json:"product_code"`
	Name          string    `gorm:"column:name;type:varchar(200);not null;index" json:"name"`
	BaseUnit      string    `gorm:"column:base_unit;type:varchar(20);not null" json:"base_unit"`
	StockQty      int       `gorm:"column:stock_qty;default:0" json:"stock_qty"`
	MinStockAlert int       `gorm:"column:min_stock_alert;default:10" json:"min_stock_alert"`
	ImagePath     string    `gorm:"column:image_path;type:varchar(255)" json:"image_path"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relations
	ProductUnits []ProductUnit `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product_units,omitempty"`
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}

// ProductUnit represents the product_units table (Conversion Logic & Pricing)
type ProductUnit struct {
	ID               int             `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID        int             `gorm:"column:product_id;not null;index" json:"product_id"`
	UnitName         string          `gorm:"column:unit_name;type:varchar(20);not null" json:"unit_name"`
	ConversionFactor int             `gorm:"column:conversion_factor;not null" json:"conversion_factor"`
	BuyPrice         decimal.Decimal `gorm:"column:buy_price;type:decimal(15,2);not null" json:"buy_price"`
	SellPrice        decimal.Decimal `gorm:"column:sell_price;type:decimal(15,2);not null" json:"sell_price"`

	// Relations
	Product *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product,omitempty"`
}

// TableName specifies the table name for ProductUnit model
func (ProductUnit) TableName() string {
	return "product_units"
}
