-- Drop the partial unique index
DROP INDEX team_members_account_id_user_id_unique;

-- Recreate the original unique constraint
ALTER TABLE team_members ADD CONSTRAINT team_members_account_id_user_id_key UNIQUE(account_id, user_id);