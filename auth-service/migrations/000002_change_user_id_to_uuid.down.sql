-- Rollback: Change user ID from UUID back to serial
-- WARNING: This rollback will cause data loss of UUID values

-- Drop the UUID primary key
ALTER TABLE users DROP CONSTRAINT users_pkey;

-- Add new serial ID column
ALTER TABLE users ADD COLUMN user_serial_id SERIAL;

-- Drop the UUID id column
ALTER TABLE users DROP COLUMN id;

-- Rename serial column to id
ALTER TABLE users RENAME COLUMN user_serial_id TO id;

-- Add primary key constraint
ALTER TABLE users ADD PRIMARY KEY (id);

-- Recreate deleted_at index
DROP INDEX IF EXISTS idx_users_deleted_at;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);