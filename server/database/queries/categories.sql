-- name: ListCategories :many
SELECT
  c.id,
  c.name,
  c.created_at,
  CAST(COALESCE(SUM(t.amount_cents), 0) AS INTEGER) AS balance_cents
FROM categories c
LEFT JOIN transactions t ON t.category_id = c.id
GROUP BY c.id
ORDER BY c.name;

-- name: GetCategory :one
SELECT
  c.id,
  c.name,
  c.created_at,
  CAST(COALESCE(SUM(t.amount_cents), 0) AS INTEGER) AS balance_cents
FROM categories c
LEFT JOIN transactions t ON t.category_id = c.id
WHERE c.id = @id
GROUP BY c.id;

-- name: CreateCategory :one
INSERT INTO categories (name)
VALUES (@name)
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = @id;
