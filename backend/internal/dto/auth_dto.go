package dto

import "k-mate-backend/internal/models"

// ================== REQUEST DTOs ==================

// RegisterRequest represents the JSON body for user registration
type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=owner admin cashier"`
}

// LoginRequest represents the JSON body for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ================== RESPONSE DTOs ==================

// UserResponse represents user data without sensitive fields
type UserResponse struct {
	ID       int             `json:"id"`
	FullName string          `json:"full_name"`
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// RegisterResponse represents the response after successful registration
type RegisterResponse struct {
	User UserResponse `json:"user"`
}

// ================== HELPER FUNCTIONS ==================

// ToUserResponse converts a User model to UserResponse DTO
func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Username: user.Username,
		Role:     user.Role,
	}
}
