-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create "accounts" table
CREATE TABLE "public"."accounts" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "email" character varying(255) NOT NULL,
  "password_hash" character varying(255) NOT NULL,
  "name" character varying(255) NOT NULL,
  "phone" character varying(50) NULL,
  "is_active" boolean NULL DEFAULT true,
  "email_verified" boolean NULL DEFAULT false,
  "email_verified_at" timestamptz NULL,
  "subscription_id" uuid NULL,
  "deactivation_requested_at" timestamptz NULL,
  "first_name" character varying(255) NULL DEFAULT '',
  "last_name" character varying(255) NULL DEFAULT '',
  PRIMARY KEY ("id"),
  CONSTRAINT "accounts_email_key" UNIQUE ("email")
);

-- Create "subscription_plans" table
CREATE TABLE "public"."subscription_plans" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  "code" character varying(50) NOT NULL,
  "description" text NULL,
  "price" numeric(10,2) NOT NULL,
  "currency" character varying(3) NULL DEFAULT 'USD',
  "interval" character varying(20) NULL DEFAULT 'month',
  "is_active" boolean NULL DEFAULT true,
  "stripe_price_id" character varying(255) NULL,
  "version" integer NULL DEFAULT 1,
  "is_popular" boolean NULL DEFAULT false,
  "trial_days" integer NULL DEFAULT 0,
  "is_visible" boolean NULL DEFAULT true,
  "max_organizations" integer NOT NULL DEFAULT 1,
  "max_qr_codes" integer NOT NULL DEFAULT 5,
  "max_feedbacks_per_month" integer NOT NULL DEFAULT 50,
  "max_team_members" integer NOT NULL DEFAULT 2,
  "has_basic_analytics" boolean NOT NULL DEFAULT false,
  "has_advanced_analytics" boolean NOT NULL DEFAULT false,
  "has_feedback_explorer" boolean NOT NULL DEFAULT false,
  "has_custom_branding" boolean NOT NULL DEFAULT false,
  "has_priority_support" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "subscription_plans_code_key" UNIQUE ("code"),
  CONSTRAINT "check_max_feedbacks" CHECK (max_feedbacks_per_month >= '-1'::integer),
  CONSTRAINT "check_max_organizations" CHECK (max_organizations >= '-1'::integer),
  CONSTRAINT "check_max_qr_codes" CHECK (max_qr_codes >= '-1'::integer),
  CONSTRAINT "check_max_team_members" CHECK (max_team_members >= '-1'::integer)
);

-- Create "subscriptions" table
CREATE TABLE "public"."subscriptions" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "account_id" uuid NOT NULL,
  "plan_id" uuid NOT NULL,
  "status" character varying(50) NOT NULL,
  "current_period_start" timestamptz NOT NULL,
  "current_period_end" timestamptz NOT NULL,
  "cancel_at" timestamptz NULL,
  "cancelled_at" timestamptz NULL,
  "stripe_customer_id" character varying(255) NULL,
  "stripe_subscription_id" character varying(255) NULL,
  "trial_ends_at" timestamptz NULL,
  "payment_failed_at" timestamptz NULL,
  "usage_reset_at" timestamptz NULL,
  PRIMARY KEY ("id")
);

-- Create "organizations" table
CREATE TABLE "public"."organizations" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "account_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "logo" character varying(500) NULL,
  "website" character varying(500) NULL,
  "phone" character varying(50) NULL,
  "email" character varying(255) NULL,
  "is_active" boolean NULL DEFAULT true,
  "settings" jsonb NULL DEFAULT '{}',
  PRIMARY KEY ("id")
);

-- Create "locations" table
CREATE TABLE "public"."locations" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "organization_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "address" character varying(500) NULL,
  "city" character varying(100) NULL,
  "state" character varying(100) NULL,
  "country" character varying(100) NULL,
  "postal_code" character varying(20) NULL,
  "latitude" numeric(10,8) NULL,
  "longitude" numeric(11,8) NULL,
  "is_active" boolean NULL DEFAULT true,
  PRIMARY KEY ("id")
);

