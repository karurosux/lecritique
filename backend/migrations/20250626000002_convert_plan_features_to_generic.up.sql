-- Migration to convert plan features to generic structure

-- First, let's add a temporary column to store the new format
ALTER TABLE subscription_plans ADD COLUMN features_new JSONB;

-- Convert existing features to new format
UPDATE subscription_plans 
SET features_new = jsonb_build_object(
    'limits', jsonb_build_object(
        'max_organizations', COALESCE((features->>'max_organizations')::int, 0),
        'max_qr_codes', COALESCE((features->>'max_qr_codes')::int, 0),
        'max_feedbacks_per_month', COALESCE((features->>'max_feedbacks_per_month')::int, 0),
        'max_team_members', COALESCE((features->>'max_team_members')::int, 0)
    ),
    'flags', jsonb_build_object(
        'basic_analytics', COALESCE((features->>'basic_analytics')::boolean, false),
        'advanced_analytics', COALESCE((features->>'advanced_analytics')::boolean, false),
        'feedback_explorer', COALESCE((features->>'feedback_explorer')::boolean, false),
        'custom_branding', COALESCE((features->>'custom_branding')::boolean, false),
        'priority_support', COALESCE((features->>'priority_support')::boolean, false)
    ),
    'custom', '{}'::jsonb
);

-- Drop old column and rename new one
ALTER TABLE subscription_plans DROP COLUMN features;
ALTER TABLE subscription_plans RENAME COLUMN features_new TO features;
