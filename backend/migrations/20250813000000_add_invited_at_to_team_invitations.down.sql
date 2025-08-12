-- Remove invited_at column from team_invitations table
ALTER TABLE "public"."team_invitations" DROP COLUMN "invited_at";