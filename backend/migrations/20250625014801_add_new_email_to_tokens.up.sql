-- Add new_email column to verification_tokens table for email change functionality
ALTER TABLE verification_tokens 
ADD COLUMN new_email VARCHAR(255) DEFAULT NULL;

-- Add index on new_email for better performance when checking if email is already in use
CREATE INDEX idx_verification_tokens_new_email ON verification_tokens(new_email) WHERE new_email IS NOT NULL;