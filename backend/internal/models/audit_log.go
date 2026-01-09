package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONMap represents a JSON object type for GORM
type JSONMap map[string]interface{}

// Value implements the driver.Valuer interface for JSONMap
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONMap
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// AuditLog represents the audit_logs table (Security & Tracking)
type AuditLog struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      *int      `gorm:"column:user_id;index" json:"user_id"`
	Action      string    `gorm:"column:action;type:varchar(50);not null" json:"action"`
	TargetTable *string   `gorm:"column:target_table;type:varchar(50)" json:"target_table"`
	TargetID    *int      `gorm:"column:target_id" json:"target_id"`
	OldValue    *JSONMap  `gorm:"column:old_value;type:json" json:"old_value"`
	NewValue    *JSONMap  `gorm:"column:new_value;type:json" json:"new_value"`
	IPAddress   *string   `gorm:"column:ip_address;type:varchar(45)" json:"ip_address"`
	UserAgent   *string   `gorm:"column:user_agent;type:text" json:"user_agent"`
	Timestamp   time.Time `gorm:"column:timestamp;autoCreateTime" json:"timestamp"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}

// Common action constants
const (
	ActionLogin       = "LOGIN"
	ActionLoginFailed = "LOGIN_FAILED"
	ActionLogout      = "LOGOUT"
	ActionCreate      = "CREATE"
	ActionUpdate      = "UPDATE"
	ActionDelete      = "DELETE"
	ActionCheckout    = "CHECKOUT"
	ActionReturn      = "RETURN"
	ActionPriceChange = "PRICE_CHANGE"
)
