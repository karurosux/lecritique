-- Convert plan features from JSON to individual columns

-- Add new columns for limits
ALTER TABLE subscription_plans
ADD COLUMN max_organizations INTEGER DEFAULT 1,
ADD COLUMN max_qr_codes INTEGER DEFAULT 5,
ADD COLUMN max_feedbacks_per_month INTEGER DEFAULT 50,
ADD COLUMN max_team_members INTEGER DEFAULT 2;

-- Add new columns for feature flags
ALTER TABLE subscription_plans
ADD COLUMN has_basic_analytics BOOLEAN DEFAULT false,
ADD COLUMN has_advanced_analytics BOOLEAN DEFAULT false,
ADD COLUMN has_feedback_explorer BOOLEAN DEFAULT false,
ADD COLUMN has_custom_branding BOOLEAN DEFAULT false,
ADD COLUMN has_priority_support BOOLEAN DEFAULT false;

-- Migrate existing data from JSON to columns
UPDATE subscription_plans
SET 
    max_organizations = COALESCE((features->'limits'->>'max_organizations')::int, 1),
    max_qr_codes = COALESCE((features->'limits'->>'max_qr_codes')::int, 5),
    max_feedbacks_per_month = COALESCE((features->'limits'->>'max_feedbacks_per_month')::int, 50),
    max_team_members = COALESCE((features->'limits'->>'max_team_members')::int, 2),
    has_basic_analytics = COALESCE((features->'flags'->>'basic_analytics')::boolean, false),
    has_advanced_analytics = COALESCE((features->'flags'->>'advanced_analytics')::boolean, false),
    has_feedback_explorer = COALESCE((features->'flags'->>'feedback_explorer')::boolean, false),
    has_custom_branding = COALESCE((features->'flags'->>'custom_branding')::boolean, false),
    has_priority_support = COALESCE((features->'flags'->>'priority_support')::boolean, false);

-- Add constraints to ensure valid values
ALTER TABLE subscription_plans
ADD CONSTRAINT check_max_organizations CHECK (max_organizations >= -1),
ADD CONSTRAINT check_max_qr_codes CHECK (max_qr_codes >= -1),
ADD CONSTRAINT check_max_feedbacks CHECK (max_feedbacks_per_month >= -1),
ADD CONSTRAINT check_max_team_members CHECK (max_team_members >= -1);

-- Add NOT NULL constraints after data migration
ALTER TABLE subscription_plans
ALTER COLUMN max_organizations SET NOT NULL,
ALTER COLUMN max_qr_codes SET NOT NULL,
ALTER COLUMN max_feedbacks_per_month SET NOT NULL,
ALTER COLUMN max_team_members SET NOT NULL,
ALTER COLUMN has_basic_analytics SET NOT NULL,
ALTER COLUMN has_advanced_analytics SET NOT NULL,
ALTER COLUMN has_feedback_explorer SET NOT NULL,
ALTER COLUMN has_custom_branding SET NOT NULL,
ALTER COLUMN has_priority_support SET NOT NULL;

-- Drop the old features column
ALTER TABLE subscription_plans DROP COLUMN features;
