package service

import (
	"context"
	"errors"
	"todo-app/backend/internal/db/sqlc"
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
		CreatedAt: g.CreatedAt.Format(timeLayout),
		UpdatedAt: g.UpdatedAt.Format(timeLayout),
	}
}

const timeLayout = "2006-01-02T15:04:05Z07:00"

func (s *GroupService) Create(ctx context.Context, userID int64, req model.GroupRequest) (model.GroupDTO, error) {
	if !util.Required(req.Name) {
		return model.GroupDTO{}, errors.New("group name is required")
	}
	if req.Color == "" {
		req.Color = "#4f46e5"
	}
	id, err := s.repo.Create(ctx, sqlc.CreateGroupParams{UserID: userID, Name: req.Name, Color: req.Color})
	if err != nil {
		return model.GroupDTO{}, err
	}
	g, err := s.repo.Get(ctx, sqlc.GetGroupByIDParams{ID: id, UserID: userID})
	if err != nil {
		return model.GroupDTO{}, err
	}
	return groupDTO(g), nil
}

func (s *GroupService) List(ctx context.Context, userID int64) ([]model.GroupDTO, error) {
	gs, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]model.GroupDTO, 0, len(gs))
	for _, g := range gs {
		out = append(out, groupDTO(g))
	}
	return out, nil
}

func (s *GroupService) Update(ctx context.Context, userID, id int64, req model.GroupRequest) error {
	if req.Color == "" {
		req.Color = "#4f46e5"
	}
	return s.repo.Update(ctx, sqlc.UpdateGroupParams{ID: id, UserID: userID, Name: req.Name, Color: req.Color})
}

func (s *GroupService) Delete(ctx context.Context, userID, id int64) error {
	return s.repo.Delete(ctx, sqlc.DeleteGroupParams{ID: id, UserID: userID})
}
