-- Rollback: Remove first_name and last_name fields from accounts table

-- Remove comments
COMMENT ON COLUMN accounts.first_name IS NULL;
COMMENT ON COLUMN accounts.last_name IS NULL;
COMMENT ON COLUMN accounts.company_name IS NULL;

-- Drop the columns
ALTER TABLE accounts 
    DROP COLUMN IF EXISTS first_name,
    DROP COLUMN IF EXISTS last_name;