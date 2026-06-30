package test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"todo-app/backend/internal/controller/dto"
	"todo-app/backend/internal/db/sqlc"
	"todo-app/backend/internal/service"
	"todo-app/backend/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupReminderTestHelper initializes both mocks and the service.
func setupReminderTestHelper(t *testing.T) (*mocks.ReminderRepository, *mocks.TaskRepository, *service.ReminderService) {
	mockReminderRepo := mocks.NewReminderRepository(t)
	mockTaskRepo := mocks.NewTaskRepository(t)

	svc := service.NewReminderService(mockReminderRepo, mockTaskRepo)
	return mockReminderRepo, mockTaskRepo, svc
}

func TestReminderService_Create(t *testing.T) {
	mockRepo, mockTaskRepo, svc := setupReminderTestHelper(t)
	ctx := context.Background()
	userID := int64(1)
	taskID := int64(10)
	reminderID := int64(100)

	futureTime := time.Now().Add(24 * time.Hour)
	timeString := futureTime.Format(time.RFC3339)

	req := dto.ReminderRequest{
		TaskID:   taskID,
		RemindAt: timeString,
		Message:  "Don't forget this!",
	}

	// 1. Expect Task verification (Get task by ID)
	mockTaskRepo.On("Get", ctx, sqlc.GetTaskByIDParams{ID: taskID, UserID: userID}).Return(sqlc.Task{ID: taskID}, nil)

	// 2. Expect Reminder Create
	mockRepo.On("Create", ctx, mock.AnythingOfType("sqlc.CreateReminderParams")).Return(reminderID, nil)

	// 3. Expect List to be called to fetch the newly created reminder's DTO details
	mockReminders := []sqlc.Reminder{
		{
			ID:        reminderID,
			UserID:    userID,
			TaskID:    taskID,
			RemindAt:  futureTime,
			Message:   sql.NullString{String: req.Message, Valid: true},
			Sent:      false,
			CreatedAt: time.Now(),
		},
	}
	mockRepo.On("List", ctx, userID).Return(mockReminders, nil)

	// Execute
	resp, err := svc.Create(ctx, userID, req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, reminderID, resp.ID)
	assert.Equal(t, taskID, resp.TaskID)
	assert.Equal(t, "Don't forget this!", resp.Message)
	assert.False(t, resp.Sent)

	mockRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestReminderService_List(t *testing.T) {
	mockRepo, _, svc := setupReminderTestHelper(t)
	ctx := context.Background()
	userID := int64(1)

	mockReminders := []sqlc.Reminder{
		{ID: 100, UserID: userID, TaskID: 10, Message: sql.NullString{String: "First", Valid: true}, CreatedAt: time.Now(), RemindAt: time.Now()},
		{ID: 101, UserID: userID, TaskID: 11, Message: sql.NullString{String: "Second", Valid: true}, CreatedAt: time.Now(), RemindAt: time.Now()},
	}

	// 1. Expect List
	mockRepo.On("List", ctx, userID).Return(mockReminders, nil)

	// Execute
	resp, err := svc.List(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "First", resp[0].Message)
	assert.Equal(t, "Second", resp[1].Message)

	mockRepo.AssertExpectations(t)
}

func TestReminderService_Delete(t *testing.T) {
	mockRepo, _, svc := setupReminderTestHelper(t)
	ctx := context.Background()
	userID := int64(1)
	reminderID := int64(100)

	// 1. Expect Delete
	mockRepo.On("Delete", ctx, sqlc.DeleteReminderParams{ID: reminderID, UserID: userID}).Return(nil)

	// Execute
	err := svc.Delete(ctx, userID, reminderID)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReminderService_Due(t *testing.T) {
	mockRepo, _, svc := setupReminderTestHelper(t)
	ctx := context.Background()
	limit := int32(10)

	mockDueReminders := []sqlc.Reminder{
		{ID: 100, TaskID: 10},
		{ID: 101, TaskID: 11},
	}

	// 1. Expect Due to be called. We use mock.AnythingOfType because the argument includes time.Now()
	mockRepo.On("Due", ctx, mock.AnythingOfType("sqlc.ListDueRemindersParams")).Return(mockDueReminders, nil)

	// Execute
	resp, err := svc.Due(ctx, limit)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, resp, 2)

	mockRepo.AssertExpectations(t)
}

func TestReminderService_MarkSent(t *testing.T) {
	mockRepo, _, svc := setupReminderTestHelper(t)
	ctx := context.Background()
	reminderID := int64(100)

	// 1. Expect MarkSent
	mockRepo.On("MarkSent", ctx, reminderID).Return(nil)

	// Execute
	err := svc.MarkSent(ctx, reminderID)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
