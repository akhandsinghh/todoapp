package repository

import (
	"context"
	"todo-app/backend/internal/db/sqlc"
)

type ReminderRepository struct {
	q *sqlc.Queries
}

func NewReminderRepository(q *sqlc.Queries) *ReminderRepository {
	return &ReminderRepository{q: q}
}
func (r *ReminderRepository) Create(ctx context.Context, p sqlc.CreateReminderParams) (int64, error) {
	return r.q.CreateReminder(ctx, p)
}
func (r *ReminderRepository) List(ctx context.Context, userID int64) ([]sqlc.Reminder, error) {
	return r.q.ListRemindersByUser(ctx, userID)
}
func (r *ReminderRepository) Due(ctx context.Context, p sqlc.ListDueRemindersParams) ([]sqlc.Reminder, error) {
	return r.q.ListDueReminders(ctx, p)
}
func (r *ReminderRepository) MarkSent(ctx context.Context, id int64) error {
	return r.q.MarkReminderSent(ctx, id)
}
func (r *ReminderRepository) Delete(ctx context.Context, p sqlc.DeleteReminderParams) error {
	return r.q.DeleteReminder(ctx, p)
}
