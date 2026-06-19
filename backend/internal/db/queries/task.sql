-- name: CreateTask :execlastid
INSERT INTO tasks (user_id, group_id, title, description, priority, due_at) VALUES (?, ?, ?, ?, ?, ?);

-- name: ListTasksByUser :many
SELECT id, user_id, group_id, title, description, status, priority, due_at, completed_at, created_at, updated_at
FROM tasks
WHERE (
    user_id = ?
    OR (group_id IS NOT NULL AND EXISTS (
      SELECT 1
      FROM task_groups g
      LEFT JOIN group_shares gs ON gs.group_id = g.id AND gs.user_id = ?
      WHERE g.id = tasks.group_id AND (g.user_id = ? OR gs.user_id IS NOT NULL)
    ))
  )
  AND (? = '' OR status = ?)
  AND ((? = 1 AND group_id IS NULL) OR (? = 0 AND (? = 0 OR group_id = ?)))
ORDER BY
  CASE WHEN ? = 'priority' AND ? = 'asc' THEN FIELD(priority, 'low', 'medium', 'high') END ASC,
  CASE WHEN ? = 'priority' AND ? = 'desc' THEN FIELD(priority, 'high', 'medium', 'low') END ASC,
  CASE WHEN ? = 'due_at' AND ? = 'asc' THEN COALESCE(due_at, '9999-12-31') END ASC,
  CASE WHEN ? = 'due_at' AND ? = 'desc' THEN COALESCE(due_at, '1000-01-01') END DESC,
  COALESCE(due_at, created_at) ASC,
  created_at DESC
LIMIT ? OFFSET ?;

-- name: CountTasksByUser :one
SELECT COUNT(*)
FROM tasks
WHERE (
    user_id = ?
    OR (group_id IS NOT NULL AND EXISTS (
      SELECT 1
      FROM task_groups g
      LEFT JOIN group_shares gs ON gs.group_id = g.id AND gs.user_id = ?
      WHERE g.id = tasks.group_id AND (g.user_id = ? OR gs.user_id IS NOT NULL)
    ))
  )
  AND (? = '' OR status = ?)
  AND ((? = 1 AND group_id IS NULL) OR (? = 0 AND (? = 0 OR group_id = ?)));

-- name: GetTaskByID :one
SELECT id, user_id, group_id, title, description, status, priority, due_at, completed_at, created_at, updated_at
FROM tasks
WHERE id = ?
  AND (
    user_id = ?
    OR (group_id IS NOT NULL AND EXISTS (
      SELECT 1
      FROM task_groups g
      LEFT JOIN group_shares gs ON gs.group_id = g.id AND gs.user_id = ?
      WHERE g.id = tasks.group_id AND (g.user_id = ? OR gs.user_id IS NOT NULL)
    ))
  )
LIMIT 1;

-- name: UpdateTask :exec
UPDATE tasks SET group_id = ?, title = ?, description = ?, status = ?, priority = ?, due_at = ?, completed_at = ?
WHERE id = ?
  AND (
    user_id = ?
    OR (group_id IS NOT NULL AND EXISTS (
      SELECT 1
      FROM task_groups g
      LEFT JOIN group_shares gs ON gs.group_id = g.id AND gs.user_id = ?
      WHERE g.id = tasks.group_id AND (g.user_id = ? OR gs.user_id IS NOT NULL)
    ))
  );

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ?
  AND (
    user_id = ?
    OR (group_id IS NOT NULL AND EXISTS (
      SELECT 1
      FROM task_groups g
      LEFT JOIN group_shares gs ON gs.group_id = g.id AND gs.user_id = ?
      WHERE g.id = tasks.group_id AND (g.user_id = ? OR gs.user_id IS NOT NULL)
    ))
  );
