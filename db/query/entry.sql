-- name: CreateEntry :one
INSERT INTO entries (
  account_id, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatEntry :one
UPDATE entries
  SET account_id = $2, amount = $3
WHERE id = $1
RETURNING *;

-- name: DeleteEntriy :exec
DELETE FROM entries
WHERE id = $1;