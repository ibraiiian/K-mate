package services

import (
	"errors"

	"k-mate-backend/internal/dto"
	"k-mate-backend/internal/models"
	"k-mate-backend/internal/repository"

	"github.com/shopspring/decimal"
)

// InventoryService handles inventory business logic
type InventoryService struct {
	productRepo *repository.ProductRepository
}

// NewInventoryService creates a new InventoryService instance
func NewInventoryService(productRepo *repository.ProductRepository) *InventoryService {
	return &InventoryService{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product with its unit variants
func (s *InventoryService) CreateProduct(req *dto.CreateProductRequest) (*models.Product, error) {
	// 1. Validate units
	if len(req.Units) == 0 {
		return nil, errors.New("minimal harus ada 1 satuan/unit")
	}

	for _, unit := range req.Units {
		if unit.ConversionFactor <= 0 {
			return nil, errors.New("conversion_factor harus lebih dari 0")
		}
		if unit.BuyPrice.LessThan(decimal.Zero) {
			return nil, errors.New("buy_price tidak boleh negatif")
		}
		if unit.SellPrice.LessThan(decimal.Zero) {
			return nil, errors.New("sell_price tidak boleh negatif")
		}
	}

	// 2. Create product model
	product := &models.Product{
		ProductCode:   req.ProductCode,
		Name:          req.Name,
		BaseUnit:      req.BaseUnit,
		StockQty:      req.StockQty,
		MinStockAlert: req.MinStockAlert,
		ImagePath:     req.ImagePath,
	}

	// Set default min stock alert
	if product.MinStockAlert == 0 {
		product.MinStockAlert = 10
	}

	// 3. Create product units
	units := make([]models.ProductUnit, len(req.Units))
	for i, u := range req.Units {
		units[i] = models.ProductUnit{
			UnitName:         u.UnitName,
			ConversionFactor: u.ConversionFactor,
			BuyPrice:         u.BuyPrice,
			SellPrice:        u.SellPrice,
		}
	}

	// 4. Save to database with transaction
	if err := s.productRepo.CreateProductWithUnits(product, units); err != nil {
		return nil, errors.New("gagal menyimpan produk: " + err.Error())
	}

	return product, nil
}

// GetAllProducts retrieves all products
func (s *InventoryService) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.FindAll()
}

// GetProductsPaginated retrieves products with pagination
func (s *InventoryService) GetProductsPaginated(page, limit int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.productRepo.FindAllPaginated(page, limit)
}

// GetProductByID retrieves a product by ID
func (s *InventoryService) GetProductByID(id int) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return product, nil
}

// SearchProducts searches products by name
func (s *InventoryService) SearchProducts(name string) ([]models.Product, error) {
	return s.productRepo.FindByName(name)
}

// UpdateProduct updates a product
func (s *InventoryService) UpdateProduct(id int, req *dto.UpdateProductRequest) (*models.Product, error) {
	// 1. Find existing product
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	// 2. Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.BaseUnit != "" {
		product.BaseUnit = req.BaseUnit
	}
	if req.ProductCode != nil {
		product.ProductCode = req.ProductCode
	}
	if req.StockQty >= 0 {
		product.StockQty = req.StockQty
	}
	if req.MinStockAlert > 0 {
		product.MinStockAlert = req.MinStockAlert
	}
	if req.ImagePath != "" {
		product.ImagePath = req.ImagePath
	}

	// 3. Save product changes
	if err := s.productRepo.Update(product); err != nil {
		return nil, errors.New("gagal mengupdate produk: " + err.Error())
	}

	// 4. Update units if provided
	if len(req.Units) > 0 {
		for _, unitReq := range req.Units {
			if unitReq.ID > 0 {
				// Update existing unit
				unit, err := s.productRepo.FindUnitByID(unitReq.ID)
				if err == nil && unit != nil {
					unit.UnitName = unitReq.UnitName
					if unitReq.ConversionFactor > 0 {
						unit.ConversionFactor = unitReq.ConversionFactor
					}
					unit.BuyPrice = unitReq.BuyPrice
					unit.SellPrice = unitReq.SellPrice
					s.productRepo.UpdateUnit(unit)
				}
			}
		}
		// Reload product with updated units
		product, _ = s.productRepo.FindByID(id)
	}

	return product, nil
}

// DeleteProduct deletes a product
func (s *InventoryService) DeleteProduct(id int) error {
	// Check if product exists
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	if err := s.productRepo.Delete(id); err != nil {
		return errors.New("gagal menghapus produk")
	}

	return nil
}

// GetLowStockProducts retrieves products with low stock
func (s *InventoryService) GetLowStockProducts() ([]models.Product, error) {
	return s.productRepo.FindLowStock()
}

// ================== PRODUCT UNIT OPERATIONS ==================

// AddProductUnit adds a new unit variant to a product
func (s *InventoryService) AddProductUnit(productID int, req *dto.ProductUnitRequest) (*models.ProductUnit, error) {
	// Validate
	if req.ConversionFactor <= 0 {
		return nil, errors.New("conversion_factor harus lebih dari 0")
	}

	// Check product exists
	_, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	unit := &models.ProductUnit{
		ProductID:        productID,
		UnitName:         req.UnitName,
		ConversionFactor: req.ConversionFactor,
		BuyPrice:         req.BuyPrice,
		SellPrice:        req.SellPrice,
	}

	if err := s.productRepo.CreateUnit(unit); err != nil {
		return nil, errors.New("gagal menambah unit")
	}

	return unit, nil
}

// UpdateProductUnit updates a product unit
func (s *InventoryService) UpdateProductUnit(unitID int, req *dto.ProductUnitRequest) (*models.ProductUnit, error) {
	unit, err := s.productRepo.FindUnitByID(unitID)
	if err != nil {
		return nil, errors.New("unit tidak ditemukan")
	}

	if req.ConversionFactor > 0 {
		unit.ConversionFactor = req.ConversionFactor
	}
	if req.UnitName != "" {
		unit.UnitName = req.UnitName
	}
	if !req.BuyPrice.IsZero() {
		unit.BuyPrice = req.BuyPrice
	}
	if !req.SellPrice.IsZero() {
		unit.SellPrice = req.SellPrice
	}

	if err := s.productRepo.UpdateUnit(unit); err != nil {
		return nil, errors.New("gagal mengupdate unit")
	}

	return unit, nil
}

// DeleteProductUnit deletes a product unit
func (s *InventoryService) DeleteProductUnit(unitID int) error {
	_, err := s.productRepo.FindUnitByID(unitID)
	if err != nil {
		return errors.New("unit tidak ditemukan")
	}

	return s.productRepo.DeleteUnit(unitID)
}
