-- Migration: Change user_id from integer to UUID in subscriptions table
-- This migration will:
-- 1. Enable UUID extension
-- 2. Change user_id column type to UUID
-- 3. Update indexes

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop existing index on user_id
DROP INDEX IF EXISTS idx_subscriptions_user_id;

-- Change user_id column type from INTEGER to UUID
-- Note: This will require data migration if there's existing data
-- For new installations, this is safe
ALTER TABLE subscriptions ALTER COLUMN user_id TYPE UUID USING uuid_generate_v4();

-- Recreate index on user_id
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);

-- Add a comment to clarify the relationship
COMMENT ON COLUMN subscriptions.user_id IS 'UUID reference to user from auth-service';