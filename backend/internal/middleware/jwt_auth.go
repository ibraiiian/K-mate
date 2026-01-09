package middleware

import (
	"strings"

	"k-mate-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Context keys for storing user data
const (
	ContextKeyUserID   = "user_id"
	ContextKeyUsername = "username"
	ContextKeyRole     = "role"
)

// JWTAuth is a middleware that validates JWT tokens from Authorization header
// It extracts user data from the token and stores it in Fiber context
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.UnauthorizedResponse(c, "Token tidak ditemukan. Silakan login terlebih dahulu")
		}

		// 2. Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.UnauthorizedResponse(c, "Format token tidak valid. Gunakan: Bearer <token>")
		}

		// 3. Extract token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return utils.UnauthorizedResponse(c, "Token kosong")
		}

		// 4. Validate token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return utils.UnauthorizedResponse(c, "Token tidak valid atau sudah expired")
		}

		// 5. Store user data in context for later use in handlers
		c.Locals(ContextKeyUserID, claims.UserID)
		c.Locals(ContextKeyUsername, claims.Username)
		c.Locals(ContextKeyRole, claims.Role)

		// 6. Continue to next handler
		return c.Next()
	}
}

// RoleRequired is a middleware that checks if the user has one of the required roles
// Must be used AFTER JWTAuth middleware
func RoleRequired(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get role from context (set by JWTAuth middleware)
		role, ok := c.Locals(ContextKeyRole).(string)
		if !ok || role == "" {
			return utils.UnauthorizedResponse(c, "Data role tidak ditemukan. Pastikan sudah login")
		}

		// Check if user's role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		// User doesn't have permission
		return utils.ForbiddenResponse(c, "Akses ditolak. Anda tidak memiliki izin untuk mengakses resource ini")
	}
}

// OwnerOnly is a shortcut middleware that only allows 'owner' role
func OwnerOnly() fiber.Handler {
	return RoleRequired("owner")
}

// AdminOrOwner is a shortcut middleware that allows 'admin' or 'owner' roles
func AdminOrOwner() fiber.Handler {
	return RoleRequired("owner", "admin")
}

// AllStaff is a shortcut middleware that allows all authenticated staff
func AllStaff() fiber.Handler {
	return RoleRequired("owner", "admin", "cashier")
}

// ================== HELPER FUNCTIONS ==================

// GetUserID retrieves the authenticated user's ID from context
// Returns 0 if not found
func GetUserID(c *fiber.Ctx) int {
	if userID, ok := c.Locals(ContextKeyUserID).(int); ok {
		return userID
	}
	return 0
}

// GetUsername retrieves the authenticated user's username from context
// Returns empty string if not found
func GetUsername(c *fiber.Ctx) string {
	if username, ok := c.Locals(ContextKeyUsername).(string); ok {
		return username
	}
	return ""
}

// GetRole retrieves the authenticated user's role from context
// Returns empty string if not found
func GetRole(c *fiber.Ctx) string {
	if role, ok := c.Locals(ContextKeyRole).(string); ok {
		return role
	}
	return ""
}
