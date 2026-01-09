package handlers

import (
	"strconv"

	"k-mate-backend/internal/dto"
	"k-mate-backend/internal/middleware"
	"k-mate-backend/internal/models"
	"k-mate-backend/internal/services"
	"k-mate-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ExpenseHandler handles expense HTTP requests
type ExpenseHandler struct {
	expenseService *services.ExpenseService
}

// NewExpenseHandler creates a new ExpenseHandler instance
func NewExpenseHandler(expenseService *services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

// CreateExpense handles POST /expenses
func (h *ExpenseHandler) CreateExpense(c *fiber.Ctx) error {
	var req dto.CreateExpenseRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Format JSON tidak valid")
	}

	// Validation
	if req.Title == "" {
		return utils.BadRequestResponse(c, "Judul pengeluaran wajib diisi")
	}
	if req.Category == "" {
		return utils.BadRequestResponse(c, "Kategori wajib diisi")
	}
	if req.ExpenseDate == "" {
		return utils.BadRequestResponse(c, "Tanggal pengeluaran wajib diisi")
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.UnauthorizedResponse(c, "User tidak terautentikasi")
	}

	expense, err := h.expenseService.CreateExpense(userID, &req)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.CreatedResponse(c, "Pengeluaran berhasil dicatat", toExpenseResponse(expense))
}

// GetAllExpenses handles GET /expenses
func (h *ExpenseHandler) GetAllExpenses(c *fiber.Ctx) error {
	// Check for category filter
	category := c.Query("category")
	if category != "" {
		expenses, err := h.expenseService.GetExpensesByCategory(category)
		if err != nil {
			return utils.InternalErrorResponse(c, "Gagal mengambil data pengeluaran")
		}
		return utils.SuccessResponse(c, "Daftar pengeluaran", toExpenseListResponse(expenses))
	}

	// Check for pagination
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr != "" || limitStr != "" {
		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)

		expenses, total, err := h.expenseService.GetExpensesPaginated(page, limit)
		if err != nil {
			return utils.InternalErrorResponse(c, "Gagal mengambil data pengeluaran")
		}

		return utils.SuccessResponse(c, "Daftar pengeluaran", fiber.Map{
			"expenses": toExpenseListResponse(expenses),
			"total":    total,
			"page":     page,
			"limit":    limit,
		})
	}

	expenses, err := h.expenseService.GetAllExpenses()
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data pengeluaran")
	}

	return utils.SuccessResponse(c, "Daftar pengeluaran", toExpenseListResponse(expenses))
}

// GetExpenseByID handles GET /expenses/:id
func (h *ExpenseHandler) GetExpenseByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "ID tidak valid")
	}

	expense, err := h.expenseService.GetExpenseByID(id)
	if err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Detail pengeluaran", toExpenseResponse(expense))
}

// GetTodayExpenses handles GET /expenses/today
func (h *ExpenseHandler) GetTodayExpenses(c *fiber.Ctx) error {
	expenses, err := h.expenseService.GetTodayExpenses()
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data pengeluaran")
	}

	return utils.SuccessResponse(c, "Pengeluaran hari ini", toExpenseListResponse(expenses))
}

// UpdateExpense handles PUT /expenses/:id
func (h *ExpenseHandler) UpdateExpense(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "ID tidak valid")
	}

	var req dto.UpdateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Format JSON tidak valid")
	}

	expense, err := h.expenseService.UpdateExpense(id, &req)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Pengeluaran berhasil diupdate", toExpenseResponse(expense))
}

// DeleteExpense handles DELETE /expenses/:id
func (h *ExpenseHandler) DeleteExpense(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "ID tidak valid")
	}

	if err := h.expenseService.DeleteExpense(id); err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Pengeluaran berhasil dihapus", nil)
}

// ================== HELPER FUNCTIONS ==================

func toExpenseResponse(e *models.Expense) dto.ExpenseResponse {
	userName := ""
	if e.User != nil {
		userName = e.User.FullName
	}

	return dto.ExpenseResponse{
		ID:          e.ID,
		UserID:      e.UserID,
		UserName:    userName,
		Title:       e.Title,
		Category:    string(e.Category),
		Amount:      e.Amount,
		ExpenseDate: e.ExpenseDate.Format("2006-01-02"),
		Description: e.Description,
		CreatedAt:   e.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toExpenseListResponse(expenses []models.Expense) []dto.ExpenseResponse {
	result := make([]dto.ExpenseResponse, len(expenses))
	for i, e := range expenses {
		result[i] = toExpenseResponse(&e)
	}
	return result
}