-- Create "products" table
CREATE TABLE "public"."products" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "organization_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "category" character varying(100) NULL,
  "price" numeric(10,2) NULL,
  "currency" character varying(3) NULL DEFAULT 'USD',
  "image" character varying(500) NULL,
  "tags" text[] NULL,
  "is_available" boolean NULL DEFAULT true,
  "is_active" boolean NULL DEFAULT true,
  "display_order" integer NULL DEFAULT 0,
  PRIMARY KEY ("id")
);

-- Create "qr_codes" table
CREATE TABLE "public"."qr_codes" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "organization_id" uuid NOT NULL,
  "code" character varying(100) NOT NULL,
  "label" character varying(255) NULL,
  "type" character varying(50) NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "scans_count" integer NULL DEFAULT 0,
  "last_scanned_at" timestamptz NULL,
  "expires_at" timestamptz NULL,
  "location" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "qr_codes_code_key" UNIQUE ("code")
);

-- Create "questionnaires" table
CREATE TABLE "public"."questionnaires" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "organization_id" uuid NOT NULL,
  "product_id" uuid NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "is_default" boolean NULL DEFAULT false,
  "is_active" boolean NULL DEFAULT true,
  PRIMARY KEY ("id")
);

-- Create "questions" table
CREATE TABLE "public"."questions" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "text" text NOT NULL,
  "type" character varying(50) NOT NULL,
  "is_required" boolean NULL DEFAULT true,
  "display_order" integer NULL DEFAULT 0,
  "options" text[] NULL,
  "min_value" integer NULL,
  "max_value" integer NULL,
  "min_label" character varying(100) NULL,
  "max_label" character varying(100) NULL,
  "product_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);

-- Create "question_templates" table
CREATE TABLE "public"."question_templates" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "category" character varying(100) NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "text" text NOT NULL,
  "type" character varying(50) NOT NULL,
  "options" text[] NULL,
  "min_value" integer NULL,
  "max_value" integer NULL,
  "min_label" character varying(100) NULL,
  "max_label" character varying(100) NULL,
  "tags" text[] NULL,
  "is_active" boolean NULL DEFAULT true,
  PRIMARY KEY ("id")
);

-- Create "feedbacks" table
CREATE TABLE "public"."feedbacks" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "organization_id" uuid NOT NULL,
  "product_id" uuid NOT NULL,
  "qr_code_id" uuid NOT NULL,
  "customer_name" character varying(255) NULL,
  "customer_email" character varying(255) NULL,
  "customer_phone" character varying(50) NULL,
  "overall_rating" integer NULL,
  "responses" jsonb NOT NULL DEFAULT '[]',
  "device_info" jsonb NULL DEFAULT '{}',
  "is_complete" boolean NULL DEFAULT true,
  PRIMARY KEY ("id")
);

-- Create "verification_tokens" table
CREATE TABLE "public"."verification_tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "account_id" uuid NOT NULL,
  "token" character varying(128) NOT NULL,
  "type" character varying(50) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "used_at" timestamptz NULL,
  "created_at" timestamptz NULL DEFAULT now(),
  "updated_at" timestamptz NULL DEFAULT now(),
  "new_email" character varying(255) NULL DEFAULT NULL::character varying,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "verification_tokens_token_key" UNIQUE ("token"),
  CONSTRAINT "chk_verification_tokens_type" CHECK ((type)::text = ANY ((ARRAY['EMAIL_VERIFICATION'::character varying, 'PASSWORD_RESET'::character varying, 'TEAM_INVITE'::character varying])::text[]))
);

-- Create "subscription_usage" table
CREATE TABLE "public"."subscription_usage" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  "subscription_id" uuid NOT NULL,
  "period_start" timestamptz NOT NULL,
  "period_end" timestamptz NOT NULL,
  "feedbacks_count" integer NULL DEFAULT 0,
  "organizations_count" integer NULL DEFAULT 0,
  "locations_count" integer NULL DEFAULT 0,
  "qr_codes_count" integer NULL DEFAULT 0,
  "team_members_count" integer NULL DEFAULT 0,
  "last_updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "subscription_usage_subscription_id_period_start_period_end_key" UNIQUE ("subscription_id", "period_start", "period_end")
);

