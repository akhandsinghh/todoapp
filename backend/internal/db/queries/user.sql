-- name: CreateUser :execlastid
INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?);

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = ? LIMIT 1;

-- name: GetUserByID :one
SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE id = ? LIMIT 1;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = ? WHERE id = ?;