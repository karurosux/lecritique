-- Create verification_tokens table
CREATE TABLE verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    token VARCHAR(128) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_verification_tokens_account_id 
        FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_verification_tokens_account_id ON verification_tokens(account_id);
CREATE INDEX idx_verification_tokens_token ON verification_tokens(token);
CREATE INDEX idx_verification_tokens_type ON verification_tokens(type);
CREATE INDEX idx_verification_tokens_expires_at ON verification_tokens(expires_at);

-- Add check constraint for token types
ALTER TABLE verification_tokens 
ADD CONSTRAINT chk_verification_tokens_type 
CHECK (type IN ('EMAIL_VERIFICATION', 'PASSWORD_RESET', 'TEAM_INVITE'));