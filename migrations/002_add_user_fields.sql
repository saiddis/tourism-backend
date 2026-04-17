-- Migration: 002_add_user_fields.sql
-- Add avatar_url and balance columns to users table

ALTER TABLE users ADD COLUMN avatar_url VARCHAR(500);
ALTER TABLE users ADD COLUMN balance DECIMAL(12,2) DEFAULT 0.00;
