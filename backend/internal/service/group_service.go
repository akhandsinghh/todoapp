package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"todo-app/backend/internal/db/sqlc"
	apperr "todo-app/backend/internal/errors"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/repository"
	"todo-app/backend/internal/util"
)

type GroupService struct {
	repo *repository.GroupRepository
}

func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func groupDTO(g sqlc.TaskGroup) model.GroupDTO {
	return model.GroupDTO{
		ID:        g.ID,
		UserID:    g.UserID,
		Name:      g.Name,
		Color:     g.Color,
		Role:      "creator",
		CreatedAt: g.CreatedAt.Format(timeLayout),
		UpdatedAt: g.UpdatedAt.Format(timeLayout),
	}
}

func accessibleGroupDTO(g sqlc.AccessibleGroup) model.GroupDTO {
	return model.GroupDTO{
		ID:        g.ID,
		UserID:    g.UserID,
		Name:      g.Name,
		Color:     g.Color,
		Role:      g.Role,
		CreatedAt: g.CreatedAt.Format(timeLayout),
		UpdatedAt: g.UpdatedAt.Format(timeLayout),
	}
}

const timeLayout = "2006-01-02T15:04:05Z07:00"

func (s *GroupService) Create(ctx context.Context, userID int64, req model.GroupRequest) (model.GroupDTO, error) {
	if !util.Required(req.Name) {
		return model.GroupDTO{}, apperr.BadRequest("group name is required")
	}
	if req.Color == "" {
		req.Color = "#4f46e5"
	}
	id, err := s.repo.Create(ctx, sqlc.CreateGroupParams{UserID: userID, Name: req.Name, Color: req.Color})
	if err != nil {
		return model.GroupDTO{}, apperr.Internal("failed to create group")
	}
	g, err := s.repo.Get(ctx, sqlc.GetGroupByIDParams{ID: id, UserID: userID})
	if err != nil {
		return model.GroupDTO{}, apperr.Internal("failed to fetch group")
	}
	return groupDTO(g), nil
}

func (s *GroupService) List(ctx context.Context, userID int64) ([]model.GroupDTO, error) {
	gs, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, apperr.Internal("failed to fetch groups")
	}
	out := make([]model.GroupDTO, 0, len(gs))
	for _, g := range gs {
		out = append(out, accessibleGroupDTO(g))
	}
	return out, nil
}

func (s *GroupService) Share(ctx context.Context, userID, id int64, req model.ShareGroupRequest) error {
	email := strings.TrimSpace(req.Email)
	if email == "" {
		return apperr.BadRequest("email is required")
	}
	if _, err := s.repo.Get(ctx, sqlc.GetGroupByIDParams{ID: id, UserID: userID}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperr.Forbidden("only the creator can share this group")
		}
		return apperr.Internal("failed to fetch group")
	}
	user, err := s.repo.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperr.NotFound("user not found")
		}
		return apperr.Internal("failed to fetch user")
	}
	if user.ID == userID {
		return apperr.BadRequest("cannot share a group with yourself")
	}
	err = s.repo.Share(ctx, sqlc.CreateGroupShareParams{GroupID: id, UserID: user.ID})
	if err != nil {
		return apperr.Internal("failed to share group")
	}
	return nil
}

func (s *GroupService) Update(ctx context.Context, userID, id int64, req model.GroupRequest) error {
	if req.Color == "" {
		req.Color = "#4f46e5"
	}
	err := s.repo.Update(ctx, sqlc.UpdateGroupParams{ID: id, UserID: userID, Name: req.Name, Color: req.Color})
	if err != nil {
		return apperr.Internal("failed to update group")
	}
	return nil
}

func (s *GroupService) Delete(ctx context.Context, userID, id int64) error {
	err := s.repo.Delete(ctx, sqlc.DeleteGroupParams{ID: id, UserID: userID})
	if err != nil {
		return apperr.Internal("failed to delete group")
	}
	return nil
}
