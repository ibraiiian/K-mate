package repository

import (
	"k-mate-backend/config"
	"k-mate-backend/internal/models"
)

// AuditLogRepository handles database operations for AuditLog model
type AuditLogRepository struct{}

// NewAuditLogRepository creates a new AuditLogRepository instance
func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{}
}

// FindAllPaginated retrieves audit logs with pagination and filters
func (r *AuditLogRepository) FindAllPaginated(page, limit int, userID *int, action string) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := config.DB.Model(&models.AuditLog{})

	// Apply filters
	if userID != nil && *userID > 0 {
		query = query.Where("user_id = ?", *userID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginate and fetch
	offset := (page - 1) * limit
	err := query.Preload("User").
		Offset(offset).
		Limit(limit).
		Order("timestamp DESC").
		Find(&logs).Error

	return logs, total, err
}

// FindByID retrieves an audit log by ID
func (r *AuditLogRepository) FindByID(id int) (*models.AuditLog, error) {
	var log models.AuditLog
	err := config.DB.Preload("User").First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetDistinctActions returns all distinct action types
func (r *AuditLogRepository) GetDistinctActions() ([]string, error) {
	var actions []string
	err := config.DB.Model(&models.AuditLog{}).
		Distinct("action").
		Pluck("action", &actions).Error
	return actions, err
}
