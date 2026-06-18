-- name: CreateGroup :execlastid
INSERT INTO task_groups (user_id, name, color) VALUES (?, ?, ?);

-- name: ListGroupsByUser :many
SELECT id, user_id, name, color, created_at, updated_at FROM task_groups WHERE user_id = ? ORDER BY name ASC;

-- name: GetGroupByID :one
SELECT id, user_id, name, color, created_at, updated_at FROM task_groups WHERE id = ? AND user_id = ? LIMIT 1;

-- name: UpdateGroup :exec
UPDATE task_groups SET name = ?, color = ? WHERE id = ? AND user_id = ?;

-- name: DeleteGroup :exec
DELETE FROM task_groups WHERE id = ? AND user_id = ?;
