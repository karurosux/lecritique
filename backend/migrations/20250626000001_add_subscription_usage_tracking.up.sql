-- Create subscription_usage table for tracking usage metrics
CREATE TABLE IF NOT EXISTS subscription_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    feedbacks_count INTEGER DEFAULT 0,
    restaurants_count INTEGER DEFAULT 0,
    locations_count INTEGER DEFAULT 0,
    qr_codes_count INTEGER DEFAULT 0,
    team_members_count INTEGER DEFAULT 0,
    last_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(subscription_id, period_start, period_end)
);

-- Create indexes for performance
CREATE INDEX idx_subscription_usage_subscription_id ON subscription_usage(subscription_id);
CREATE INDEX idx_subscription_usage_period ON subscription_usage(period_start, period_end);

-- Create usage_events table for audit trail
CREATE TABLE IF NOT EXISTS usage_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- create, delete, update
    resource_type VARCHAR(50) NOT NULL, -- feedback, restaurant, location, etc
    resource_id UUID,
    metadata JSONB,
    CHECK (event_type IN ('create', 'delete', 'update')),
    CHECK (resource_type IN ('feedback', 'restaurant', 'location', 'qr_code', 'team_member'))
);

-- Create indexes for usage_events
CREATE INDEX idx_usage_events_subscription_id ON usage_events(subscription_id);
CREATE INDEX idx_usage_events_created_at ON usage_events(created_at);

-- Add subscription enforcement columns to existing tables (if not exists)
ALTER TABLE subscription_plans ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;
ALTER TABLE subscription_plans ADD COLUMN IF NOT EXISTS is_popular BOOLEAN DEFAULT FALSE;
ALTER TABLE subscription_plans ADD COLUMN IF NOT EXISTS trial_days INTEGER DEFAULT 0;

-- Add columns for better subscription management
ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS trial_ends_at TIMESTAMPTZ;
ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS payment_failed_at TIMESTAMPTZ;
ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS usage_reset_at TIMESTAMPTZ;

-- Create a function to update last_updated_at
CREATE OR REPLACE FUNCTION update_subscription_usage_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for subscription_usage
CREATE TRIGGER update_subscription_usage_timestamp
    BEFORE UPDATE ON subscription_usage
    FOR EACH ROW
    EXECUTE FUNCTION update_subscription_usage_timestamp();

-- Initialize usage records for existing active subscriptions
INSERT INTO subscription_usage (
    subscription_id,
    period_start,
    period_end,
    restaurants_count,
    feedbacks_count,
    locations_count,
    qr_codes_count,
    team_members_count
)
SELECT 
    s.id,
    s.current_period_start,
    s.current_period_end,
    (SELECT COUNT(*) FROM restaurants WHERE account_id = s.account_id AND deleted_at IS NULL),
    0, -- We don't have historical feedback count
    (SELECT COUNT(*) FROM locations l JOIN restaurants r ON l.restaurant_id = r.id WHERE r.account_id = s.account_id AND l.deleted_at IS NULL),
    (SELECT COUNT(*) FROM qr_codes q JOIN locations l ON q.location_id = l.id JOIN restaurants r ON l.restaurant_id = r.id WHERE r.account_id = s.account_id AND q.deleted_at IS NULL),
    (SELECT COUNT(*) FROM team_members WHERE account_id = s.account_id AND deleted_at IS NULL)
FROM subscriptions s
WHERE s.status = 'active' 
AND s.deleted_at IS NULL
ON CONFLICT (subscription_id, period_start, period_end) DO NOTHING;