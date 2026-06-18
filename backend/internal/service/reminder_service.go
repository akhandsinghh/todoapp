package service

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"todo-app/backend/internal/db/sqlc"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/repository"
)

type ReminderService struct {
	repo  *repository.ReminderRepository
	tasks *repository.TaskRepository
}

func NewReminderService(repo *repository.ReminderRepository, tasks *repository.TaskRepository) *ReminderService {
	return &ReminderService{repo: repo, tasks: tasks}
}

func reminderDTO(r sqlc.Reminder) model.ReminderDTO {
	return model.ReminderDTO{
		ID:        r.ID,
		UserID:    r.UserID,
		TaskID:    r.TaskID,
		RemindAt:  r.RemindAt.Format(timeLayout),
		Message:   r.Message.String,
		Sent:      r.Sent,
		CreatedAt: r.CreatedAt.Format(timeLayout),
	}
}

func (s *ReminderService) Create(ctx context.Context, userID int64, req model.ReminderRequest) (model.ReminderDTO, error) {
	if req.TaskID == 0 || req.RemindAt == "" {
		return model.ReminderDTO{}, errors.New("task_id and remind_at are required")
	}
	if _, err := s.tasks.Get(ctx, sqlc.GetTaskByIDParams{ID: req.TaskID, UserID: userID}); err != nil {
		return model.ReminderDTO{}, err
	}
	at, err := time.Parse(time.RFC3339, req.RemindAt)
	if err != nil {
		return model.ReminderDTO{}, errors.New("remind_at must be RFC3339")
	}
	id, err := s.repo.Create(ctx, sqlc.CreateReminderParams{
		UserID:   userID,
		TaskID:   req.TaskID,
		RemindAt: at,
		Message:  sql.NullString{String: req.Message, Valid: req.Message != ""},
	})
	if err != nil {
		return model.ReminderDTO{}, err
	}
	rs, err := s.repo.List(ctx, userID)
	if err != nil {
		return model.ReminderDTO{}, err
	}
	for _, r := range rs {
		if r.ID == id {
			return reminderDTO(r), nil
		}
	}
	return model.ReminderDTO{}, sql.ErrNoRows
}

func (s *ReminderService) List(ctx context.Context, userID int64) ([]model.ReminderDTO, error) {
	rs, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]model.ReminderDTO, 0, len(rs))
	for _, r := range rs {
		out = append(out, reminderDTO(r))
	}
	return out, nil
}

func (s *ReminderService) Delete(ctx context.Context, userID, id int64) error {
	return s.repo.Delete(ctx, sqlc.DeleteReminderParams{ID: id, UserID: userID})
}

func (s *ReminderService) Due(ctx context.Context, limit int32) ([]sqlc.Reminder, error) {
	return s.repo.Due(ctx, sqlc.ListDueRemindersParams{RemindAt: time.Now(), Limit: limit})
}

func (s *ReminderService) MarkSent(ctx context.Context, id int64) error {
	return s.repo.MarkSent(ctx, id)
}
