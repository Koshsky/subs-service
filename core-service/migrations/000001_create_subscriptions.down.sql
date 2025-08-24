-- Rollback subscriptions table creation
DROP INDEX IF EXISTS idx_subscriptions_service_name;
DROP INDEX IF EXISTS idx_subscriptions_deleted_at;
DROP INDEX IF EXISTS idx_subscriptions_user_id;
DROP TABLE IF EXISTS subscriptions;