package routes

import (
	"k-mate-backend/internal/handlers"
	"k-mate-backend/internal/middleware"
	"k-mate-backend/internal/repository"
	"k-mate-backend/internal/services"
	"k-mate-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App) {
	// ==================== DEPENDENCY INJECTION ====================
	// Auth
	userRepo := repository.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Inventory
	productRepo := repository.NewProductRepository()
	inventoryService := services.NewInventoryService(productRepo)
	productHandler := handlers.NewProductHandler(inventoryService)

	// POS
	transactionRepo := repository.NewTransactionRepository()
	posService := services.NewPOSService(transactionRepo, productRepo)
	posHandler := handlers.NewPOSHandler(posService)

	// Expense
	expenseRepo := repository.NewExpenseRepository()
	expenseService := services.NewExpenseService(expenseRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseService)

	// Report/Dashboard
	reportRepo := repository.NewReportRepository()
	reportHandler := handlers.NewReportHandler(reportRepo)

	// Return
	returnRepo := repository.NewReturnRepository()
	returnService := services.NewReturnService(returnRepo, transactionRepo, productRepo)
	returnHandler := handlers.NewReturnHandler(returnService)

	// Audit Log
	auditRepo := repository.NewAuditLogRepository()
	auditHandler := handlers.NewAuditHandler(auditRepo)

	// AI Service
	aiService, err := services.NewAIService(productRepo, reportRepo)
	if err != nil {
		// Log error but allow app to start (AI features won't work)
		// in real app use proper logger
		println("Warning: AI Service failed to initialize:", err.Error())
	}
	aiHandler := handlers.NewAIHandler(aiService)

	// API v1 Group
	api := app.Group("/api/v1")

	// ==================== PUBLIC ROUTES ====================
	// Auth Routes - No JWT Required
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// ==================== PROTECTED ROUTES ====================
	// All routes below require valid JWT token

	// User Profile - All authenticated users can access
	protected := api.Group("/", middleware.JWTAuth())
	protected.Get("/me", func(c *fiber.Ctx) error {
		return utils.SuccessResponse(c, "Data user berhasil diambil", fiber.Map{
			"user_id":  middleware.GetUserID(c),
			"username": middleware.GetUsername(c),
			"role":     middleware.GetRole(c),
		})
	})

	// ==================== PRODUCT ROUTES ====================
	products := api.Group("/products", middleware.JWTAuth())

	// GET - All staff can access (owner, admin, cashier)
	products.Get("/", productHandler.GetAllProducts)
	products.Get("/low-stock", middleware.AdminOrOwner(), productHandler.GetLowStockProducts)
	products.Get("/:id", productHandler.GetProductByID)

	// POST, PUT, DELETE - Only admin and owner can access
	products.Post("/", middleware.AdminOrOwner(), productHandler.CreateProduct)
	products.Put("/:id", middleware.AdminOrOwner(), productHandler.UpdateProduct)
	products.Delete("/:id", middleware.AdminOrOwner(), productHandler.DeleteProduct)

	// ==================== POS ROUTES ====================
	pos := api.Group("/pos", middleware.JWTAuth())

	// Checkout - All staff can access
	pos.Post("/checkout", posHandler.Checkout)

	// Transaction History
	pos.Get("/transactions", posHandler.GetAllTransactions)
	pos.Get("/transactions/today", posHandler.GetTodayTransactions)
	pos.Get("/transactions/my", posHandler.GetMyTransactions)
	pos.Get("/transactions/invoice/:invoiceNo", posHandler.GetTransactionByInvoice)
	pos.Get("/transactions/:id", posHandler.GetTransactionByID)

	// Summary - Admin/Owner only
	pos.Get("/summary/today", middleware.AdminOrOwner(), posHandler.GetTodaySummary)

	// ==================== EXPENSE ROUTES ====================
	expenses := api.Group("/expenses", middleware.JWTAuth(), middleware.AdminOrOwner())

	expenses.Post("/", expenseHandler.CreateExpense)
	expenses.Get("/", expenseHandler.GetAllExpenses)
	expenses.Get("/today", expenseHandler.GetTodayExpenses)
	expenses.Get("/:id", expenseHandler.GetExpenseByID)
	expenses.Put("/:id", expenseHandler.UpdateExpense)
	expenses.Delete("/:id", expenseHandler.DeleteExpense)

	// ==================== DASHBOARD ROUTES (OWNER ONLY) ====================
	dashboard := api.Group("/dashboard", middleware.JWTAuth(), middleware.OwnerOnly())

	dashboard.Get("/summary", reportHandler.GetDailySummary)
	dashboard.Get("/chart", reportHandler.GetSalesChart)
	dashboard.Get("/top-products", reportHandler.GetTopProducts)
	dashboard.Get("/low-stock", reportHandler.GetLowStockAlert)

	// ==================== REPORT ROUTES (ADMIN & OWNER) ====================
	reports := api.Group("/reports", middleware.JWTAuth(), middleware.AdminOrOwner())
	reports.Get("/", reportHandler.GetFinancialReport)

	// ==================== RETURN ROUTES ====================
	returns := api.Group("/returns", middleware.JWTAuth())

	// POST - All staff can create returns
	returns.Post("/", returnHandler.CreateReturn)

	// GET - Admin/Owner can view all returns
	returns.Get("/", middleware.AdminOrOwner(), returnHandler.GetAllReturns)
	returns.Get("/:id", middleware.AdminOrOwner(), returnHandler.GetReturnByID)

	// ==================== AUDIT LOG ROUTES (OWNER ONLY) ====================
	auditLogs := api.Group("/audit-logs", middleware.JWTAuth(), middleware.OwnerOnly())

	auditLogs.Get("/actions", auditHandler.GetActions)
	auditLogs.Get("/", auditHandler.GetAllLogs)
	auditLogs.Get("/:id", auditHandler.GetLogByID)

	// ==================== AI CHAT ROUTES ====================
	// Available to all authenticated users
	api.Post("/chat", middleware.JWTAuth(), aiHandler.Chat)
}
