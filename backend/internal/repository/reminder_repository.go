package repository

import (
	"context"
	"todo-app/backend/internal/db/sqlc"
)

type ReminderRepository interface {
	Create(ctx context.Context, p sqlc.CreateReminderParams) (int64, error)
	List(ctx context.Context, userID int64) ([]sqlc.Reminder, error)
	Due(ctx context.Context, p sqlc.ListDueRemindersParams) ([]sqlc.Reminder, error)
	MarkSent(ctx context.Context, id int64) error
	Delete(ctx context.Context, p sqlc.DeleteReminderParams) error
}

type reminderRepository struct {
	q *sqlc.Queries
}

func NewReminderRepository(q *sqlc.Queries) ReminderRepository {
	return &reminderRepository{q: q}
}
func (r *reminderRepository) Create(ctx context.Context, p sqlc.CreateReminderParams) (int64, error) {
	return r.q.CreateReminder(ctx, p)
}
func (r *reminderRepository) List(ctx context.Context, userID int64) ([]sqlc.Reminder, error) {
	return r.q.ListRemindersByUser(ctx, userID)
}
func (r *reminderRepository) Due(ctx context.Context, p sqlc.ListDueRemindersParams) ([]sqlc.Reminder, error) {
	return r.q.ListDueReminders(ctx, p)
}
func (r *reminderRepository) MarkSent(ctx context.Context, id int64) error {
	return r.q.MarkReminderSent(ctx, id)
}
func (r *reminderRepository) Delete(ctx context.Context, p sqlc.DeleteReminderParams) error {
	return r.q.DeleteReminder(ctx, p)
}
