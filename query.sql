-- name: GetTransactions :many
SELECT amount, operation, description, timestamp
FROM transactions
WHERE account_id = $1
ORDER BY timestamp DESC
LIMIT 10;

-- name: GetAccount :one
SELECT balance, account_limit
FROM accounts
WHERE id = $1;

-- name: UpdateAccount :one
UPDATE accounts 
SET balance = balance + $1 
WHERE id = $2 
RETURNING balance, account_limit;

-- name: CreateTransaction :exec
INSERT INTO transactions (
  account_id, amount, operation, description
) 
VALUES ($1,$2,$3,$4);
