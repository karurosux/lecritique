-- Add deleted_at column to verification_tokens table
ALTER TABLE verification_tokens ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;