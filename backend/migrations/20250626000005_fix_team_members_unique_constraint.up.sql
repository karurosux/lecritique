-- Drop the existing unique constraint
ALTER TABLE team_members DROP CONSTRAINT team_members_account_id_user_id_key;

-- Create a partial unique index that only applies to non-deleted records
CREATE UNIQUE INDEX team_members_account_id_user_id_unique 
ON team_members(account_id, user_id) 
WHERE deleted_at IS NULL;