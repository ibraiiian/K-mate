package repository

import (
	"k-mate-backend/config"
	"k-mate-backend/internal/models"

	"gorm.io/gorm"
)

// ProductRepository handles database operations for Product and ProductUnit models
type ProductRepository struct{}

// NewProductRepository creates a new ProductRepository instance
func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// CreateProductWithUnits creates a product with its unit variants in a transaction
func (r *ProductRepository) CreateProductWithUnits(product *models.Product, units []models.ProductUnit) error {
	// Start database transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Defer rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create product
	if err := tx.Create(product).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Set product_id for each unit and create them
	for i := range units {
		units[i].ProductID = product.ID
	}

	if len(units) > 0 {
		if err := tx.Create(&units).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 3. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 4. Load units into product for response
	product.ProductUnits = units

	return nil
}

// FindAll retrieves all products with their units
func (r *ProductRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := config.DB.Preload("ProductUnits").Find(&products).Error
	return products, err
}

// FindAllPaginated retrieves products with pagination
func (r *ProductRepository) FindAllPaginated(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// Count total
	if err := config.DB.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := config.DB.Preload("ProductUnits").
		Offset(offset).
		Limit(limit).
		Order("id DESC").
		Find(&products).Error

	return products, total, err
}

// FindByID retrieves a product by ID with its units
func (r *ProductRepository) FindByID(id int) (*models.Product, error) {
	var product models.Product
	err := config.DB.Preload("ProductUnits").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindByName searches products by name (partial match)
func (r *ProductRepository) FindByName(name string) ([]models.Product, error) {
	var products []models.Product
	err := config.DB.Preload("ProductUnits").
		Where("name LIKE ?", "%"+name+"%").
		Find(&products).Error
	return products, err
}

// Update updates a product (only specific fields to avoid created_at issues)
func (r *ProductRepository) Update(product *models.Product) error {
	return config.DB.Model(&models.Product{}).
		Where("id = ?", product.ID).
		Updates(map[string]interface{}{
			"product_code":    product.ProductCode,
			"name":            product.Name,
			"base_unit":       product.BaseUnit,
			"stock_qty":       product.StockQty,
			"min_stock_alert": product.MinStockAlert,
			"image_path":      product.ImagePath,
		}).Error
}

// UpdateStock updates product stock quantity
func (r *ProductRepository) UpdateStock(productID int, newStock int) error {
	return config.DB.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("stock_qty", newStock).Error
}

// AdjustStock adjusts stock by delta (positive or negative)
func (r *ProductRepository) AdjustStock(tx *gorm.DB, productID int, delta int) error {
	return tx.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("stock_qty", gorm.Expr("stock_qty + ?", delta)).Error
}

// Delete deletes a product (cascade deletes units)
func (r *ProductRepository) Delete(id int) error {
	return config.DB.Delete(&models.Product{}, id).Error
}

// CountAll returns total number of products
func (r *ProductRepository) CountAll() (int64, error) {
	var count int64
	err := config.DB.Model(&models.Product{}).Count(&count).Error
	return count, err
}

// FindLowStock finds products with stock below minimum alert
func (r *ProductRepository) FindLowStock() ([]models.Product, error) {
	var products []models.Product
	err := config.DB.Preload("ProductUnits").
		Where("stock_qty <= min_stock_alert").
		Find(&products).Error
	return products, err
}

// ================== PRODUCT UNIT REPOSITORY ==================

// FindUnitByID retrieves a product unit by ID
func (r *ProductRepository) FindUnitByID(id int) (*models.ProductUnit, error) {
	var unit models.ProductUnit
	err := config.DB.Preload("Product").First(&unit, id).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// CreateUnit creates a new product unit
func (r *ProductRepository) CreateUnit(unit *models.ProductUnit) error {
	return config.DB.Create(unit).Error
}

// UpdateUnit updates a product unit
func (r *ProductRepository) UpdateUnit(unit *models.ProductUnit) error {
	return config.DB.Save(unit).Error
}

// DeleteUnit deletes a product unit
func (r *ProductRepository) DeleteUnit(id int) error {
	return config.DB.Delete(&models.ProductUnit{}, id).Error
}
