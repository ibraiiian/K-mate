package config

import (
	"log"
	"os"
	"time"

	"k-mate-backend/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	// 1. Ambil DSN dari .env
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Println("⚠️ PERINGATAN: DB_DSN kosong! Mengabaikan koneksi database agar server tetap bisa menyala. Harap set DB_DSN di Cloud Run Dashboard.")
		return
	}

	// 2. Config Log Mode (Info query SQL terlihat jika development)
	logMode := logger.Silent
	if os.Getenv("APP_ENV") == "development" {
		logMode = logger.Info
	}

	// 3. Buka Koneksi ke MySQL
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		log.Fatalf("Gagal connect ke Database: %v", err)
	}

	// 4. Set Connection Pool (Agar performa cepat)
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Gagal mengambil instance SQL DB")
	}

	// SetMaxIdleConns: Jumlah koneksi nganggur maksimal
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns: Jumlah koneksi terbuka maksimal (concurrent users)
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime: Durasi koneksi boleh hidup
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Berhasil terhubung ke Database MySQL (Docker)!")

	// Auto Migrate - Create tables if not exist
	database.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ProductUnit{},
		&models.Transaction{},
		&models.TransactionDetail{},
		&models.Expense{},
		&models.Return{},
		&models.ReturnDetail{},
		&models.AuditLog{},
	)
	log.Println("✅ Auto-migrate models selesai!")

	DB = database
}
