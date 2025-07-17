-- Remove email_accepted_at field
ALTER TABLE team_invitations
DROP COLUMN email_accepted_at;