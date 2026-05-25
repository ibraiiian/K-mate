package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"k-mate-backend/internal/dto"
	"k-mate-backend/internal/middleware"
	"k-mate-backend/internal/models"
	"k-mate-backend/internal/services"
	"k-mate-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	inventoryService *services.InventoryService
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(inventoryService *services.InventoryService) *ProductHandler {
	return &ProductHandler{
		inventoryService: inventoryService,
	}
}

// CreateProduct handles POST /products
// @Summary Create a new product with units and image
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param image formData file false "Product Image"
// @Param data formData string true "Product Data (JSON)"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	// 1. Handle Image Upload
	file, err := c.FormFile("image")
	var imagePath string

	if err != nil {
		// Log the error for debugging
		fmt.Println("No image uploaded or error:", err.Error())
	}

	if err == nil && file != nil {
		fmt.Println("Image received:", file.Filename, "Size:", file.Size)

		if os.Getenv("GCS_BUCKET_NAME") != "" {
			// Save to Google Cloud Storage
			url, err := utils.UploadToGCS(file)
			if err != nil {
				fmt.Println("Error upload to GCS:", err)
				return utils.InternalErrorResponse(c, "Gagal upload gambar ke Cloud Storage")
			}
			imagePath = url
			fmt.Println("Image uploaded to GCS:", imagePath)
		} else {
			// Fallback: Save to Local Storage
			// Ensure directory exists
			uploadDir := "./uploads/products"
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				return utils.InternalErrorResponse(c, "Gagal membuat folder upload")
			}

			// Generate unique filename
			ext := filepath.Ext(file.Filename)
			filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
			savePath := filepath.Join(uploadDir, filename)

			// Save file
			if err := c.SaveFile(file, savePath); err != nil {
				return utils.InternalErrorResponse(c, "Gagal menyimpan gambar lokal")
			}

			// Set database path (relative URL)
			imagePath = fmt.Sprintf("/uploads/products/%s", filename)
			fmt.Println("Image saved locally to:", imagePath)
		}
	}

	// 2. Parse Data JSON
	dataStr := c.FormValue("data")
	var req dto.CreateProductRequest

	if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
		return utils.BadRequestResponse(c, "Format data JSON tidak valid")
	}

	// Set image path from upload
	req.ImagePath = imagePath

	// Validation
	if req.Name == "" {
		return utils.BadRequestResponse(c, "Nama produk wajib diisi")
	}
	if req.BaseUnit == "" {
		return utils.BadRequestResponse(c, "Base unit wajib diisi")
	}
	if len(req.Units) == 0 {
		return utils.BadRequestResponse(c, "Minimal harus ada 1 varian satuan")
	}

	// 3. Create Product via Service
	product, err := h.inventoryService.CreateProduct(&req)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.CreatedResponse(c, "Produk berhasil ditambahkan", toProductResponse(product))
}

// GetAllProducts handles GET /products
// @Summary Get all products with units
// @Tags Products
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search by name"
// @Success 200 {object} utils.APIResponse
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	// Check for search query
	search := c.Query("search")
	if search != "" {
		products, err := h.inventoryService.SearchProducts(search)
		if err != nil {
			return utils.InternalErrorResponse(c, "Gagal mencari produk")
		}
		return utils.SuccessResponse(c, "Hasil pencarian produk", toProductListResponse(products))
	}

	// Check for pagination
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr != "" || limitStr != "" {
		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)

		products, total, err := h.inventoryService.GetProductsPaginated(page, limit)
		if err != nil {
			return utils.InternalErrorResponse(c, "Gagal mengambil data produk")
		}

		return utils.SuccessResponse(c, "Daftar produk", fiber.Map{
			"products": toProductListResponse(products),
			"total":    total,
			"page":     page,
			"limit":    limit,
		})
	}

	// Get all products
	products, err := h.inventoryService.GetAllProducts()
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data produk")
	}

	return utils.SuccessResponse(c, "Daftar produk", toProductListResponse(products))
}

// GetProductByID handles GET /products/:id
// @Summary Get product by ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "ID produk tidak valid")
	}

	product, err := h.inventoryService.GetProductByID(id)
	if err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Detail produk", toProductResponse(product))
}

