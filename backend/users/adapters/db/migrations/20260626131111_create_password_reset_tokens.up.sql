CREATE TABLE users.password_reset_tokens (
    token_hash VARCHAR(256) PRIMARY KEY,
    user_uuid UUID NOT NULL REFERENCES users.users(user_uuid) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_password_reset_tokens_user_uuid ON users.password_reset_tokens(user_uuid);
