# K-Mate 🏪

**Sistem Kasir & Manajemen Toko Bangunan Modern**

K-Mate adalah aplikasi Point of Sale (POS) lengkap untuk toko bangunan dengan fitur manajemen inventaris, pencatatan transaksi, laporan keuangan, dan **AI Assistant** yang membantu pemilik toko mendapatkan insight bisnis secara real-time.

![Vue.js](https://img.shields.io/badge/Vue.js-3.x-4FC08D?logo=vue.js)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white)
![TailwindCSS](https://img.shields.io/badge/Tailwind-3.x-38B2AC?logo=tailwind-css)

---

## ✨ Fitur Utama

### 🛒 Point of Sale (POS)
- Checkout cepat dengan pencarian produk
- Multi-unit support (eceran, pack, karung, dll)
- Cetak struk thermal (58mm/80mm)
- Riwayat transaksi lengkap

### 📦 Manajemen Inventaris
- CRUD produk dengan multi-satuan
- Notifikasi stok menipis
- Konversi unit otomatis

### 💰 Laporan Keuangan
- Dashboard real-time (omzet, HPP, laba)
- Laporan periode custom
- Grafik penjualan

### 🔄 Retur Barang
- Proses retur dari transaksi
- Pengembalian stok otomatis
- Pencatatan expense refund

### 🤖 AI Assistant (K-Mate Assistant)
- Cek stok barang via chat
- Lihat omzet hari ini
- Rekomendasi restock
- Powered by Groq LLaMA 3.3

### 📋 Audit Log
- Tracking semua aktivitas user
- Filter by action type
- Role-based access

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| **Frontend** | Vue 3, Vite, Tailwind CSS, Pinia, Axios |
| **Backend** | Go, Fiber, GORM |
| **Database** | MySQL 8.0 |
| **AI** | Groq API (LLaMA 3.3 70B) |

---

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- MySQL 8.0
- Groq API Key (untuk AI features)

### 1. Clone Repository
```bash
git clone https://github.com/ibraiiian/K-mate.git
cd K-mate
```

### 2. Setup Backend
```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env dengan kredensial database Anda
# DB_DSN=user:password@tcp(localhost:3306)/k_mate_db?charset=utf8mb4&parseTime=True&loc=Local
# GROQ_API_KEY=your_groq_api_key

# Install dependencies
go mod tidy

# Run server
go run ./cmd/api
```

### 3. Setup Frontend
```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

### 4. Access Application
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:3000

---

## 📁 Project Structure

```
K-Mate/
├── backend/
│   ├── cmd/api/         # Entry point
│   ├── config/          # Database config
│   ├── internal/
│   │   ├── handlers/    # HTTP handlers
│   │   ├── middleware/  # Auth, CORS
│   │   ├── models/      # GORM models
│   │   ├── repository/  # Database queries
│   │   ├── routes/      # API routes
│   │   └── services/    # Business logic
│   └── pkg/utils/       # Helper functions
│
└── frontend/
    ├── src/
    │   ├── api/         # Axios instance
    │   ├── components/  # Vue components
    │   ├── stores/      # Pinia stores
    │   ├── views/       # Page views
    │   └── router/      # Vue Router
    └── public/
```

---

## 🔐 Default Credentials

| Role | Username | Password |
|------|----------|----------|
| Owner | `owner` | `owner123` |
| Admin | `admin` | `admin123` |
| Kasir | `kasir` | `kasir123` |

---

## 📝 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Login |
| GET | `/api/v1/products` | Get all products |
| POST | `/api/v1/pos/checkout` | Create transaction |
| GET | `/api/v1/dashboard/summary` | Get daily summary |
| POST | `/api/v1/chat` | AI Chat |
| GET | `/api/v1/audit-logs` | Get audit logs |

---

## 👨‍💻 Author

**Ibrahim**

---

## 📄 License

This project is licensed under the MIT License.
