-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetIdUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING *;


-- name: ListMessages :many
SELECT * FROM messages
ORDER BY id;

-- name: CreateMessage :one
INSERT INTO messages (user_id, content)
VALUES ($1, $2)
RETURNING *;