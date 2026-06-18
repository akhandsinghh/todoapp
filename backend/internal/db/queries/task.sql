-- name: CreateTask :execlastid
INSERT INTO tasks (user_id, group_id, title, description, priority, due_at) VALUES (?, ?, ?, ?, ?, ?);

-- name: ListTasksByUser :many
SELECT id, user_id, group_id, title, description, status, priority, due_at, completed_at, created_at, updated_at
FROM tasks
WHERE user_id = ?
  AND (? = '' OR status = ?)
  AND (? = 0 OR group_id = ?)
ORDER BY COALESCE(due_at, created_at) ASC, created_at DESC
LIMIT ? OFFSET ?;

-- name: GetTaskByID :one
SELECT id, user_id, group_id, title, description, status, priority, due_at, completed_at, created_at, updated_at
FROM tasks WHERE id = ? AND user_id = ? LIMIT 1;

-- name: UpdateTask :exec
UPDATE tasks SET group_id = ?, title = ?, description = ?, status = ?, priority = ?, due_at = ?, completed_at = ? WHERE id = ? AND user_id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = ? AND user_id = ?;
