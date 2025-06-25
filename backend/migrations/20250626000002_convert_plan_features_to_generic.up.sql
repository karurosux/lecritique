-- Migration to convert plan features to generic structure

-- First, let's add a temporary column to store the new format
ALTER TABLE subscription_plans ADD COLUMN features_new JSONB;

-- Convert existing features to new format
UPDATE subscription_plans 
SET features_new = jsonb_build_object(
    'limits', jsonb_build_object(
        'max_restaurants', COALESCE((features->>'max_restaurants')::int, 0),
        'max_locations_per_restaurant', COALESCE((features->>'max_locations_per_restaurant')::int, 0),
        'max_qr_codes_per_location', COALESCE((features->>'max_qr_codes_per_location')::int, 0),
        'max_feedbacks_per_month', COALESCE((features->>'max_feedbacks_per_month')::int, 0),
        'max_team_members', COALESCE((features->>'max_team_members')::int, 0)
    ),
    'flags', jsonb_build_object(
        'advanced_analytics', COALESCE((features->>'advanced_analytics')::boolean, false),
        'custom_branding', COALESCE((features->>'custom_branding')::boolean, false),
        'api_access', COALESCE((features->>'api_access')::boolean, false),
        'priority_support', COALESCE((features->>'priority_support')::boolean, false)
    ),
    'custom', '{}'::jsonb
);

-- Drop old column and rename new one
ALTER TABLE subscription_plans DROP COLUMN features;
ALTER TABLE subscription_plans RENAME COLUMN features_new TO features;

-- Example: Add new features to existing plans
-- This shows how easy it is to add new features without code changes

-- Add storage limit to all plans
UPDATE subscription_plans 
SET features = jsonb_set(
    features,
    '{limits,max_storage_gb}',
    CASE 
        WHEN code = 'starter' THEN '10'
        WHEN code = 'professional' THEN '50'
        WHEN code = 'enterprise' THEN '-1'
        ELSE '10'
    END::jsonb
)
WHERE features IS NOT NULL;

-- Add white label feature to enterprise plan
UPDATE subscription_plans 
SET features = jsonb_set(
    features,
    '{flags,white_label}',
    'true'::jsonb
)
WHERE code = 'enterprise';

-- Add API rate limits
UPDATE subscription_plans 
SET features = jsonb_set(
    features,
    '{limits,max_api_calls_per_hour}',
    CASE 
        WHEN code = 'starter' THEN '1000'
        WHEN code = 'professional' THEN '10000'
        WHEN code = 'enterprise' THEN '-1'
        ELSE '1000'
    END::jsonb
)
WHERE features IS NOT NULL;