-- Run this ONCE to set up the database
-- mysql -u root -p < migrations/001_create_products.sql

CREATE DATABASE IF NOT EXISTS productsdb
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE productsdb;

CREATE TABLE IF NOT EXISTS products (
  id          VARCHAR(36)    NOT NULL,          -- UUID stored as string
  name        VARCHAR(255)   NOT NULL,
  description TEXT,
  price       DECIMAL(10,2)  NOT NULL DEFAULT 0.00,
  stock       INT            NOT NULL DEFAULT 0,
  category    VARCHAR(100)   NOT NULL DEFAULT '',
  created_at  DATETIME       NOT NULL,
  updated_at  DATETIME       NOT NULL,

  PRIMARY KEY (id),
  INDEX idx_category (category),               -- fast category filter
  INDEX idx_created_at (created_at DESC)        -- fast ORDER BY created_at DESC
);

-- Seed data so you have something to work with immediately
INSERT INTO products (id, name, description, price, stock, category, created_at, updated_at) VALUES
  (UUID(), 'iPhone 15 Pro',      'Apple flagship smartphone',          999.99, 50,  'Electronics', NOW(), NOW()),
  (UUID(), 'Samsung Galaxy S24', 'Android flagship smartphone',        849.99, 35,  'Electronics', NOW(), NOW()),
  (UUID(), 'Sony WH-1000XM5',   'Noise cancelling headphones',        349.99, 80,  'Electronics', NOW(), NOW()),
  (UUID(), 'The Go Book',        'Learning Go programming',             39.99, 200, 'Books',       NOW(), NOW()),
  (UUID(), 'Standing Desk',      'Electric height adjustable desk',    499.99, 15,  'Furniture',   NOW(), NOW()),
  (UUID(), 'Mechanical Keyboard','Cherry MX Blue switches',            129.99, 60,  'Electronics', NOW(), NOW());
