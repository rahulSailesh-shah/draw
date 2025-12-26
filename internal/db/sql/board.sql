-- name: CreateBoard :one
INSERT INTO "board" (name, owner_id) VALUES ($1, $2) RETURNING *;

-- name: GetBoardByID :one
SELECT * FROM "board" WHERE id = $1 AND owner_id = $2;

-- name: UpdateBoard :one
UPDATE "board" SET name = $2, elements = $3 WHERE id = $1 AND owner_id = $4 RETURNING *;

-- name: DeleteBoard :exec
DELETE FROM "board" WHERE id = $1 AND owner_id = $2;