// UpdateProduct handles PUT /products/:id
// @Summary Update a product
// @Tags Products
// @Accept json,multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param body body dto.UpdateProductRequest true "Product Data"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "ID produk tidak valid")
	}

	var req dto.UpdateProductRequest

	// Check Content-Type to determine how to parse
	contentType := c.Get("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart/form-data (with image)

		// 1. Handle Image Upload
		file, err := c.FormFile("image")
		if err == nil && file != nil {
			fmt.Println("Update - Image received:", file.Filename, "Size:", file.Size)

			if os.Getenv("GCS_BUCKET_NAME") != "" {
				// Save to Google Cloud Storage
				url, err := utils.UploadToGCS(file)
				if err != nil {
					fmt.Println("Error upload to GCS:", err)
					return utils.InternalErrorResponse(c, "Gagal upload gambar ke Cloud Storage")
				}
				req.ImagePath = url
				fmt.Println("Update - Image uploaded to GCS:", req.ImagePath)
			} else {
				// Fallback: Save to Local Storage
				// Ensure directory exists
				uploadDir := "./uploads/products"
				if err := os.MkdirAll(uploadDir, 0755); err != nil {
					return utils.InternalErrorResponse(c, "Gagal membuat folder upload")
				}

				// Generate unique filename
				ext := filepath.Ext(file.Filename)
				filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
				savePath := filepath.Join(uploadDir, filename)

				// Save file
				if err := c.SaveFile(file, savePath); err != nil {
					return utils.InternalErrorResponse(c, "Gagal menyimpan gambar lokal")
				}

				// Set database path
				req.ImagePath = fmt.Sprintf("/uploads/products/%s", filename)
				fmt.Println("Update - Image saved locally to:", req.ImagePath)
			}
		}

		// 2. Parse JSON data from "data" field
		dataStr := c.FormValue("data")
		if dataStr != "" {
			if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
				return utils.BadRequestResponse(c, "Format data JSON tidak valid")
			}
		}
	} else {
		// Handle regular JSON body
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Format JSON tidak valid")
		}
	}

	// Get old product data for audit
	oldProduct, _ := h.inventoryService.GetProductByID(id)
	var oldData interface{}
	if oldProduct != nil {
		oldData = toProductResponse(oldProduct)
	}

	product, err := h.inventoryService.UpdateProduct(id, &req)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	// Log update activity
	utils.LogActivityAsync(utils.AuditLogParams{
		UserID:      middleware.GetUserID(c),
		Action:      models.ActionUpdate,
		TargetTable: "products",
		TargetID:    id,
		OldData:     oldData,
		NewData:     toProductResponse(product),
		Context:     c,
	})

	return utils.SuccessResponse(c, "Produk berhasil diupdate", toProductResponse(product))
}

// DeleteProduct handles DELETE /products/:id
// @Summary Delete a product
// @Tags Products
// @Param id path int true "Product ID"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "ID produk tidak valid")
	}

	// Get product data before delete for audit
	oldProduct, _ := h.inventoryService.GetProductByID(id)
	var oldData interface{}
	if oldProduct != nil {
		oldData = map[string]interface{}{
			"id":   oldProduct.ID,
			"name": oldProduct.Name,
			"code": oldProduct.ProductCode,
		}
	}

	if err := h.inventoryService.DeleteProduct(id); err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}

	// Log delete activity
	utils.LogActivityAsync(utils.AuditLogParams{
		UserID:      middleware.GetUserID(c),
		Action:      models.ActionDelete,
		TargetTable: "products",
		TargetID:    id,
		OldData:     oldData,
		Context:     c,
	})

	return utils.SuccessResponse(c, "Produk berhasil dihapus", nil)
}

// GetLowStockProducts handles GET /products/low-stock
// @Summary Get products with low stock
// @Tags Products
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Router /products/low-stock [get]
func (h *ProductHandler) GetLowStockProducts(c *fiber.Ctx) error {
	products, err := h.inventoryService.GetLowStockProducts()
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data produk")
	}

	return utils.SuccessResponse(c, "Produk dengan stok rendah", toProductListResponse(products))
}

// ================== HELPER FUNCTIONS ==================

func toProductResponse(p *models.Product) dto.ProductResponse {
	units := make([]dto.ProductUnitResponse, len(p.ProductUnits))
	for i, u := range p.ProductUnits {
		units[i] = dto.ProductUnitResponse{
			ID:               u.ID,
			UnitName:         u.UnitName,
			ConversionFactor: u.ConversionFactor,
			BuyPrice:         u.BuyPrice,
			SellPrice:        u.SellPrice,
		}
	}

	return dto.ProductResponse{
		ID:            p.ID,
		ProductCode:   p.ProductCode,
		Name:          p.Name,
		BaseUnit:      p.BaseUnit,
		StockQty:      p.StockQty,
		MinStockAlert: p.MinStockAlert,
		ImagePath:     p.ImagePath,
		Units:         units,
		CreatedAt:     p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toProductListResponse(products []models.Product) []dto.ProductResponse {
	result := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		result[i] = toProductResponse(&p)
	}
	return result
}
