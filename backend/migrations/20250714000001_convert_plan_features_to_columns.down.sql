-- Rollback: Convert plan features from columns back to JSON

-- Re-add features column
ALTER TABLE subscription_plans ADD COLUMN features JSONB;

-- Migrate data back to JSON format
UPDATE subscription_plans
SET features = jsonb_build_object(
    'limits', jsonb_build_object(
        'max_restaurants', max_restaurants,
        'max_qr_codes', max_qr_codes,
        'max_feedbacks_per_month', max_feedbacks_per_month,
        'max_team_members', max_team_members
    ),
    'flags', jsonb_build_object(
        'basic_analytics', has_basic_analytics,
        'advanced_analytics', has_advanced_analytics,
        'feedback_explorer', has_feedback_explorer,
        'custom_branding', has_custom_branding,
        'priority_support', has_priority_support
    ),
    'custom', '{}'::jsonb
);

-- Make features column NOT NULL after migration
ALTER TABLE subscription_plans ALTER COLUMN features SET NOT NULL;

-- Drop the column-based fields
ALTER TABLE subscription_plans
DROP COLUMN max_restaurants,
DROP COLUMN max_qr_codes,
DROP COLUMN max_feedbacks_per_month,
DROP COLUMN max_team_members,
DROP COLUMN has_basic_analytics,
DROP COLUMN has_advanced_analytics,
DROP COLUMN has_feedback_explorer,
DROP COLUMN has_custom_branding,
DROP COLUMN has_priority_support;