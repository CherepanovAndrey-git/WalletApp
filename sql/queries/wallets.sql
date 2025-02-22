-- name: CreateWallet :one
INSERT INTO wallets (user_id) VALUES ($1)
    RETURNING id, uuid, user_id, balance_usd, balance_rub, balance_eur, created_at, updated_at;

-- name: GetWalletByUserID :one
SELECT id, uuid, user_id, balance_usd, balance_rub, balance_eur, created_at, updated_at
FROM wallets
WHERE user_id = $1 LIMIT 1;

-- name: UpdateUSDBalance :exec
UPDATE wallets
SET balance_usd = balance_usd + CAST(sqlc.arg(amount) AS NUMERIC(18,2)),
    updated_at = now()
WHERE user_id = sqlc.arg(user_id);

-- name: UpdateRUBBalance :exec
UPDATE wallets
SET balance_rub = balance_rub + CAST(sqlc.arg(amount) AS NUMERIC(18,2)),
    updated_at = now()
WHERE user_id = sqlc.arg(user_id);

-- name: UpdateEURBalance :exec
UPDATE wallets
SET balance_eur = balance_eur + CAST(sqlc.arg(amount) AS NUMERIC(18,2)),
    updated_at = now()
WHERE user_id = sqlc.arg(user_id);