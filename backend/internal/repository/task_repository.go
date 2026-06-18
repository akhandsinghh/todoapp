package repository

import (
	"context"
	"todo-app/backend/internal/db/sqlc"
)

type TaskRepository struct{ q *sqlc.Queries }

func NewTaskRepository(q *sqlc.Queries) *TaskRepository {
	return &TaskRepository{q: q}
}
func (r *TaskRepository) Create(ctx context.Context, p sqlc.CreateTaskParams) (int64, error) {
	return r.q.CreateTask(ctx, p)
}
func (r *TaskRepository) List(ctx context.Context, p sqlc.ListTasksByUserParams) ([]sqlc.Task, error) {
	return r.q.ListTasksByUser(ctx, p)
}
func (r *TaskRepository) Get(ctx context.Context, p sqlc.GetTaskByIDParams) (sqlc.Task, error) {
	return r.q.GetTaskByID(ctx, p)
}
func (r *TaskRepository) Update(ctx context.Context, p sqlc.UpdateTaskParams) error {
	return r.q.UpdateTask(ctx, p)
}
func (r *TaskRepository) Delete(ctx context.Context, p sqlc.DeleteTaskParams) error {
	return r.q.DeleteTask(ctx, p)
}
