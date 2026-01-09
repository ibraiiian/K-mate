PROJECT REQUIREMENT DOCUMENT (PRD)

Project Name: K-Mate (Sistem Manajemen Toko Bangunan & POS)

Tech Stack: Vue.js 3, Golang (Fiber/Gin), MySQL

BAB I: PENDAHULUAN & RUANG LINGKUP

1.1 Latar Belakang Masalah

Toko bangunan memiliki karakteristik unik dibanding ritel biasa, yaitu:

Multi-Satuan: Barang dibeli dalam jumlah besar (Dus/Sak) tapi bisa dijual eceran (Pcs/Kg).

Stok Ganda: Kesulitan melacak stok jika satuan tidak terkonversi otomatis.

Buta Finansial: Pemilik sering tahu omzet, tapi tidak tahu laba bersih karena tidak menghitung pengeluaran operasional.

1.2 Tujuan Sistem

Membangun sistem terintegrasi yang mampu menangani konversi satuan otomatis, mempercepat transaksi kasir, dan menyajikan laporan laba/rugi real-time dalam satu dashboard.

1.3 Peran Pengguna (User Roles)

Owner (Super Admin): Akses penuh ke Laporan Keuangan, Laba Rugi, dan Grafik Analitik.

Admin Gudang: Mengelola data barang, restock (barang masuk), dan penyesuaian harga.

Kasir: Hanya akses halaman POS (Point of Sales), Retur Barang, dan Riwayat Transaksi Shift-nya.

BAB II: ARSITEKTUR TEKNOLOGI

Komponen

Teknologi

Keterangan

Frontend

Vue.js 3 (Composition API)

Framework UI utama.

State Management

Pinia

Pengganti Vuex (lebih ringan).

CSS Framework

Tailwind CSS

Untuk styling cepat & responsif.

Backend API

Golang (Fiber/Gin)

Menangani Business Logic & Konkurensi tinggi.

Database

MySQL 8.0

Relational Database Management.

ORM

GORM

Abstraksi query database (Wajib Parameterized).

BAB III: SPESIFIKASI MODUL & FITUR

3.1 Modul Dashboard (Executive Summary)

Kartu Ringkasan: Omzet Hari Ini, Pengeluaran Hari Ini, Laba Bersih (Omzet - HPP - Expense - Retur).

Grafik Analitik: Revenue Trend (Line Chart 7 hari).

Tabel Prioritas: Top Selling Products & Low Stock Alert.

3.2 Modul Inventory (Multi-Unit System)

Logic Satuan (CRITICAL): Stok fisik selalu disimpan dalam Base Unit (Satuan Terkecil).

Contoh: Keramik. Base Unit: Lembar. Stok DB: 400 Lembar.

Konversi Dinamis: Unit Jual Dus (Faktor 6) -> Saat terjual 1 Dus, stok berkurang 6 Lembar.

Manajemen Harga: Input Harga Beli (HPP) dan Harga Jual per satuan.

3.3 Modul Point of Sales (Kasir)

UI Kasir: Tampilan full-width, fokus kecepatan.

Unit Selector: Pop-up otomatis saat barang dipilih: "Jual per Dus atau Pcs?".

Hold Cart: Menyimpan keranjang sementara.

Cetak Struk: Integrasi print thermal (via Browser Print).

3.4 Modul Keuangan (Expense Tracker)

Input Pengeluaran (Gaji, Listrik, Operasional).

Kategori Biaya untuk laporan Laba Rugi.

3.5 Modul Retur Barang (Returns & Refunds) - NEW

Fitur untuk menangani pengembalian barang.

Validasi Invoice: Retur wajib berdasarkan No Invoice valid.

Restock Otomatis: Barang diretur menambah stok (jika kondisi bagus).

Pengurangan Omzet: Nilai refund otomatis memotong omzet harian.

3.6 Modul Keamanan (Security Enforcers) - NEW

Role-Based Access Control (RBAC): Validasi role strict di Middleware.

Audit Logging: Mencatat siapa melakukan apa (Hapus barang, Edit harga).

BAB IV: DESAIN DATABASE (PENTING)

⚠️ CRITICAL INSTRUCTION FOR AI:
Do NOT use the SQL snippets in this document as the source of truth.
ALWAYS Refer to database/schema.sql for the actual, executable Database Schema.

Tabel Utama yang harus ada (Detail lihat di schema.sql):

users: Login & Role.

products: Master data (Base Unit).

product_units: Konversi satuan & Harga bertingkat.

transactions & transaction_details: Penjualan.

returns & return_details: Retur barang.

expenses: Pengeluaran operasional.

audit_logs: Log keamanan.

BAB V: ALUR SISTEM (LOGIC & FLOW)

5.1 Flowchart: Transaksi & Kalkulasi Stok

Start: Terima Request Checkout.

Loop Item: Ambil conversion_factor.

