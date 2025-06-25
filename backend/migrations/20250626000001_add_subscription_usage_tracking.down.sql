-- Drop triggers
DROP TRIGGER IF EXISTS update_subscription_usage_timestamp ON subscription_usage;

-- Drop functions
DROP FUNCTION IF EXISTS update_subscription_usage_timestamp();

-- Drop columns from subscriptions
ALTER TABLE subscriptions DROP COLUMN IF EXISTS trial_ends_at;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS payment_failed_at;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS usage_reset_at;

-- Drop columns from subscription_plans
ALTER TABLE subscription_plans DROP COLUMN IF EXISTS version;
ALTER TABLE subscription_plans DROP COLUMN IF EXISTS is_popular;
ALTER TABLE subscription_plans DROP COLUMN IF EXISTS trial_days;

-- Drop tables
DROP TABLE IF EXISTS usage_events;
DROP TABLE IF EXISTS subscription_usage;