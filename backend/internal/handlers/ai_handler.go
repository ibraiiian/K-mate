package handlers

import (
	"k-mate-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AIHandler struct {
	service *services.AIService
}

func NewAIHandler(service *services.AIService) *AIHandler {
	return &AIHandler{service: service}
}

type ChatRequest struct {
	Message string                 `json:"message"`
	History []services.ChatMessage `json:"history"`
}

type ChatResponse struct {
	Response string                 `json:"response"`
	History  []services.ChatMessage `json:"history"`
}

func (h *AIHandler) Chat(c *fiber.Ctx) error {
	// Check if AI service is available
	if h.service == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "error",
			"message": "AI Service tidak tersedia. Pastikan GROQ_API_KEY sudah dikonfigurasi.",
		})
	}

	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Message cannot be empty",
		})
	}

	response, newHistory, err := h.service.Chat(c.Context(), req.Message, req.History)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": ChatResponse{
			Response: response,
			History:  newHistory,
		},
	})
}
