-- K-MATE DATABASE SCHEMA (MySQL 8.0)
-- Version: 2.2 (Remove Mandatory Barcode)

-- 1. USERS (Role Management)
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL, -- Must use Bcrypt
    role ENUM('owner', 'admin', 'cashier') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 2. PRODUCTS (Master Data - Base Unit)
-- Revisi: Barcode (SKU) dihapus/dibuat opsional. Fokus pencarian via Nama.
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_code VARCHAR(50) UNIQUE, -- Kode Manual Internal (Opsional, bukan Barcode wajib)
    name VARCHAR(200) NOT NULL,      -- Field utama untuk pencarian kasir
    base_unit VARCHAR(20) NOT NULL,  -- Satuan Dasar (e.g., Pcs, Kg, Lembar)
    stock_qty INT DEFAULT 0,         -- STOK SELALU DALAM BASE UNIT
    min_stock_alert INT DEFAULT 10,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 3. PRODUCT_UNITS (Conversion Logic)
CREATE TABLE product_units (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    unit_name VARCHAR(20) NOT NULL,    -- e.g., "Dus", "Pak", "Truk"
    conversion_factor INT NOT NULL,    -- e.g., 1 Dus = 12 Pcs
    buy_price DECIMAL(15,2) NOT NULL,  -- HPP per satuan ini
    sell_price DECIMAL(15,2) NOT NULL, -- Harga Jual per satuan ini
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- 4. TRANSACTIONS (Sales Header)
CREATE TABLE transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    invoice_no VARCHAR(50) UNIQUE NOT NULL, -- Format: INV-YYYYMMDD-XXXX
    user_id INT NOT NULL,                   -- Kasir
    total_amount DECIMAL(15,2) NOT NULL,    -- Total tagihan ke customer
    payment_method ENUM('cash', 'transfer', 'qris', 'credit') NOT NULL, -- Credit = Kasbon
    transaction_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 5. TRANSACTION_DETAILS (Sales Items & COGS Snapshot)
CREATE TABLE transaction_details (
    id INT AUTO_INCREMENT PRIMARY KEY,
    transaction_id INT NOT NULL,
    product_unit_id INT NOT NULL,
    qty INT NOT NULL,
    price_at_moment DECIMAL(15,2) NOT NULL, -- Harga Jual saat transaksi
    cost_at_moment DECIMAL(15,2) NOT NULL,  -- HPP saat transaksi (PENTING untuk Laba Rugi)
    subtotal DECIMAL(15,2) NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (product_unit_id) REFERENCES product_units(id)
);

-- 6. EXPENSES (Operational Cost)
CREATE TABLE expenses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(100) NOT NULL,
    category ENUM('operational', 'salary', 'restock', 'maintenance', 'other') NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    expense_date DATE NOT NULL,
    description TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 7. RETURNS (Pengembalian Barang)
CREATE TABLE returns (
    id INT AUTO_INCREMENT PRIMARY KEY,
    transaction_id INT NOT NULL,
    user_id INT NOT NULL,
    return_no VARCHAR(50) UNIQUE NOT NULL, -- RET-YYYYMMDD-XXXX
    total_refund DECIMAL(15,2) NOT NULL,
    reason TEXT,
    return_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 8. RETURN_DETAILS (Item Retur)
CREATE TABLE return_details (
    id INT AUTO_INCREMENT PRIMARY KEY,
    return_id INT NOT NULL,
    product_unit_id INT NOT NULL,
    qty INT NOT NULL,
    refund_price_per_unit DECIMAL(15,2) NOT NULL,
    condition_status ENUM('good', 'damaged') NOT NULL, -- good=restock, damaged=buang
    subtotal_refund DECIMAL(15,2) NOT NULL,
    FOREIGN KEY (return_id) REFERENCES returns(id),
    FOREIGN KEY (product_unit_id) REFERENCES product_units(id)
);

-- 9. AUDIT_LOGS (Security & Tracking)
CREATE TABLE audit_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT, 
    action VARCHAR(50) NOT NULL,      
    target_table VARCHAR(50),         
    target_id INT,
    old_value JSON,                   
    new_value JSON,                   
    ip_address VARCHAR(45),
    user_agent TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for Performance
-- Index product_code tetap ada tapi opsional, fokus ke pencarian nama mungkin perlu Fulltext index nanti
CREATE INDEX idx_products_code ON products(product_code);
CREATE INDEX idx_products_name ON products(name); -- Ditambahkan index Nama untuk pencarian cepat
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_expenses_date ON expenses(expense_date);