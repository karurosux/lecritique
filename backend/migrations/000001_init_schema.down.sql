-- Drop triggers
DROP TRIGGER IF EXISTS update_feedbacks_updated_at ON feedbacks;
DROP TRIGGER IF EXISTS update_question_templates_updated_at ON question_templates;
DROP TRIGGER IF EXISTS update_questions_updated_at ON questions;
DROP TRIGGER IF EXISTS update_questionnaires_updated_at ON questionnaires;
DROP TRIGGER IF EXISTS update_qr_codes_updated_at ON qr_codes;
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP TRIGGER IF EXISTS update_locations_updated_at ON locations;
DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;
DROP TRIGGER IF EXISTS update_team_members_updated_at ON team_members;
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
DROP TRIGGER IF EXISTS update_subscription_plans_updated_at ON subscription_plans;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order
DROP TABLE IF EXISTS feedbacks;
DROP TABLE IF EXISTS question_templates;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS questionnaires;
DROP TABLE IF EXISTS qr_codes;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS users;

-- Drop foreign key constraint first
ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS fk_accounts_subscription;

DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS subscription_plans;
