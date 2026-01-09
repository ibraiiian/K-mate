package models

import "time"

// UserRole represents the allowed roles for a user
type UserRole string

const (
	UserRoleOwner   UserRole = "owner"
	UserRoleAdmin   UserRole = "admin"
	UserRoleCashier UserRole = "cashier"
)

// User represents the users table in the database
type User struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName     string    `gorm:"column:full_name;type:varchar(100);not null" json:"full_name"`
	Username     string    `gorm:"column:username;type:varchar(50);uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Role         UserRole  `gorm:"column:role;type:enum('owner','admin','cashier');not null" json:"role"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
