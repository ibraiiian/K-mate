package utils

import (
	"encoding/json"

	"k-mate-backend/config"
	"k-mate-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

// AuditLogParams contains parameters for logging an activity
type AuditLogParams struct {
	UserID      int
	Action      string
	TargetTable string
	TargetID    int
	OldData     interface{}
	NewData     interface{}
	Context     *fiber.Ctx // Optional: to get IP and User Agent
}

// internalAuditLogParams contains processed parameters (without context)
type internalAuditLogParams struct {
	UserID      int
	Action      string
	TargetTable string
	TargetID    int
	OldData     interface{}
	NewData     interface{}
	IPAddress   string
	UserAgent   string
}

// LogActivity logs an activity to the audit_logs table
func LogActivity(params AuditLogParams) error {
	// Extract IP and User Agent from context BEFORE any async operation
	var ipAddress, userAgent string
	if params.Context != nil {
		ipAddress = params.Context.IP()
		userAgent = params.Context.Get("User-Agent")
	}

	return logActivityInternal(internalAuditLogParams{
		UserID:      params.UserID,
		Action:      params.Action,
		TargetTable: params.TargetTable,
		TargetID:    params.TargetID,
		OldData:     params.OldData,
		NewData:     params.NewData,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	})
}

// logActivityInternal handles the actual database insert
func logActivityInternal(params internalAuditLogParams) error {
	// Convert old data to JSON
	var oldValueJSON *models.JSONMap
	if params.OldData != nil {
		jsonBytes, err := json.Marshal(params.OldData)
		if err == nil {
			var jsonMap models.JSONMap
			if json.Unmarshal(jsonBytes, &jsonMap) == nil {
				oldValueJSON = &jsonMap
			}
		}
	}

	// Convert new data to JSON
	var newValueJSON *models.JSONMap
	if params.NewData != nil {
		jsonBytes, err := json.Marshal(params.NewData)
		if err == nil {
			var jsonMap models.JSONMap
			if json.Unmarshal(jsonBytes, &jsonMap) == nil {
				newValueJSON = &jsonMap
			}
		}
	}

	// Prepare optional fields
	var userID *int
	if params.UserID > 0 {
		userID = &params.UserID
	}

	var targetTable *string
	if params.TargetTable != "" {
		targetTable = &params.TargetTable
	}

	var targetID *int
	if params.TargetID > 0 {
		targetID = &params.TargetID
	}

	var ipAddr, ua *string
	if params.IPAddress != "" {
		ipAddr = &params.IPAddress
	}
	if params.UserAgent != "" {
		ua = &params.UserAgent
	}

	// Create audit log entry
	auditLog := models.AuditLog{
		UserID:      userID,
		Action:      params.Action,
		TargetTable: targetTable,
		TargetID:    targetID,
		OldValue:    oldValueJSON,
		NewValue:    newValueJSON,
		IPAddress:   ipAddr,
		UserAgent:   ua,
	}

	// Save to database
	if err := config.DB.Create(&auditLog).Error; err != nil {
		return err
	}

	return nil
}

// LogActivityAsync logs activity asynchronously (fire and forget)
func LogActivityAsync(params AuditLogParams) {
	// IMPORTANT: Extract context data BEFORE starting goroutine
	// because fiber.Ctx is recycled after the request completes
	var ipAddress, userAgent string
	if params.Context != nil {
		ipAddress = params.Context.IP()
		userAgent = params.Context.Get("User-Agent")
	}

	internalParams := internalAuditLogParams{
		UserID:      params.UserID,
		Action:      params.Action,
		TargetTable: params.TargetTable,
		TargetID:    params.TargetID,
		OldData:     params.OldData,
		NewData:     params.NewData,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	go func() {
		_ = logActivityInternal(internalParams)
	}()
}
