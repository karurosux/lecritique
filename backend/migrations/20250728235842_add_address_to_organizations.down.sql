-- Recreate locations table
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

-- Add index
CREATE INDEX "idx_locations_organization_id" ON "public"."locations" ("organization_id");

-- Add foreign key constraint
ALTER TABLE "public"."locations" ADD CONSTRAINT "locations_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;

-- Remove address column from organizations table
ALTER TABLE organizations DROP COLUMN IF EXISTS address;