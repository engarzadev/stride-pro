-- Revoked tokens table for logout / token invalidation support
CREATE TABLE IF NOT EXISTS revoked_tokens (
    jti         TEXT PRIMARY KEY,
    expires_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index used by the periodic cleanup job and ValidateToken lookup
CREATE INDEX IF NOT EXISTS idx_revoked_tokens_expires_at ON revoked_tokens(expires_at);
