-- name: CreateWallet :one
INSERT INTO wallets (
    user_id,
    balance
) VALUES (
    $1,
    '0.00'
) RETURNING *;

-- name: GetWalletByUserID :one
SELECT * FROM wallets 
WHERE user_id = $1 LIMIT 1;

-- name: UpdateWalletBalance :exec
UPDATE wallets
SET 
    balance = CAST(balance AS DECIMAL(18,2)) + CAST(sqlc.arg(Amount) AS DECIMAL(18,2)),
    updated_at = now()
WHERE user_id = sqlc.arg(UserID);

-- name: GetWalletBalance :one
SELECT balance FROM wallets
WHERE user_id = $1;


ALTER TABLE wallets ADD COLUMN usd_balance NUMERIC(18,2) DEFAULT 0.0;
ALTER TABLE wallets ADD COLUMN eur_balance NUMERIC(18,2) DEFAULT 0.0;
ALTER TABLE wallets ADD COLUMN rub_balance NUMERIC(18,2) DEFAULT 0.0;