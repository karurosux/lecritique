-- Add deactivation_requested_at column to accounts table
ALTER TABLE accounts 
ADD COLUMN deactivation_requested_at TIMESTAMPTZ;

-- Create index for efficient querying of accounts pending deactivation
CREATE INDEX idx_accounts_deactivation_requested_at 
ON accounts(deactivation_requested_at) 
WHERE deactivation_requested_at IS NOT NULL;