-- Create "usage_events" table
CREATE TABLE "public"."usage_events" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  "subscription_id" uuid NOT NULL,
  "event_type" character varying(50) NOT NULL,
  "resource_type" character varying(50) NOT NULL,
  "resource_id" uuid NULL,
  "metadata" jsonb NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "usage_events_event_type_check" CHECK ((event_type)::text = ANY ((ARRAY['create'::character varying, 'delete'::character varying, 'update'::character varying])::text[])),
  CONSTRAINT "usage_events_resource_type_check" CHECK ((resource_type)::text = ANY ((ARRAY['feedback'::character varying, 'organization'::character varying, 'location'::character varying, 'qr_code'::character varying, 'team_member'::character varying])::text[]))
);

-- Create "team_invitations" table
CREATE TABLE "public"."team_invitations" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "account_id" uuid NOT NULL,
  "email" character varying(255) NOT NULL,
  "role" character varying(50) NOT NULL,
  "invited_by" uuid NOT NULL,
  "token" character varying(255) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "accepted_at" timestamptz NULL,
  "email_accepted_at" timestamp NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "team_invitations_email_idx" UNIQUE ("account_id", "email", "deleted_at"),
  CONSTRAINT "team_invitations_token_key" UNIQUE ("token")
);

-- Create "team_members" table
CREATE TABLE "public"."team_members" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "account_id" uuid NOT NULL,
  "member_id" uuid NOT NULL,
  "role" character varying(50) NOT NULL,
  "invited_by" uuid NOT NULL,
  "invited_at" timestamptz NOT NULL,
  "accepted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);

-- Create indexes
CREATE INDEX "idx_accounts_deactivation_requested_at" ON "public"."accounts" ("deactivation_requested_at") WHERE (deactivation_requested_at IS NOT NULL);
CREATE INDEX "idx_accounts_email" ON "public"."accounts" ("email");
CREATE INDEX "idx_accounts_subscription_id" ON "public"."accounts" ("subscription_id");
CREATE INDEX "idx_feedbacks_created_at" ON "public"."feedbacks" ("created_at");
CREATE INDEX "idx_feedbacks_organization_id" ON "public"."feedbacks" ("organization_id");
CREATE INDEX "idx_feedbacks_product_id" ON "public"."feedbacks" ("product_id");
CREATE INDEX "idx_locations_organization_id" ON "public"."locations" ("organization_id");
CREATE INDEX "idx_organizations_account_id" ON "public"."organizations" ("account_id");
CREATE INDEX "idx_products_category" ON "public"."products" ("category");
CREATE INDEX "idx_products_organization_id" ON "public"."products" ("organization_id");
CREATE INDEX "idx_qr_codes_code" ON "public"."qr_codes" ("code");
CREATE INDEX "idx_qr_codes_organization_id" ON "public"."qr_codes" ("organization_id");
CREATE INDEX "idx_questionnaires_organization_id" ON "public"."questionnaires" ("organization_id");
CREATE INDEX "idx_questionnaires_product_id" ON "public"."questionnaires" ("product_id");
CREATE INDEX "idx_questions_product_display_order" ON "public"."questions" ("product_id", "display_order");
CREATE INDEX "idx_questions_product_id" ON "public"."questions" ("product_id");
CREATE INDEX "idx_subscription_plans_visible" ON "public"."subscription_plans" ("is_visible", "is_active") WHERE ((is_visible = true) AND (is_active = true));
CREATE INDEX "idx_subscription_usage_period" ON "public"."subscription_usage" ("period_start", "period_end");
CREATE INDEX "idx_subscription_usage_subscription_id" ON "public"."subscription_usage" ("subscription_id");
CREATE INDEX "idx_subscriptions_account_id" ON "public"."subscriptions" ("account_id");
CREATE INDEX "idx_subscriptions_status" ON "public"."subscriptions" ("status");
CREATE INDEX "idx_team_invitations_account_id" ON "public"."team_invitations" ("account_id");
CREATE INDEX "idx_team_invitations_email" ON "public"."team_invitations" ("email");
CREATE INDEX "idx_team_invitations_token" ON "public"."team_invitations" ("token");
CREATE INDEX "idx_team_members_account_id" ON "public"."team_members" ("account_id");
CREATE INDEX "idx_team_members_account_member" ON "public"."team_members" ("account_id", "member_id");
CREATE INDEX "idx_team_members_member_id" ON "public"."team_members" ("member_id");
CREATE INDEX "idx_team_members_user_id" ON "public"."team_members" ("member_id");
CREATE UNIQUE INDEX "team_members_account_id_user_id_unique" ON "public"."team_members" ("account_id", "member_id") WHERE (deleted_at IS NULL);
CREATE INDEX "idx_usage_events_created_at" ON "public"."usage_events" ("created_at");
CREATE INDEX "idx_usage_events_subscription_id" ON "public"."usage_events" ("subscription_id");
CREATE INDEX "idx_verification_tokens_account_id" ON "public"."verification_tokens" ("account_id");
CREATE INDEX "idx_verification_tokens_expires_at" ON "public"."verification_tokens" ("expires_at");
CREATE INDEX "idx_verification_tokens_new_email" ON "public"."verification_tokens" ("new_email") WHERE (new_email IS NOT NULL);
CREATE INDEX "idx_verification_tokens_token" ON "public"."verification_tokens" ("token");
CREATE INDEX "idx_verification_tokens_type" ON "public"."verification_tokens" ("type");

