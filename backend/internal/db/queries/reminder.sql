-- name: CreateReminder :execlastid
INSERT INTO reminders (user_id, task_id, remind_at, message) VALUES (?, ?, ?, ?);

-- name: ListRemindersByUser :many
SELECT id, user_id, task_id, remind_at, message, sent, created_at, updated_at FROM reminders WHERE user_id = ? ORDER BY remind_at ASC;

-- name: ListDueReminders :many
SELECT id, user_id, task_id, remind_at, message, sent, created_at, updated_at FROM reminders WHERE sent = FALSE AND remind_at <= ? ORDER BY remind_at ASC LIMIT ?;

-- name: MarkReminderSent :exec
UPDATE reminders SET sent = TRUE WHERE id = ?;

-- name: DeleteReminder :exec
DELETE FROM reminders WHERE id = ? AND user_id = ?;
