-- Migration: Change user ID from serial to UUID
-- This migration will:
-- 1. Add UUID extension if not exists
-- 2. Add a new UUID column
-- 3. Populate UUIDs for existing users
-- 4. Drop the old serial ID
-- 5. Rename UUID column to id

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Add new UUID column
ALTER TABLE users ADD COLUMN user_uuid UUID DEFAULT uuid_generate_v4();

-- Update existing records to have UUIDs (if any exist)
UPDATE users SET user_uuid = uuid_generate_v4() WHERE user_uuid IS NULL;

-- Make UUID column NOT NULL
ALTER TABLE users ALTER COLUMN user_uuid SET NOT NULL;

-- Drop the old serial ID column
ALTER TABLE users DROP COLUMN id;

-- Rename UUID column to id
ALTER TABLE users RENAME COLUMN user_uuid TO id;

-- Add primary key constraint
ALTER TABLE users ADD PRIMARY KEY (id);

-- Update deleted_at index
DROP INDEX IF EXISTS idx_users_deleted_at;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);