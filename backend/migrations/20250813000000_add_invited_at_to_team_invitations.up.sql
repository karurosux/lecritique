-- Add invited_at column to team_invitations table
ALTER TABLE "public"."team_invitations" ADD COLUMN "invited_at" timestamptz NOT NULL DEFAULT now();