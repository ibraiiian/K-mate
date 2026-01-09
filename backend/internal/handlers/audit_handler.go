package handlers

import (
	"strconv"

	"k-mate-backend/internal/repository"
	"k-mate-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AuditHandler handles HTTP requests for audit logs
type AuditHandler struct {
	auditRepo *repository.AuditLogRepository
}

// NewAuditHandler creates a new AuditHandler instance
func NewAuditHandler(auditRepo *repository.AuditLogRepository) *AuditHandler {
	return &AuditHandler{
		auditRepo: auditRepo,
	}
}

// GetAllLogs handles GET /audit-logs - get all audit logs with pagination
func (h *AuditHandler) GetAllLogs(c *fiber.Ctx) error {
	// Parse pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse filters
	var userID *int
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if uid, err := strconv.Atoi(userIDStr); err == nil {
			userID = &uid
		}
	}
	action := c.Query("action")

	// Fetch logs
	logs, total, err := h.auditRepo.FindAllPaginated(page, limit, userID, action)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data audit log")
	}

	// Transform response
	var response []fiber.Map
	for _, log := range logs {
		userName := ""
		if log.User != nil {
			userName = log.User.FullName
		}

		response = append(response, fiber.Map{
			"id":           log.ID,
			"user_id":      log.UserID,
			"user_name":    userName,
			"action":       log.Action,
			"target_table": log.TargetTable,
			"target_id":    log.TargetID,
			"old_value":    log.OldValue,
			"new_value":    log.NewValue,
			"ip_address":   log.IPAddress,
			"timestamp":    log.Timestamp.Format("2006-01-02 15:04:05"),
		})
	}

	return utils.SuccessResponse(c, "Daftar audit log", fiber.Map{
		"logs":  response,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetLogByID handles GET /audit-logs/:id - get audit log detail
func (h *AuditHandler) GetLogByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "ID tidak valid")
	}

	log, err := h.auditRepo.FindByID(id)
	if err != nil {
		return utils.NotFoundResponse(c, "Log tidak ditemukan")
	}

	userName := ""
	if log.User != nil {
		userName = log.User.FullName
	}

	return utils.SuccessResponse(c, "Detail audit log", fiber.Map{
		"id":           log.ID,
		"user_id":      log.UserID,
		"user_name":    userName,
		"action":       log.Action,
		"target_table": log.TargetTable,
		"target_id":    log.TargetID,
		"old_value":    log.OldValue,
		"new_value":    log.NewValue,
		"ip_address":   log.IPAddress,
		"user_agent":   log.UserAgent,
		"timestamp":    log.Timestamp.Format("2006-01-02 15:04:05"),
	})
}

// GetActions handles GET /audit-logs/actions - get distinct action types
func (h *AuditHandler) GetActions(c *fiber.Ctx) error {
	actions, err := h.auditRepo.GetDistinctActions()
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil daftar aksi")
	}

	return utils.SuccessResponse(c, "Daftar aksi", actions)
}
