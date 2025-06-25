-- Revert to old structure

-- Add temporary column for old format
ALTER TABLE subscription_plans ADD COLUMN features_old JSONB;

-- Convert back to old format
UPDATE subscription_plans 
SET features_old = jsonb_build_object(
    'max_restaurants', COALESCE((features->'limits'->>'max_restaurants')::int, 0),
    'max_locations_per_restaurant', COALESCE((features->'limits'->>'max_locations_per_restaurant')::int, 0),
    'max_qr_codes_per_location', COALESCE((features->'limits'->>'max_qr_codes_per_location')::int, 0),
    'max_feedbacks_per_month', COALESCE((features->'limits'->>'max_feedbacks_per_month')::int, 0),
    'max_team_members', COALESCE((features->'limits'->>'max_team_members')::int, 0),
    'advanced_analytics', COALESCE((features->'flags'->>'advanced_analytics')::boolean, false),
    'custom_branding', COALESCE((features->'flags'->>'custom_branding')::boolean, false),
    'api_access', COALESCE((features->'flags'->>'api_access')::boolean, false),
    'priority_support', COALESCE((features->'flags'->>'priority_support')::boolean, false)
);

-- Drop new column and rename old one
ALTER TABLE subscription_plans DROP COLUMN features;
ALTER TABLE subscription_plans RENAME COLUMN features_old TO features;