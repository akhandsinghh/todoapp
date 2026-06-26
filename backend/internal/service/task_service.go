package service

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"todo-app/backend/internal/db/sqlc"
	apperr "todo-app/backend/internal/errors"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/repository"
	"todo-app/backend/internal/util"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func taskDTO(t sqlc.Task) model.TaskDTO {
	var gid *int64
	if t.GroupID.Valid {
		v := t.GroupID.Int64
		gid = &v
	}
	due := timePtr(t.DueAt)
	done := timePtr(t.CompletedAt)
	return model.TaskDTO{
		ID:          t.ID,
		UserID:      t.UserID,
		GroupID:     gid,
		Title:       t.Title,
		Description: t.Description.String,
		Status:      string(t.Status),
		Priority:    string(t.Priority),
		DueAt:       due,
		CompletedAt: done,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}
}

func timePtr(t sql.NullTime) *string {
	if !t.Valid {
		return nil
	}
	v := t.Time.Format(time.RFC3339)
	return &v
}

func parseTime(v string) (sql.NullTime, error) {
	if v == "" {
		return sql.NullTime{}, nil
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}

func nullGroup(v *int64) sql.NullInt64 {
	if v == nil || *v == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *v, Valid: true}
}

func (s *TaskService) ensureGroupAccess(ctx context.Context, userID int64, groupID *int64) error {
	if groupID == nil || *groupID == 0 {
		return nil
	}
	if _, err := s.repo.GetAccessibleGroup(ctx, sqlc.GetAccessibleGroupByIDParams{ID: *groupID, UserID: userID}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperr.Forbidden("group is not accessible")
		}
		return apperr.Forbidden("group is not accessible")
	}
	return nil
}

func (s *TaskService) Create(ctx context.Context, userID int64, req model.TaskRequest) (model.TaskDTO, error) {
	if !util.Required(req.Title) {
		return model.TaskDTO{}, apperr.BadRequest("title is required")
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}
	if !util.ValidPriority(req.Priority) {
		return model.TaskDTO{}, apperr.BadRequest("invalid priority")
	}
	if err := s.ensureGroupAccess(ctx, userID, req.GroupID); err != nil {
		return model.TaskDTO{}, err
	}
	due, err := parseTime(req.DueAt)
	if err != nil {
		return model.TaskDTO{}, apperr.BadRequest("due_at must be RFC3339")
	}
	id, err := s.repo.Create(ctx, sqlc.CreateTaskParams{
		UserID:      userID,
		GroupID:     nullGroup(req.GroupID),
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Priority:    sqlc.TasksPriority(req.Priority),
		DueAt:       due,
	})
	if err != nil {
		return model.TaskDTO{}, apperr.Internal("failed to create task")
	}
	t, err := s.repo.Get(ctx, sqlc.GetTaskByIDParams{ID: id, UserID: userID})
	if err != nil {
		return model.TaskDTO{}, err
	}
	return taskDTO(t), nil
}

func normalizeSort(sortBy, sortOrder string) (string, string) {
	if sortBy != "priority" && sortBy != "due_at" {
		sortBy = ""
	}
	if sortOrder != "desc" {
		sortOrder = "asc"
	}
	return sortBy, sortOrder
}