-- Add foreign key constraints
ALTER TABLE "public"."accounts" ADD CONSTRAINT "fk_accounts_subscription" FOREIGN KEY ("subscription_id") REFERENCES "public"."subscriptions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "public"."feedbacks" ADD CONSTRAINT "feedbacks_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."feedbacks" ADD CONSTRAINT "feedbacks_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."feedbacks" ADD CONSTRAINT "feedbacks_qr_code_id_fkey" FOREIGN KEY ("qr_code_id") REFERENCES "public"."qr_codes" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "public"."locations" ADD CONSTRAINT "locations_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."organizations" ADD CONSTRAINT "organizations_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."products" ADD CONSTRAINT "products_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."qr_codes" ADD CONSTRAINT "qr_codes_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."questionnaires" ADD CONSTRAINT "questionnaires_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."questionnaires" ADD CONSTRAINT "questionnaires_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."questions" ADD CONSTRAINT "fk_questions_product" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."subscription_usage" ADD CONSTRAINT "subscription_usage_subscription_id_fkey" FOREIGN KEY ("subscription_id") REFERENCES "public"."subscriptions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."subscriptions" ADD CONSTRAINT "subscriptions_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."subscriptions" ADD CONSTRAINT "subscriptions_plan_id_fkey" FOREIGN KEY ("plan_id") REFERENCES "public"."subscription_plans" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "public"."team_invitations" ADD CONSTRAINT "team_invitations_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."team_members" ADD CONSTRAINT "team_members_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."team_members" ADD CONSTRAINT "team_members_member_id_fkey" FOREIGN KEY ("member_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."usage_events" ADD CONSTRAINT "usage_events_subscription_id_fkey" FOREIGN KEY ("subscription_id") REFERENCES "public"."subscriptions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "public"."verification_tokens" ADD CONSTRAINT "fk_verification_tokens_account_id" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;

-- Add table and column comments
COMMENT ON TABLE "public"."team_members" IS 'Team members including the owner. Each account should have at least one OWNER member record';
COMMENT ON COLUMN "public"."accounts"."name" IS 'Company/organization name - used for organization accounts';
COMMENT ON COLUMN "public"."accounts"."first_name" IS 'First name for individual accounts/team members';
COMMENT ON COLUMN "public"."accounts"."last_name" IS 'Last name for individual accounts/team members';
COMMENT ON COLUMN "public"."team_members"."account_id" IS 'The organization account ID';
COMMENT ON COLUMN "public"."team_members"."member_id" IS 'The member account ID';