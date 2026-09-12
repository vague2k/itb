-- name: CreateTransaction :one
INSERT INTO transactions (category_id, amount_cents, note)
VALUES (@category_id, @amount_cents, @note)
RETURNING *;

-- name: ListTransactionsByCategory :many
SELECT * FROM transactions
WHERE category_id = @category_id
ORDER BY created_at DESC, id DESC;
