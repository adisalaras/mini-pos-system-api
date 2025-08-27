-- ================================================
-- Mini POS Database Migration Script
-- ================================================

-- Create database (run this first)
-- CREATE DATABASE minipos;
-- \c minipos;

-- ================================================
-- 1. CREATE TABLES
-- ================================================

-- tabel produk
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(15,2) NOT NULL CHECK (price >= 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

-- tabel transaksi
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(15,2) NOT NULL CHECK (total_amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

-- tabel item_transaksi
CREATE TABLE transaction_items (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    product_name VARCHAR(100) NOT NULL,
    price DECIMAL(15,2) NOT NULL CHECK (price >= 0),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    subtotal DECIMAL(15,2) NOT NULL CHECK (subtotal >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transaction_items_transaction_id 
        FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    CONSTRAINT fk_transaction_items_product_id 
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

-- ================================================
-- 2. CREATE INDEXES
-- ================================================

-- Index untuk produk
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_deleted_at ON products(deleted_at);
CREATE INDEX idx_products_created_at ON products(created_at);

-- Index untuk transaksi
CREATE INDEX idx_transactions_transaction_date ON transactions(transaction_date);
CREATE INDEX idx_transactions_deleted_at ON transactions(deleted_at);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);

-- Index untuk item_transaksi
CREATE INDEX idx_transaction_items_transaction_id ON transaction_items(transaction_id);
CREATE INDEX idx_transaction_items_product_id ON transaction_items(product_id);
CREATE INDEX idx_transaction_items_created_at ON transaction_items(created_at);

-- ================================================
-- 3. CREATE TRIGGERS FOR UPDATED_AT
-- ================================================

-- Function untuk mengupdate kolom updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for produk
CREATE TRIGGER trigger_produk_updated_at
    BEFORE UPDATE ON produk
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Triggers untuk transaksi
CREATE TRIGGER trigger_transaksi_updated_at
    BEFORE UPDATE ON transaksi
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Triggers for item_transaksi
CREATE TRIGGER trigger_item_transaksi_updated_at
    BEFORE UPDATE ON item_transaksi
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ================================================
-- 4. INSERT data dummyy
-- ================================================

-- products dummy data
INSERT INTO products (name, price, stock) VALUES
('Laptop Dell Inspiron 15', 8500000.00, 5),
('Mouse Wireless Logitech', 250000.00, 25),
('Keyboard Mechanical RGB', 750000.00, 15),
('Monitor LED 24 inch', 2200000.00, 8),
('Headset Gaming', 450000.00, 12),
('Webcam HD 1080p', 350000.00, 20),
('Speaker Bluetooth', 180000.00, 30),
('Hard Drive External 1TB', 650000.00, 10),
('USB Flash Drive 32GB', 75000.00, 50),
('Power Bank 10000mAh', 150000.00, 40);

-- transaksi dummy data
INSERT INTO transactions (transaction_date, total_amount) VALUES
('2024-01-15 10:30:00', 8750000.00),
('2024-01-15 14:45:00', 1200000.00),
('2024-01-16 09:15:00', 500000.00);

--  transaction items dummy data
INSERT INTO transaction_items (transaction_id, product_id, product_name, price, quantity, subtotal) VALUES
-- Transaction 1
(1, 1, 'Laptop Dell Inspiron 15', 8500000.00, 1, 8500000.00),
(1, 2, 'Mouse Wireless Logitech', 250000.00, 1, 250000.00),

-- Transaction 2  
(2, 4, 'Monitor LED 24 inch', 2200000.00, 1, 2200000.00),
(2, 3, 'Keyboard Mechanical RGB', 750000.00, 1, 750000.00),
(2, 2, 'Mouse Wireless Logitech', 250000.00, 1, 250000.00),

-- Transaction 3
(3, 5, 'Headset Gaming', 450000.00, 1, 450000.00),
(3, 9, 'USB Flash Drive 32GB', 75000.00, 2, 150000.00);

-- ================================================
-- 5. UPDATE STOCK setelah transaksi
-- ================================================

-- Update stock setelah terjadi transaksi
UPDATE products SET stock = stock - 1 WHERE id = 1; -- Laptop
UPDATE products SET stock = stock - 2 WHERE id = 2; -- Mouse (sold 2 times)
UPDATE products SET stock = stock - 1 WHERE id = 3; -- Keyboard
UPDATE products SET stock = stock - 1 WHERE id = 4; -- Monitor
UPDATE products SET stock = stock - 1 WHERE id = 5; -- Headset
UPDATE products SET stock = stock - 2 WHERE id = 9; -- USB Flash Drive

-- ================================================
-- 6. CREATE VIEWS FOR REPORTING
-- ================================================

-- View untuk summary transaksi
CREATE VIEW v_transaction_summary AS
SELECT 
    t.id,
    t.transaction_date,
    t.total_amount,
    COUNT(ti.id) as total_items,
    SUM(ti.quantity) as total_quantity
FROM transactions t
LEFT JOIN transaction_items ti ON t.id = ti.transaction_id
WHERE t.deleted_at IS NULL
GROUP BY t.id, t.transaction_date, t.total_amount
ORDER BY t.transaction_date DESC;

-- View untuk laporan penjualan produk
CREATE VIEW v_product_sales_report AS
SELECT 
    p.id,
    p.name as product_name,
    p.price as current_price,
    p.stock as current_stock,
    COALESCE(SUM(ti.quantity), 0) as total_sold,
    COALESCE(SUM(ti.subtotal), 0) as total_revenue
FROM products p
LEFT JOIN transaction_items ti ON p.id = ti.product_id
LEFT JOIN transactions t ON ti.transaction_id = t.id
WHERE p.deleted_at IS NULL 
    AND (t.deleted_at IS NULL OR t.deleted_at IS NOT NULL)
GROUP BY p.id, p.name, p.price, p.stock
ORDER BY total_sold DESC;

-- View untuk alert stok rendah
CREATE VIEW v_low_stock_alert AS
SELECT 
    id,
    name,
    price,
    stock,
    CASE 
        WHEN stock = 0 THEN 'OUT_OF_STOCK'
        WHEN stock <= 5 THEN 'LOW_STOCK'
        WHEN stock <= 10 THEN 'WARNING'
        ELSE 'NORMAL'
    END as stock_status
FROM products 
WHERE deleted_at IS NULL 
    AND stock <= 10
ORDER BY stock ASC;

-- ================================================
-- 7. CREATE STORED PROCEDURES/FUNCTIONS
-- ================================================

-- Function untuk menghitung total transaksi
CREATE OR REPLACE FUNCTION calculate_transaction_total(p_transaction_id INTEGER)
RETURNS DECIMAL(15,2) AS $$
DECLARE
    total_amount DECIMAL(15,2);
BEGIN
    SELECT COALESCE(SUM(subtotal), 0) 
    INTO total_amount
    FROM transaction_items 
    WHERE transaction_id = p_transaction_id;
    
    RETURN total_amount;
END;
$$ LANGUAGE plpgsql;

-- Function untuk mengecek stok produk yang tersedia
CREATE OR REPLACE FUNCTION check_stock_availability(p_product_id INTEGER, p_quantity INTEGER)
RETURNS BOOLEAN AS $$
DECLARE
    current_stock INTEGER;
BEGIN
    SELECT stock INTO current_stock 
    FROM products 
    WHERE id = p_product_id AND deleted_at IS NULL;
    
    IF current_stock IS NULL THEN
        RETURN FALSE;
    END IF;
    
    RETURN current_stock >= p_quantity;
END;
$$ LANGUAGE plpgsql;

-- Function untuk mengupdate stok produk setelah transaksi
CREATE OR REPLACE FUNCTION update_product_stock(p_product_id INTEGER, p_quantity INTEGER)
RETURNS BOOLEAN AS $$
DECLARE
    current_stock INTEGER;
BEGIN
    -- ambil stock saat ini
    SELECT stock INTO current_stock 
    FROM products 
    WHERE id = p_product_id AND deleted_at IS NULL;
    
    IF current_stock IS NULL OR current_stock < p_quantity THEN
        RETURN FALSE;
    END IF;
    
    -- Update stock
    UPDATE products 
    SET stock = stock - p_quantity,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_product_id;
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;


-- ================================================
-- 8. USEFUL QUERIES FOR TESTING
-- ================================================

-- Check all tables and their row counts
SELECT 
    schemaname,
    tablename,
    attname,
    n_distinct,
    correlation
FROM pg_stats
WHERE schemaname = 'public'
ORDER BY tablename, attname;

-- ambil ukuran tabel
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

