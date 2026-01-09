package utils

import "github.com/gofiber/fiber/v2"

// APIResponse is a standardized JSON response wrapper
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse returns a success response with data
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// CreatedResponse returns a 201 created response
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse returns an error response with custom status code
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status:  "error",
		Message: message,
	})
}

// BadRequestResponse returns a 400 bad request response
func BadRequestResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message)
}

// UnauthorizedResponse returns a 401 unauthorized response
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message)
}

// ForbiddenResponse returns a 403 forbidden response
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message)
}

// NotFoundResponse returns a 404 not found response
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message)
}

// InternalErrorResponse returns a 500 internal server error response
func InternalErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message)
}
