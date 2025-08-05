-- Rollback: Change user_id from UUID back to integer
-- WARNING: This rollback will cause data loss

-- Drop existing index on user_id
DROP INDEX IF EXISTS idx_subscriptions_user_id;

-- Change user_id column type from UUID back to INTEGER
-- Note: This will cause data loss and may fail if there are UUID values
ALTER TABLE subscriptions ALTER COLUMN user_id TYPE INTEGER USING 1;

-- Recreate index on user_id
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);

-- Remove comment
COMMENT ON COLUMN subscriptions.user_id IS NULL;