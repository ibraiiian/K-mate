package dto

import "github.com/shopspring/decimal"

// ================== REQUEST DTOs ==================

// ProductUnitRequest represents a unit variant in create/update request
type ProductUnitRequest struct {
	UnitName         string          `json:"unit_name" validate:"required"`
	ConversionFactor int             `json:"conversion_factor" validate:"required,gt=0"`
	BuyPrice         decimal.Decimal `json:"buy_price" validate:"required"`
	SellPrice        decimal.Decimal `json:"sell_price" validate:"required"`
}

// CreateProductRequest represents the JSON body for creating a product
type CreateProductRequest struct {
	ProductCode   *string              `json:"product_code"`
	Name          string               `json:"name" validate:"required"`
	BaseUnit      string               `json:"base_unit" validate:"required"`
	StockQty      int                  `json:"stock_qty"`
	MinStockAlert int                  `json:"min_stock_alert"`
	ImagePath     string               `json:"image_path"`
	Units         []ProductUnitRequest `json:"units" validate:"required,min=1,dive"`
}

// UpdateProductRequest represents the JSON body for updating a product
type UpdateProductRequest struct {
	ProductCode   *string             `json:"product_code"`
	Name          string              `json:"name"`
	BaseUnit      string              `json:"base_unit"`
	StockQty      int                 `json:"stock_qty"`
	MinStockAlert int                 `json:"min_stock_alert"`
	ImagePath     string              `json:"image_path"`
	Units         []UpdateUnitRequest `json:"units"`
}

// UpdateUnitRequest represents unit data for update
type UpdateUnitRequest struct {
	ID               int             `json:"id"`
	UnitName         string          `json:"unit_name"`
	ConversionFactor int             `json:"conversion_factor"`
	BuyPrice         decimal.Decimal `json:"buy_price"`
	SellPrice        decimal.Decimal `json:"sell_price"`
}

// ================== RESPONSE DTOs ==================

// ProductUnitResponse represents a unit in response
type ProductUnitResponse struct {
	ID               int             `json:"id"`
	UnitName         string          `json:"unit_name"`
	ConversionFactor int             `json:"conversion_factor"`
	BuyPrice         decimal.Decimal `json:"buy_price"`
	SellPrice        decimal.Decimal `json:"sell_price"`
}

// ProductResponse represents product data in response
type ProductResponse struct {
	ID            int                   `json:"id"`
	ProductCode   *string               `json:"product_code"`
	Name          string                `json:"name"`
	BaseUnit      string                `json:"base_unit"`
	StockQty      int                   `json:"stock_qty"`
	MinStockAlert int                   `json:"min_stock_alert"`
	ImagePath     string                `json:"image_path"`
	Units         []ProductUnitResponse `json:"units"`
	CreatedAt     string                `json:"created_at"`
	UpdatedAt     string                `json:"updated_at"`
}

// ProductListResponse represents paginated product list
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
}
