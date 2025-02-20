-- +goose Up
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance NUMERIC(18, 2) NOT NULL DEFAULT 0.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add an index for faster UUID lookups
CREATE UNIQUE INDEX idx_wallet_uuid ON wallets (uuid);

-- +goose Down
DROP TABLE IF EXISTS wallets;