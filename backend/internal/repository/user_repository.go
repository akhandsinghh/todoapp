package repository

import (
	"context"
	"todo-app/backend/internal/db/sqlc"
)

type UserRepository interface {
	Create(ctx context.Context, name, email, hash string) (int64, error)
	ByEmail(ctx context.Context, email string) (sqlc.User, error)
	ByID(ctx context.Context, id int64) (sqlc.User, error)
	UpdatePassword(ctx context.Context, id int64, newHash string) error
}

type userRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) UserRepository {
	return &userRepository{q: q}
}
func (r *userRepository) Create(ctx context.Context, name, email, hash string) (int64, error) {
	return r.q.CreateUser(ctx, sqlc.CreateUserParams{Name: name, Email: email, PasswordHash: hash})
}
func (r *userRepository) ByEmail(ctx context.Context, email string) (sqlc.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}
func (r *userRepository) ByID(ctx context.Context, id int64) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *userRepository) UpdatePassword(ctx context.Context, id int64, newHash string) error {
	return r.q.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		PasswordHash: newHash,
		ID:           id,
	})
}
