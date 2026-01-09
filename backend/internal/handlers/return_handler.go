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

// ReturnHandler handles HTTP requests for returns
type ReturnHandler struct {
	returnService *services.ReturnService
}

// NewReturnHandler creates a new ReturnHandler instance
func NewReturnHandler(returnService *services.ReturnService) *ReturnHandler {
	return &ReturnHandler{
		returnService: returnService,
	}
}

// CreateReturn handles POST /returns - create a new return
func (h *ReturnHandler) CreateReturn(c *fiber.Ctx) error {
	// Get user ID from JWT
	userID := middleware.GetUserID(c)

	// Parse request body
	var req dto.CreateReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Format request tidak valid")
	}

	// Validate required fields
	if req.TransactionID == 0 {
		return utils.BadRequestResponse(c, "transaction_id wajib diisi")
	}
	if len(req.Items) == 0 {
		return utils.BadRequestResponse(c, "Minimal harus ada 1 item untuk diretur")
	}

	// Process return
	result, err := h.returnService.ProcessReturn(userID, &req)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	// Log return activity
	utils.LogActivityAsync(utils.AuditLogParams{
		UserID:      userID,
		Action:      models.ActionReturn,
		TargetTable: "returns",
		TargetID:    result.ID,
		NewData: map[string]interface{}{
			"return_no":      result.ReturnNo,
			"transaction_id": req.TransactionID,
			"total_refund":   result.TotalRefund,
			"items_count":    len(result.Items),
		},
		Context: c,
	})

	return utils.CreatedResponse(c, "Retur berhasil diproses", result)
}

// GetAllReturns handles GET /returns - get all returns
func (h *ReturnHandler) GetAllReturns(c *fiber.Ctx) error {
	returns, err := h.returnService.GetAllReturns()
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data return")
	}

	return utils.SuccessResponse(c, "Daftar return", fiber.Map{
		"returns": returns,
		"total":   len(returns),
	})
}

// GetReturnByID handles GET /returns/:id - get return by ID
func (h *ReturnHandler) GetReturnByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "ID return tidak valid")
	}

	returnRecord, err := h.returnService.GetReturnByID(id)
	if err != nil {
		return utils.NotFoundResponse(c, "Return tidak ditemukan")
	}

	return utils.SuccessResponse(c, "Detail return", returnRecord)
}
