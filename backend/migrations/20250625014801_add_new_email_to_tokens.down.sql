-- Remove index
DROP INDEX IF EXISTS idx_verification_tokens_new_email;

-- Remove new_email column from verification_tokens table
ALTER TABLE verification_tokens 
DROP COLUMN IF EXISTS new_email;