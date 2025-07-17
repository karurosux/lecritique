-- Migration: Add first_name and last_name fields to accounts table

-- Add first_name and last_name columns to accounts table
ALTER TABLE accounts 
    ADD COLUMN IF NOT EXISTS first_name VARCHAR(255) DEFAULT '',
    ADD COLUMN IF NOT EXISTS last_name VARCHAR(255) DEFAULT '';

-- For existing accounts, we can optionally parse company_name to extract names
-- This is a simple heuristic that might help for individual accounts
UPDATE accounts 
SET 
    first_name = CASE 
        WHEN company_name LIKE '% %' THEN SPLIT_PART(company_name, ' ', 1)
        ELSE ''
    END,
    last_name = CASE 
        WHEN company_name LIKE '% %' THEN SUBSTRING(company_name FROM POSITION(' ' IN company_name) + 1)
        ELSE ''
    END
WHERE first_name = '' AND last_name = '' AND company_name != '';

-- Add comment to clarify usage
COMMENT ON COLUMN accounts.first_name IS 'First name for individual accounts/team members';
COMMENT ON COLUMN accounts.last_name IS 'Last name for individual accounts/team members';
COMMENT ON COLUMN accounts.company_name IS 'Company/organization name - used for organization accounts';