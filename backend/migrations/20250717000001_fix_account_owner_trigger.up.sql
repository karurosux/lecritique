-- Fix the account owner trigger to use member_id instead of user_id

-- Drop the old trigger and function
DROP TRIGGER IF EXISTS create_account_owner_trigger ON accounts;
DROP FUNCTION IF EXISTS create_account_owner();

-- Create updated function that uses member_id
CREATE OR REPLACE FUNCTION create_account_owner()
RETURNS TRIGGER AS $$
BEGIN
    -- Create team member as owner with member_id = account_id (self-reference)
    INSERT INTO team_members (account_id, member_id, role, invited_by, invited_at, accepted_at)
    VALUES (
        NEW.id,
        NEW.id,  -- Owner is their own member
        'OWNER',
        NEW.id,  -- Self-invited
        NOW(),
        NOW()
    );
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for new accounts
CREATE TRIGGER create_account_owner_trigger
AFTER INSERT ON accounts
FOR EACH ROW
EXECUTE FUNCTION create_account_owner();

-- Create owner records for any accounts that don't have them
INSERT INTO team_members (account_id, member_id, role, invited_by, invited_at, accepted_at, created_at, updated_at)
SELECT 
    a.id,
    a.id,  -- Owner is their own member
    'OWNER',
    a.id,  -- Self-invited
    a.created_at,
    a.created_at,
    NOW(),
    NOW()
FROM accounts a
WHERE NOT EXISTS (
    SELECT 1 FROM team_members tm 
    WHERE tm.account_id = a.id 
    AND tm.member_id = a.id
    AND tm.role = 'OWNER'
);