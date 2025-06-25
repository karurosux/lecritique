-- Remove the custom plan example
DELETE FROM subscription_plans WHERE code = 'enterprise_custom_acme';

-- Drop the index
DROP INDEX IF EXISTS idx_subscription_plans_visible;

-- Remove the visibility column
ALTER TABLE subscription_plans 
DROP COLUMN IF EXISTS is_visible;