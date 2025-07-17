-- Migration: Create team member record for account owners

-- For each account that doesn't have any team members, create an owner record
-- This ensures the owner appears in the team members list
INSERT INTO team_members (id, account_id, member_id, role, invited_by, invited_at, accepted_at, created_at, updated_at)
SELECT 
    gen_random_uuid() as id,
    a.id as account_id,
    a.id as member_id,  -- Owner is their own member
    'OWNER' as role,
    a.id as invited_by,  -- Self-invited
    a.created_at as invited_at,
    a.created_at as accepted_at,
    NOW() as created_at,
    NOW() as updated_at
FROM accounts a
WHERE NOT EXISTS (
    SELECT 1 
    FROM team_members tm 
    WHERE tm.account_id = a.id 
    AND tm.role = 'OWNER'
);

-- Add comment
COMMENT ON TABLE team_members IS 'Team members including the owner. Each account should have at least one OWNER member record';