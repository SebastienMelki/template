-- name: CreateExample :one
INSERT INTO examples (id, name, description, created_at, updated_at)
VALUES (@id, @name, @description, NOW(), NOW())
RETURNING id, name, description, created_at, updated_at;

-- name: GetExampleByID :one
SELECT id, name, description, created_at, updated_at
FROM examples
WHERE id = @id;

-- name: ListExamples :many
SELECT id, name, description, created_at, updated_at
FROM examples
ORDER BY created_at DESC
LIMIT @result_limit OFFSET @result_offset;

-- name: CountExamples :one
SELECT COUNT(*) FROM examples;

-- name: DeleteExample :exec
DELETE FROM examples WHERE id = @id;
