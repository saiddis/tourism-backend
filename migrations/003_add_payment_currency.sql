-- Migration: 003_add_payment_currency.sql
-- Add currency column to payments table

ALTER TABLE payments ADD COLUMN currency VARCHAR(3) DEFAULT 'TJS';
