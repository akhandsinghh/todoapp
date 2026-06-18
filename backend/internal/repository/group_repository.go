package repository

import (
	"context"
	"todo-app/backend/internal/db/sqlc"
)

type GroupRepository struct {
	q *sqlc.Queries
}

func NewGroupRepository(q *sqlc.Queries) *GroupRepository {
	return &GroupRepository{q: q}
}
func (r *GroupRepository) Create(ctx context.Context, p sqlc.CreateGroupParams) (int64, error) {
	return r.q.CreateGroup(ctx, p)
}
func (r *GroupRepository) List(ctx context.Context, userID int64) ([]sqlc.TaskGroup, error) {
	return r.q.ListGroupsByUser(ctx, userID)
}
func (r *GroupRepository) Get(ctx context.Context, p sqlc.GetGroupByIDParams) (sqlc.TaskGroup, error) {
	return r.q.GetGroupByID(ctx, p)
}
func (r *GroupRepository) Update(ctx context.Context, p sqlc.UpdateGroupParams) error {
	return r.q.UpdateGroup(ctx, p)
}
func (r *GroupRepository) Delete(ctx context.Context, p sqlc.DeleteGroupParams) error {
	return r.q.DeleteGroup(ctx, p)
}
