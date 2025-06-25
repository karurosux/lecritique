-- Add visibility flag to subscription plans
ALTER TABLE subscription_plans 
ADD COLUMN IF NOT EXISTS is_visible BOOLEAN DEFAULT true;

-- Create index for efficient querying of visible plans
CREATE INDEX idx_subscription_plans_visible 
ON subscription_plans(is_visible, is_active) 
WHERE is_visible = true AND is_active = true;

-- Example: Create a custom hidden plan for a specific enterprise customer
INSERT INTO subscription_plans (
    code,
    name,
    description,
    price,
    currency,
    interval,
    is_active,
    is_visible,
    features
) VALUES (
    'enterprise_custom_acme',
    'Enterprise Custom - ACME Corp',
    'Custom plan for ACME Corporation with special pricing and features',
    299.00,
    'USD',
    'month',
    true,
    false, -- Hidden from public listing
    jsonb_build_object(
        'limits', jsonb_build_object(
            'max_restaurants', -1,
            'max_locations_per_restaurant', -1,
            'max_qr_codes_per_location', -1,
            'max_feedbacks_per_month', -1,
            'max_team_members', -1,
            'max_storage_gb', 500,
            'max_api_calls_per_hour', 50000
        ),
        'flags', jsonb_build_object(
            'advanced_analytics', true,
            'custom_branding', true,
            'api_access', true,
            'priority_support', true,
            'white_label', true,
            'custom_domain', true,
            'dedicated_support_manager', true,
            'sla_guarantee', true
        ),
        'custom', jsonb_build_object(
            'discount_percentage', 25,
            'contract_months', 24,
            'dedicated_server', true
        )
    )
) ON CONFLICT (code) DO NOTHING;