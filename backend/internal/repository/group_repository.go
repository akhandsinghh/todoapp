package repository

import (
	"context"
	"todo-app/backend/internal/db/sqlc"
)

type GroupRepository interface {
	Create(ctx context.Context, p sqlc.CreateGroupParams) (int64, error)
	Share(ctx context.Context, p sqlc.CreateGroupShareParams) error
	List(ctx context.Context, userID int64) ([]sqlc.AccessibleGroup, error)
	Get(ctx context.Context, p sqlc.GetGroupByIDParams) (sqlc.TaskGroup, error)
	GetAccessible(ctx context.Context, p sqlc.GetAccessibleGroupByIDParams) (sqlc.AccessibleGroup, error)
	UserByEmail(ctx context.Context, email string) (sqlc.User, error)
	Update(ctx context.Context, p sqlc.UpdateGroupParams) error
	Delete(ctx context.Context, p sqlc.DeleteGroupParams) error
}

type groupRepository struct {
	q *sqlc.Queries
}

func NewGroupRepository(q *sqlc.Queries) GroupRepository {
	return &groupRepository{q: q}
}
func (r *groupRepository) Create(ctx context.Context, p sqlc.CreateGroupParams) (int64, error) {
	return r.q.CreateGroup(ctx, p)
}
func (r *groupRepository) Share(ctx context.Context, p sqlc.CreateGroupShareParams) error {
	return r.q.CreateGroupShare(ctx, p)
}
func (r *groupRepository) List(ctx context.Context, userID int64) ([]sqlc.AccessibleGroup, error) {
	return r.q.ListAccessibleGroups(ctx, userID)
}
func (r *groupRepository) Get(ctx context.Context, p sqlc.GetGroupByIDParams) (sqlc.TaskGroup, error) {
	return r.q.GetGroupByID(ctx, p)
}
func (r *groupRepository) GetAccessible(ctx context.Context, p sqlc.GetAccessibleGroupByIDParams) (sqlc.AccessibleGroup, error) {
	return r.q.GetAccessibleGroupByID(ctx, p)
}
func (r *groupRepository) UserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}
func (r *groupRepository) Update(ctx context.Context, p sqlc.UpdateGroupParams) error {
	return r.q.UpdateGroup(ctx, p)
}
func (r *groupRepository) Delete(ctx context.Context, p sqlc.DeleteGroupParams) error {
	return r.q.DeleteGroup(ctx, p)
}
