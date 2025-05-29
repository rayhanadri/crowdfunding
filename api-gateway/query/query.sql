-- Tabel "users"."users" (Pengguna)
CREATE TABLE "users"."users" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Campaigns (Kampanye Donasi)
CREATE TABLE "campaigns"."campaigns" (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(255) NOT NULL,
    target_amount NUMERIC(15,2) NOT NULL,
    collected_amount NUMERIC(15,2) DEFAULT 0,
    deadline DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'ACTIVE', -- e.g., ACTIVE, COMPLETED
	category VARCHAR(50),
	min_donation NUMERIC(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel "donations"."donations" (Riwayat Donasi)
CREATE TABLE "donations"."donations" (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    campaign_id INTEGER NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    message VARCHAR(255),
	status VARCHAR(50) DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel "transactions"."transactions" (Transaksi Keuangan)
CREATE TABLE "transactions"."transactions" (
    id SERIAL PRIMARY KEY,
    donation_id INTEGER NOT NULL,
	invoice_id VARCHAR(255),
	invoice_url VARCHAR(255),
	invoice_description VARCHAR(255),
    payment_method VARCHAR(50),
	amount NUMERIC(15,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Example Data
INSERT INTO "users"."users" (name, email, password)
VALUES
('Andi Wijaya', 'andi@example.com', 'hashed_password_1'),
('Siti Rahma', 'siti@example.com', 'hashed_password_2'),
('Budi Santoso', 'budi@example.com', 'hashed_password_3'),
('Dewi Lestari', 'dewi@example.com', 'hashed_password_4'),
('Rizky Maulana', 'rizky@example.com', 'hashed_password_5');

INSERT INTO campaigns (user_id, title, description, target_amount, collected_amount, deadline, status, category, min_donation)
VALUES
(1, 'Bantu Biaya Pengobatan Anak Yatim', 'Menggalang dana untuk pengobatan anak yatim', 10000000, 2500000, '2025-12-31', 'ACTIVE', 'Kesehatan', 10000),
(2, 'Renovasi Mushola Desa', 'Kami butuh bantuan untuk merenovasi mushola desa kami', 15000000, 5000000, '2025-10-15', 'ACTIVE', 'Tempat Ibadah', 20000),
(3, 'Pendidikan Anak Kurang Mampu', 'Donasi untuk membeli perlengkapan sekolah', 8000000, 3000000, '2025-09-01', 'ACTIVE', 'Pendidikan', 15000),
(4, 'Bantu Korban Banjir', 'Penggalangan dana untuk korban banjir di Kalimantan', 20000000, 12000000, '2025-11-20', 'ACTIVE', 'Bencana Alam', 10000),
(5, 'Dukung Usaha Kecil Ibu Rumah Tangga', 'Modal usaha kecil untuk ibu rumah tangga', 5000000, 1500000, '2025-08-01', 'ACTIVE', 'Ekonomi', 5000);

INSERT INTO "donations"."donations" (user_id, campaign_id, amount, message, status)
VALUES
(2, 1, 50000, 'Semoga lekas sembuh', 'PAID'),
(3, 2, 100000, 'Semoga mushola cepat direnovasi', 'PAID'),
(4, 3, 25000, 'Bantu pendidikan generasi penerus', 'PENDING'),
(5, 4, 75000, 'Bantuan untuk saudara di Kalimantan', 'PAID'),
(1, 5, 30000, 'Semoga usaha berjalan lancar', 'PENDING');

INSERT INTO "transactions"."transactions" (
    donation_id, invoice_id, invoice_url, invoice_description,
    payment_method, amount, status
) VALUES
(1, 'INV-20240520-001', 'https://example.com/invoice/INV-20240520-001', 'Donasi untuk Pengobatan Anak Yatim', 'bank_transfer', 50000, 'PAID'),
(2, 'INV-20240520-002', 'https://example.com/invoice/INV-20240520-002', 'Donasi untuk Renovasi Mushola', 'ewallet', 100000, 'PAID'),
(3, 'INV-20240520-003', 'https://example.com/invoice/INV-20240520-003', 'Donasi Pendidikan Anak Kurang Mampu', 'bank_transfer', 25000, 'PENDING'),
(4, 'INV-20240520-004', 'https://example.com/invoice/INV-20240520-004', 'Bantuan Korban Banjir', 'credit_card', 75000, 'PAID'),
(5, 'INV-20240520-005', 'https://example.com/invoice/INV-20240520-005', 'Donasi Modal Usaha Ibu Rumah Tangga', 'ewallet', 30000, 'PENDING');



