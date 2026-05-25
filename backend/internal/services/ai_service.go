package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"k-mate-backend/internal/repository"
)

type AIService struct {
	apiKey      string
	productRepo *repository.ProductRepository
	reportRepo  *repository.ReportRepository
}

func NewAIService(productRepo *repository.ProductRepository, reportRepo *repository.ReportRepository) (*AIService, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY not set")
	}

	return &AIService{
		apiKey:      apiKey,
		productRepo: productRepo,
		reportRepo:  reportRepo,
	}, nil
}

// Groq API types (OpenAI-compatible)
type ChatMessage struct {
	Role       string     `json:"role"`
	
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Tools       []Tool        `json:"tools,omitempty"`
	Temperature float32       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type Tool struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

// Define available tools
func getTools() []Tool {
	return []Tool{
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "get_product_stock",
				Description: "Cek stok barang berdasarkan nama produk. Gunakan ini ketika user bertanya tentang ketersediaan atau jumlah stok suatu barang.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_name": map[string]interface{}{
							"type":        "string",
							"description": "Nama produk yang ingin dicek stoknya, contoh: semen, paku, cat",
						},
					},
					"required": []string{"product_name"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "get_daily_sales",
				Description: "Mendapatkan total omzet/penjualan hari ini. Gunakan ini ketika user bertanya tentang pendapatan atau omzet hari ini.",
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "get_low_stock_products",
				Description: "Mendapatkan daftar produk yang stoknya menipis atau di bawah batas minimum. Gunakan ini ketika user bertanya barang apa yang perlu direstock atau hampir habis.",
			},
		},
	}
}

func (s *AIService) Chat(ctx context.Context, userMessage string, history []ChatMessage) (string, []ChatMessage, error) {
	// Build messages
	messages := []ChatMessage{
		{
			Role:    "system",
			Content: "Kamu adalah K-Mate Assistant, asisten pintar untuk toko bangunan. Kamu membantu pemilik toko mengecek stok barang, melihat omzet penjualan, dan memberikan insight bisnis. Jawab dengan bahasa Indonesia yang ramah dan ringkas. Jika user bertanya tentang stok atau penjualan, gunakan tools yang tersedia.",
		},
	}
	messages = append(messages, history...)
	messages = append(messages, ChatMessage{Role: "user", Content: userMessage})

	// Make API call with tool support
	response, err := s.callGroqAPI(messages)
	if err != nil {
		return "", nil, err
	}

	// Check if model wants to call a tool
	if len(response.Choices) > 0 && len(response.Choices[0].Message.ToolCalls) > 0 {
		// Execute tools and get results
		assistantMsg := response.Choices[0].Message
		messages = append(messages, assistantMsg)

		for _, toolCall := range assistantMsg.ToolCalls {
			result, err := s.executeFunction(toolCall.Function)
			if err != nil {
				result = fmt.Sprintf("Error: %s", err.Error())
			}

			// Add tool result to messages
			messages = append(messages, ChatMessage{
				Role:       "tool",
				Content:    result,
				ToolCallID: toolCall.ID,
			})
		}

		// Call API again with tool results
		response, err = s.callGroqAPI(messages)
		if err != nil {
			return "", nil, err
		}
	}

	// Extract final response
	if len(response.Choices) > 0 {
		finalMsg := response.Choices[0].Message
		messages = append(messages, finalMsg)

		// Return only user and assistant messages for history (exclude system)
		var cleanHistory []ChatMessage
		for _, m := range messages {
			if m.Role == "user" || m.Role == "assistant" {
				cleanHistory = append(cleanHistory, m)
			}
		}

		return finalMsg.Content, cleanHistory, nil
	}

	return "Maaf, saya tidak bisa menjawab saat ini.", nil, nil
}

func (s *AIService) callGroqAPI(messages []ChatMessage) (*ChatResponse, error) {
	reqBody := ChatRequest{
		Model:       "llama-3.3-70b-versatile",
		Messages:    messages,
		Tools:       getTools(),
		Temperature: 0.5,
		MaxTokens:   1024,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Groq API error: %s", string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, err
	}

	return &chatResp, nil
}

func (s *AIService) executeFunction(call FunctionCall) (string, error) {
	log.Printf("Executing tool: %s with args: %s", call.Name, call.Arguments)

	switch call.Name {
	case "get_product_stock":
		var args struct {
			ProductName string `json:"product_name"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", err
		}

		products, err := s.productRepo.FindByName(args.ProductName)
		if err != nil {
			return "", err
		}
		if len(products) == 0 {
			return "Produk tidak ditemukan dengan nama tersebut.", nil
		}

		var result string
		for _, p := range products {
			result += fmt.Sprintf("- %s: %d %s tersedia\n", p.Name, p.StockQty, p.BaseUnit)
		}
		return result, nil

	case "get_daily_sales":
		summary, err := s.reportRepo.GetDailySummary(time.Now())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Omzet hari ini: Rp %s dari %d transaksi.",
			summary.TotalRevenue.StringFixed(0), summary.TransactionCount), nil

	case "get_low_stock_products":
		products, err := s.reportRepo.GetLowStockProducts()
		if err != nil {
			return "", err
		}
		if len(products) == 0 {
			return "Semua stok aman, tidak ada yang di bawah batas minimum.", nil
		}

		var result string
		for _, p := range products {
			result += fmt.Sprintf("- %s (Sisa: %d %s, Minimum: %d)\n", p.Name, p.StockQty, p.BaseUnit, p.MinStockAlert)
		}
		return result, nil

	default:
		return "", fmt.Errorf("unknown function: %s", call.Name)
	}
}