Kalkulasi: Total Pengurangan = Qty Beli * conversion_factor.

Cek Stok: IF stock < Total Pengurangan -> Error.

Commit: Kurangi Stok, Simpan Transaksi, Simpan HPP saat itu (Snapshot).

5.2 Flowchart: Proses Retur Barang - NEW

Input: Scan Invoice Lama.

Validasi: Cek batas waktu retur.

Update Stok:

IF Good Condition: Stock = Stock + (Qty Retur * conversion_factor).

IF Damaged: Masuk log barang rusak (Stok jual tidak nambah).

Keuangan: Catat Refund di tabel returns.

5.3 Flowchart: Laba Bersih (Net Profit)

Net Profit = Total Sales - Total HPP - Total Expenses - Total Refunds.

BAB VI: IMPLEMENTASI KEAMANAN (CYBER SECURITY) - NEW

Autentikasi:

JWT dengan Short-lived access token.

Password hash menggunakan Bcrypt.

Rate Limiting pada Login (Max 5 fail/menit).

Proteksi Data:

Anti SQL Injection: Wajib GORM Parameterized Queries.

Anti IDOR: User A tidak boleh melihat transaksi User B (kecuali Admin).

Input Validation: conversion_factor > 0, price >= 0.

Logging:

Semua aksi sensitif (Delete/Update) masuk audit_logs.

BAB VII: STRUKTUR FOLDER PROJECT (MANDATORY)

AI Agent WAJIB mengikuti struktur folder ini saat membuat file baru.

1. Backend Structure (Golang - Gin/Fiber + GORM)

k-mate-backend/
├── cmd/
│   └── api/
│       └── main.go           # Entry point aplikasi
├── config/
│   ├── database.go           # Setup koneksi MySQL & GORM
│   └── env.go                # Load variable dari .env
├── internal/
│   ├── models/               # Struct GORM (Sesuai schema.sql)
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── transaction.go
│   │   └── expense.go
│   ├── dto/                  # Data Transfer Object (Req/Res JSON)
│   │   ├── auth_dto.go
│   │   ├── cart_dto.go
│   │   └── report_dto.go
│   ├── handlers/             # Controller (HTTP Layer)
│   │   ├── auth_handler.go
│   │   ├── product_handler.go
│   │   ├── pos_handler.go
│   │   └── report_handler.go
│   ├── routes/               # Definisi Endpoint
│   │   └── api_routes.go
│   ├── services/             # BUSINESS LOGIC (Inventory/Conversion Logic Here)
│   │   ├── auth_service.go
│   │   ├── inventory_service.go
│   │   └── finance_service.go
│   ├── repository/           # Data Access Layer (GORM Queries)
│   │   ├── product_repo.go
│   │   └── transaction_repo.go
│   └── middleware/
│       ├── jwt_auth.go
│       └── cors.go
├── pkg/
│   └── utils/
│       ├── password.go       # Hash utils
│       └── response.go       # JSON Response wrapper
├── uploads/
├── .env
├── go.mod
└── go.sum


2. Frontend Structure (Vue.js 3 + Vite + Pinia)

k-mate-frontend/
├── public/
├── src/
│   ├── assets/
│   │   └── main.css          # Tailwind Import
│   ├── api/                  # Axios Instances & Endpoints
│   │   ├── axiosInstance.js
│   │   ├── authApi.js
│   │   ├── productApi.js
│   │   └── transactionApi.js
│   ├── components/
│   │   ├── common/           # UI Dasar (Button, Input, Modal)
│   │   │   ├── BaseButton.vue
│   │   │   ├── BaseInput.vue
│   │   │   └── ModalDialog.vue
│   │   ├── layout/           # Sidebar, Navbar
│   │   │   ├── Sidebar.vue
│   │   │   ├── Navbar.vue
│   │   │   └── AppLayout.vue
│   │   ├── pos/              # Logic Kasir
│   │   │   ├── ProductCard.vue
│   │   │   ├── CartItem.vue
│   │   │   └── UnitSelector.vue
│   │   ├── dashboard/        # Charts & Stats
│   │   │   ├── StatCard.vue
│   │   │   └── SalesChart.vue
│   │   └── receipt/          # Template Print
│   │       └── ThermalReceipt.vue
│   ├── composables/          # Reusable Hooks
│   │   ├── useCurrency.js
│   │   └── usePrint.js
│   ├── stores/               # Pinia State
│   │   ├── authStore.js
│   │   ├── cartStore.js
│   │   └── productStore.js
│   ├── router/
│   │   └── index.js
│   ├── views/                # Halaman Utama
│   │   ├── LoginView.vue
│   │   ├── DashboardView.vue
│   │   ├── InventoryView.vue
│   │   ├── POSView.vue
│   │   └── FinanceView.vue
│   ├── App.vue
│   └── main.js
├── .env
├── package.json
├── tailwind.config.js
└── vite.config.js
