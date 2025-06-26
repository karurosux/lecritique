-- Remove deleted_at column from verification_tokens table
ALTER TABLE verification_tokens DROP COLUMN deleted_at;