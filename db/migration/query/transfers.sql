-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount, created_at) 
VALUES ($1, $2, $3, NOW())
RETURNING id, from_account_id, to_account_id, amount, created_at;

-- name: GetTransfer :one
SELECT * FROM transfers WHERE id = $1;

-- name: ListTransfersFromAccount :many
SELECT * FROM transfers 
WHERE from_account_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: ListTransfersToAccount :many
SELECT * FROM transfers 
WHERE to_account_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE 
    from_account_id = $1 OR
    to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;