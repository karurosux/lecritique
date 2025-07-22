-- Drop all foreign key constraints first
ALTER TABLE "public"."accounts" DROP CONSTRAINT IF EXISTS "fk_accounts_subscription";
ALTER TABLE "public"."feedbacks" DROP CONSTRAINT IF EXISTS "feedbacks_organization_id_fkey";
ALTER TABLE "public"."feedbacks" DROP CONSTRAINT IF EXISTS "feedbacks_product_id_fkey";
ALTER TABLE "public"."feedbacks" DROP CONSTRAINT IF EXISTS "feedbacks_qr_code_id_fkey";
ALTER TABLE "public"."locations" DROP CONSTRAINT IF EXISTS "locations_organization_id_fkey";
ALTER TABLE "public"."organizations" DROP CONSTRAINT IF EXISTS "organizations_account_id_fkey";
ALTER TABLE "public"."products" DROP CONSTRAINT IF EXISTS "products_organization_id_fkey";
ALTER TABLE "public"."qr_codes" DROP CONSTRAINT IF EXISTS "qr_codes_organization_id_fkey";
ALTER TABLE "public"."questionnaires" DROP CONSTRAINT IF EXISTS "questionnaires_organization_id_fkey";
ALTER TABLE "public"."questionnaires" DROP CONSTRAINT IF EXISTS "questionnaires_product_id_fkey";
ALTER TABLE "public"."questions" DROP CONSTRAINT IF EXISTS "fk_questions_product";
ALTER TABLE "public"."subscription_usage" DROP CONSTRAINT IF EXISTS "subscription_usage_subscription_id_fkey";
ALTER TABLE "public"."subscriptions" DROP CONSTRAINT IF EXISTS "subscriptions_account_id_fkey";
ALTER TABLE "public"."subscriptions" DROP CONSTRAINT IF EXISTS "subscriptions_plan_id_fkey";
ALTER TABLE "public"."team_invitations" DROP CONSTRAINT IF EXISTS "team_invitations_account_id_fkey";
ALTER TABLE "public"."team_members" DROP CONSTRAINT IF EXISTS "team_members_account_id_fkey";
ALTER TABLE "public"."team_members" DROP CONSTRAINT IF EXISTS "team_members_member_id_fkey";
ALTER TABLE "public"."usage_events" DROP CONSTRAINT IF EXISTS "usage_events_subscription_id_fkey";
ALTER TABLE "public"."verification_tokens" DROP CONSTRAINT IF EXISTS "fk_verification_tokens_account_id";

-- Drop all tables
DROP TABLE IF EXISTS "public"."usage_events";
DROP TABLE IF EXISTS "public"."team_members";
DROP TABLE IF EXISTS "public"."team_invitations";
DROP TABLE IF EXISTS "public"."subscription_usage";
DROP TABLE IF EXISTS "public"."verification_tokens";
DROP TABLE IF EXISTS "public"."feedbacks";
DROP TABLE IF EXISTS "public"."question_templates";
DROP TABLE IF EXISTS "public"."questions";
DROP TABLE IF EXISTS "public"."questionnaires";
DROP TABLE IF EXISTS "public"."qr_codes";
DROP TABLE IF EXISTS "public"."products";
DROP TABLE IF EXISTS "public"."locations";
DROP TABLE IF EXISTS "public"."organizations";
DROP TABLE IF EXISTS "public"."subscriptions";
DROP TABLE IF EXISTS "public"."subscription_plans";
DROP TABLE IF EXISTS "public"."accounts";

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";