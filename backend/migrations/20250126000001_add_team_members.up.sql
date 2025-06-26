-- Only create the trigger and function since tables already exist in initial migration

-- Create the account owner when an account is created
CREATE OR REPLACE FUNCTION create_account_owner()
RETURNS TRIGGER AS $$
BEGIN
    -- Create a user with the same email as the account
    INSERT INTO users (email, password_hash, is_active)
    VALUES (NEW.email, NEW.password_hash, NEW.is_active)
    ON CONFLICT (email) DO UPDATE SET
        password_hash = EXCLUDED.password_hash,
        is_active = EXCLUDED.is_active;
    
    -- Create team member as owner
    INSERT INTO team_members (account_id, user_id, role, invited_by, invited_at, accepted_at)
    SELECT 
        NEW.id,
        u.id,
        'OWNER',
        u.id,
        NOW(),
        NOW()
    FROM users u
    WHERE u.email = NEW.email;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for new accounts
CREATE TRIGGER create_account_owner_trigger
AFTER INSERT ON accounts
FOR EACH ROW
EXECUTE FUNCTION create_account_owner();

-- Migrate existing accounts to have owners
INSERT INTO users (email, password_hash, is_active, created_at, updated_at)
SELECT DISTINCT email, password_hash, is_active, created_at, updated_at
FROM accounts
ON CONFLICT (email) DO NOTHING;

INSERT INTO team_members (account_id, user_id, role, invited_by, invited_at, accepted_at, created_at, updated_at)
SELECT 
    a.id,
    u.id,
    'OWNER',
    u.id,
    a.created_at,
    a.created_at,
    a.created_at,
    a.updated_at
FROM accounts a
JOIN users u ON u.email = a.email
WHERE NOT EXISTS (
    SELECT 1 FROM team_members tm 
    WHERE tm.account_id = a.id AND tm.user_id = u.id
);