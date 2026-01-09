package handlers

import (
	"k-mate-backend/internal/dto"
	"k-mate-backend/internal/models"
	"k-mate-backend/internal/services"
	"k-mate-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles POST /auth/register
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "Register Request"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// 1. Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Format JSON tidak valid")
	}

	// 2. Basic validation
	if req.FullName == "" || req.Username == "" || req.Password == "" || req.Role == "" {
		return utils.BadRequestResponse(c, "Semua field wajib diisi (full_name, username, password, role)")
	}

	if len(req.Password) < 6 {
		return utils.BadRequestResponse(c, "Password minimal 6 karakter")
	}

	// Validate role
	validRoles := map[string]bool{"owner": true, "admin": true, "cashier": true}
	if !validRoles[req.Role] {
		return utils.BadRequestResponse(c, "Role tidak valid (owner/admin/cashier)")
	}

	// 3. Call service
	user, err := h.authService.Register(&req)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	// 4. Return response
	return utils.CreatedResponse(c, "Registrasi berhasil", dto.RegisterResponse{
		User: dto.ToUserResponse(user),
	})
}

// Login handles POST /auth/login
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Login Request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// 1. Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Format JSON tidak valid")
	}

	// 2. Basic validation
	if req.Username == "" || req.Password == "" {
		return utils.BadRequestResponse(c, "Username dan password wajib diisi")
	}

	// 3. Call service
	token, user, err := h.authService.Login(&req)
	if err != nil {
		// Log failed login attempt
		utils.LogActivityAsync(utils.AuditLogParams{
			Action:      models.ActionLoginFailed,
			TargetTable: "users",
			NewData:     map[string]interface{}{"username": req.Username, "reason": err.Error()},
			Context:     c,
		})
		return utils.UnauthorizedResponse(c, err.Error())
	}

	// Log successful login
	utils.LogActivityAsync(utils.AuditLogParams{
		UserID:      user.ID,
		Action:      models.ActionLogin,
		TargetTable: "users",
		TargetID:    user.ID,
		NewData:     map[string]interface{}{"username": user.Username, "role": user.Role},
		Context:     c,
	})

	// 4. Return response with token
	return utils.SuccessResponse(c, "Login berhasil", dto.LoginResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	})
}
