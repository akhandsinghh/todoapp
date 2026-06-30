package repository

import (
	"context"
	"todo-app/backend/internal/db/sqlc"
)

type TaskRepository interface {
	Create(ctx context.Context, p sqlc.CreateTaskParams) (int64, error)
	List(ctx context.Context, p sqlc.ListTasksByUserParams) ([]sqlc.Task, error)
	Count(ctx context.Context, p sqlc.CountTasksByUserParams) (int64, error)
	GetAccessibleGroup(ctx context.Context, p sqlc.GetAccessibleGroupByIDParams) (sqlc.AccessibleGroup, error)
	Get(ctx context.Context, p sqlc.GetTaskByIDParams) (sqlc.Task, error)
	Update(ctx context.Context, p sqlc.UpdateTaskParams) error
	Delete(ctx context.Context, p sqlc.DeleteTaskParams) error
}

type taskRepository struct{ q *sqlc.Queries }

func NewTaskRepository(q *sqlc.Queries) TaskRepository {
	return &taskRepository{q: q}
}
func (r *taskRepository) Create(ctx context.Context, p sqlc.CreateTaskParams) (int64, error) {
	return r.q.CreateTask(ctx, p)
}
func (r *taskRepository) List(ctx context.Context, p sqlc.ListTasksByUserParams) ([]sqlc.Task, error) {
	return r.q.ListTasksByUser(ctx, p)
}
func (r *taskRepository) Count(ctx context.Context, p sqlc.CountTasksByUserParams) (int64, error) {
	return r.q.CountTasksByUser(ctx, p)
}
func (r *taskRepository) GetAccessibleGroup(ctx context.Context, p sqlc.GetAccessibleGroupByIDParams) (sqlc.AccessibleGroup, error) {
	return r.q.GetAccessibleGroupByID(ctx, p)
}
func (r *taskRepository) Get(ctx context.Context, p sqlc.GetTaskByIDParams) (sqlc.Task, error) {
	return r.q.GetTaskByID(ctx, p)
}
func (r *taskRepository) Update(ctx context.Context, p sqlc.UpdateTaskParams) error {
	return r.q.UpdateTask(ctx, p)
}
func (r *taskRepository) Delete(ctx context.Context, p sqlc.DeleteTaskParams) error {
	return r.q.DeleteTask(ctx, p)
}
