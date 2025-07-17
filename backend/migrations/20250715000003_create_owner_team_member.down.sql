-- Rollback: Remove auto-created owner team member records

-- Delete owner team member records where member_id equals account_id
-- This removes the self-referential owner records created by the migration
DELETE FROM team_members 
WHERE role = 'OWNER' 
AND member_id = account_id
AND invited_by = account_id;

-- Remove comment
COMMENT ON TABLE team_members IS NULL;