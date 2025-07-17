-- Migration: Simplify team members by removing User table and using Account only

-- Step 1: Rename user_id to member_id in team_members
ALTER TABLE team_members 
    DROP CONSTRAINT IF EXISTS team_members_user_id_fkey;

ALTER TABLE team_members 
    RENAME COLUMN user_id TO member_id;

ALTER TABLE team_members 
    ADD CONSTRAINT team_members_member_id_fkey 
    FOREIGN KEY (member_id) REFERENCES accounts(id) ON DELETE CASCADE;

-- Step 2: Migrate existing User data to Account (if any exists)
-- This assumes users table exists and has data
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users') THEN
        -- Insert users as accounts where email doesn't already exist
        INSERT INTO accounts (id, email, password_hash, company_name, is_active, email_verified, created_at, updated_at)
        SELECT 
            u.id,
            u.email,
            u.password_hash,
            COALESCE(u.first_name || ' ' || u.last_name, u.email) as company_name,
            u.is_active,
            false as email_verified,
            u.created_at,
            u.updated_at
        FROM users u
        WHERE NOT EXISTS (
            SELECT 1 FROM accounts a WHERE a.email = u.email
        );
        
        -- Update team_members to use account IDs
        UPDATE team_members tm
        SET member_id = a.id
        FROM accounts a
        WHERE tm.member_id IN (SELECT id FROM users)
        AND a.email = (SELECT email FROM users WHERE id = tm.member_id);
    END IF;
END $$;

-- Step 3: Drop users table if it exists
DROP TABLE IF EXISTS users CASCADE;

-- Step 4: Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_team_members_member_id ON team_members(member_id);
CREATE INDEX IF NOT EXISTS idx_team_members_account_member ON team_members(account_id, member_id);

-- Add a comment to clarify the structure
COMMENT ON COLUMN team_members.account_id IS 'The organization account ID';
COMMENT ON COLUMN team_members.member_id IS 'The member account ID';