func (s *TaskService) List(ctx context.Context, userID int64, status string, groupID int64, ungrouped bool, limit, offset int32, sortBy, sortOrder string) (model.TaskListResponse, error) {
	sortBy, sortOrder = normalizeSort(sortBy, sortOrder)
	var ungroupedFlag int32
	if ungrouped {
		ungroupedFlag = 1
		groupID = 0
	}
	countParams := sqlc.CountTasksByUserParams{
		UserID:     userID,
		Column3:    status,
		Status:     sqlc.TasksStatus(status),
		Ungrouped1: ungroupedFlag,
		Ungrouped2: ungroupedFlag,
		Column7:    groupID,
		GroupID: sql.NullInt64{
			Int64: groupID,
			Valid: groupID != 0,
		},
	}
	total, err := s.repo.Count(ctx, countParams)
	if err != nil {
		return model.TaskListResponse{}, apperr.Internal("failed to count tasks")
	}
	ts, err := s.repo.List(ctx, sqlc.ListTasksByUserParams{
		UserID: userID,

		// Map the raw strings to the dynamically generated sqlc interface fields
		Column3: status,
		Status:  sqlc.TasksStatus(status),

		Ungrouped1: ungroupedFlag,
		Ungrouped2: ungroupedFlag,
		Column7:    groupID,
		GroupID: sql.NullInt64{
			Int64: groupID,
			Valid: groupID != 0,
		},

		SortBy1: sortBy,
		Order1:  sortOrder,
		SortBy2: sortBy,
		Order2:  sortOrder,
		SortBy3: sortBy,
		Order3:  sortOrder,
		SortBy4: sortBy,
		Order4:  sortOrder,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return model.TaskListResponse{}, apperr.Internal("failed to fetch tasks")
	}

	out := make([]model.TaskDTO, 0, len(ts))
	for _, t := range ts {
		out = append(out, taskDTO(t))
	}
	return model.TaskListResponse{Items: out, Total: total}, nil
}

func (s *TaskService) Update(ctx context.Context, userID, id int64, req model.TaskRequest) (model.TaskDTO, error) {
	cur, err := s.repo.Get(ctx, sqlc.GetTaskByIDParams{ID: id, UserID: userID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.TaskDTO{}, apperr.NotFound("task not found")
		}
		return model.TaskDTO{}, apperr.Internal("failed to fetch task")
	}
	if req.Title == "" {
		req.Title = cur.Title
	}
	if req.Status == "" {
		req.Status = string(cur.Status)
	}
	if req.Priority == "" {
		req.Priority = string(cur.Priority)
	}
	if !util.ValidStatus(req.Status) {
		return model.TaskDTO{}, apperr.BadRequest("invalid status")
	}
	if !util.ValidPriority(req.Priority) {
		return model.TaskDTO{}, apperr.BadRequest("invalid priority")
	}
	if err := s.ensureGroupAccess(ctx, userID, req.GroupID); err != nil {
		return model.TaskDTO{}, err
	}
	due, err := parseTime(req.DueAt)
	if err != nil {
		return model.TaskDTO{}, apperr.BadRequest("due_at must be RFC3339")
	}
	if req.DueAt == "" {
		due = cur.DueAt
	}
	done := cur.CompletedAt
	if req.Status == "completed" && !done.Valid {
		done = sql.NullTime{Time: time.Now(), Valid: true}
	}
	if req.Status == "pending" {
		done = sql.NullTime{}
	}
	desc := sql.NullString{String: req.Description, Valid: req.Description != ""}
	if req.Description == "" {
		desc = cur.Description
	}
	if err := s.repo.Update(ctx, sqlc.UpdateTaskParams{
		ID:          id,
		UserID:      userID,
		GroupID:     nullGroup(req.GroupID),
		Title:       req.Title,
		Description: desc,
		Status:      sqlc.TasksStatus(req.Status),
		Priority:    sqlc.TasksPriority(req.Priority),
		DueAt:       due,
		CompletedAt: done,
	}); err != nil {
		return model.TaskDTO{}, apperr.Internal("failed to update task")
	}
	t, err := s.repo.Get(ctx, sqlc.GetTaskByIDParams{ID: id, UserID: userID})
	if err != nil {
		return model.TaskDTO{}, apperr.Internal("failed to fetch updated task")
	}
	return taskDTO(t), nil
}

func (s *TaskService) Delete(ctx context.Context, userID, id int64) error {
	err := s.repo.Delete(ctx, sqlc.DeleteTaskParams{ID: id, UserID: userID})
	if err != nil {
		return apperr.Internal("failed to delete task")
	}
	return nil
}
