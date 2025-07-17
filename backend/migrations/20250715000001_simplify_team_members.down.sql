-- Rollback: Restore User table and team_members structure

-- Step 1: Create users table (if it was dropped)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Step 2: Rename member_id back to user_id in team_members
ALTER TABLE team_members 
    DROP CONSTRAINT IF EXISTS team_members_member_id_fkey;

ALTER TABLE team_members 
    RENAME COLUMN member_id TO user_id;

ALTER TABLE team_members 
    ADD CONSTRAINT team_members_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Step 3: Migrate data back from accounts to users (for non-organization accounts)
INSERT INTO users (id, email, password_hash, first_name, last_name, is_active, created_at, updated_at)
SELECT 
    tm.user_id as id,
    a.email,
    a.password_hash,
    a.first_name,
    a.last_name,
    a.is_active,
    a.created_at,
    a.updated_at
FROM team_members tm
JOIN accounts a ON tm.user_id = a.id
WHERE tm.role != 'OWNER'
ON CONFLICT (id) DO NOTHING;

-- Step 4: Drop indexes
DROP INDEX IF EXISTS idx_team_members_member_id;
DROP INDEX IF EXISTS idx_team_members_account_member;

-- Step 5: Remove comments
COMMENT ON COLUMN team_members.account_id IS NULL;
COMMENT ON COLUMN team_members.user_id IS NULL;