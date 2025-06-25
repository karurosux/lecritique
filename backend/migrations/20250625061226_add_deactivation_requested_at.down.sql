-- Drop index
DROP INDEX IF EXISTS idx_accounts_deactivation_requested_at;

-- Remove deactivation_requested_at column from accounts table
ALTER TABLE accounts 
DROP COLUMN IF EXISTS deactivation_requested_at;