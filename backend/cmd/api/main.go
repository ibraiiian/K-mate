package main

import (
	"log"
	"os"

	"k-mate-backend/config"
	"k-mate-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Cek Posisi Direktori Kerja (Debugging)
	cwd, _ := os.Getwd()
	log.Printf("Info: Menjalankan program dari direktori: %s", cwd)

	// 2. Load Environment Variables
	// Prioritas 1: Coba load .env langsung (jika dijalankan dari root 'backend')
	if err := godotenv.Load(".env"); err != nil {
		// Prioritas 2: Coba load dari folder luar (jika dijalankan dari dalam cmd/api)
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("⚠️ Peringatan: Tidak menemukan file .env di root maupun ../../.env")
			log.Println("👉 Pastikan file bernama '.env' (bukan .env.txt) ada di folder backend/")
		} else {
			log.Println("✅ Berhasil load env dari ../../.env")
		}
	} else {
		log.Println("✅ Berhasil load env dari .env")
	}

	// 3. Connect Database
	config.ConnectDB()

	// 4. Init Fiber App
	app := fiber.New(fiber.Config{
		AppName: "K-Mate Backend v2.2",
	})

	// 5. Middleware
	app.Use(logger.New()) // Log setiap request ke terminal
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // URL Frontend Vue (Vite)
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 6. Setup API Routes
	routes.SetupRoutes(app)

	// 7. Static File Serving (Uploads)
	app.Static("/uploads", "./uploads")

	// 8. Health Check Route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "K-Mate Backend is Running 🚀",
			"db":      "Connected to MySQL",
		})
	})

	// 7. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port :" + port)
	log.Fatal(app.Listen(":" + port))
}
