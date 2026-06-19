-- name: CreateGroup :execlastid
INSERT INTO task_groups (user_id, name, color) VALUES (?, ?, ?);

-- name: CreateGroupShare :exec
INSERT IGNORE INTO group_shares (group_id, user_id) VALUES (?, ?);

-- name: ListGroupsByUser :many
SELECT id, user_id, name, color, created_at, updated_at FROM task_groups WHERE user_id = ? ORDER BY name ASC;

-- name: ListAccessibleGroups :many
SELECT g.id, g.user_id, g.name, g.color, g.created_at, g.updated_at,
  CASE WHEN g.user_id = ? THEN 'creator' ELSE 'shared' END AS role
FROM task_groups g
LEFT JOIN group_shares gs ON gs.group_id = g.id AND gs.user_id = ?
WHERE g.user_id = ? OR gs.user_id IS NOT NULL
ORDER BY g.name ASC;

-- name: GetGroupByID :one
SELECT id, user_id, name, color, created_at, updated_at FROM task_groups WHERE id = ? AND user_id = ? LIMIT 1;

-- name: GetAccessibleGroupByID :one
SELECT g.id, g.user_id, g.name, g.color, g.created_at, g.updated_at,
  CASE WHEN g.user_id = ? THEN 'creator' ELSE 'shared' END AS role
FROM task_groups g
LEFT JOIN group_shares gs ON gs.group_id = g.id AND gs.user_id = ?
WHERE g.id = ? AND (g.user_id = ? OR gs.user_id IS NOT NULL)
LIMIT 1;

-- name: UpdateGroup :exec
UPDATE task_groups SET name = ?, color = ? WHERE id = ? AND user_id = ?;

-- name: DeleteGroup :exec
DELETE FROM task_groups WHERE id = ? AND user_id = ?;
