-- Add email_accepted_at field to track when the invitation link was clicked
ALTER TABLE team_invitations
ADD COLUMN email_accepted_at TIMESTAMP;