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

func setupTaskTestHelper(t *testing.T) (*mocks.TaskRepository, *service.TaskService) {
	mockRepo := mocks.NewTaskRepository(t)
	svc := service.NewTaskService(mockRepo)
	return mockRepo, svc
}

func TestTaskService_Create(t *testing.T) {
	mockRepo, svc := setupTaskTestHelper(t)
	ctx := context.Background()
	userID := int64(1)
	groupID := int64(5)
	taskID := int64(10)

	req := dto.TaskRequest{
		Title:       "Test Task",
		Description: "A simple task",
		Priority:    "high",
		GroupID:     &groupID,
	}

	mockTask := sqlc.Task{
		ID:        taskID,
		UserID:    userID,
		GroupID:   sql.NullInt64{Int64: groupID, Valid: true},
		Title:     req.Title,
		Priority:  sqlc.TasksPriority(req.Priority),
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 1. Expect Group Access Check (Because GroupID is provided)
	mockRepo.On("GetAccessibleGroup", ctx, sqlc.GetAccessibleGroupByIDParams{ID: groupID, UserID: userID}).
		Return(sqlc.AccessibleGroup{ID: groupID}, nil)

	// 2. Expect Create
	mockRepo.On("Create", ctx, mock.AnythingOfType("sqlc.CreateTaskParams")).Return(taskID, nil)

	// 3. Expect Get to fetch final DTO data
	mockRepo.On("Get", ctx, sqlc.GetTaskByIDParams{ID: taskID, UserID: userID}).Return(mockTask, nil)

	// Execute
	resp, err := svc.Create(ctx, userID, req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, taskID, resp.ID)
	assert.Equal(t, "Test Task", resp.Title)
	assert.Equal(t, "high", resp.Priority)
	assert.Equal(t, groupID, *resp.GroupID)

	mockRepo.AssertExpectations(t)
}

func TestTaskService_List(t *testing.T) {
	mockRepo, svc := setupTaskTestHelper(t)
	ctx := context.Background()
	userID := int64(1)

	mockTasks := []sqlc.Task{
		{ID: 10, Title: "Task 1", Status: "pending", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 11, Title: "Task 2", Status: "completed", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// 1. Expect Count (returns total for pagination)
	mockRepo.On("Count", ctx, mock.AnythingOfType("sqlc.CountTasksByUserParams")).Return(int64(2), nil)

	// 2. Expect List
	mockRepo.On("List", ctx, mock.AnythingOfType("sqlc.ListTasksByUserParams")).Return(mockTasks, nil)

	// Execute
	resp, err := svc.List(ctx, userID, "", 0, false, 10, 0, "", "")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(2), resp.Total) // Check new pagination total
	assert.Len(t, resp.Items, 2)          // Check the array inside the new Response DTO
	assert.Equal(t, "Task 1", resp.Items[0].Title)

	mockRepo.AssertExpectations(t)
}

func TestTaskService_Update(t *testing.T) {
	mockRepo, svc := setupTaskTestHelper(t)
	ctx := context.Background()
	userID := int64(1)
	taskID := int64(10)
	groupID := int64(5)

	req := dto.TaskUpdateRequest{
		Title:   "Updated Task",
		Status:  "completed",
		GroupID: &groupID,
	}

	existingTask := sqlc.Task{
		ID:       taskID,
		UserID:   userID,
		Title:    "Old Task",
		Status:   "pending",
		Priority: "medium",
	}

	updatedTask := sqlc.Task{
		ID:          taskID,
		UserID:      userID,
		GroupID:     sql.NullInt64{Int64: groupID, Valid: true},
		Title:       req.Title,
		Status:      sqlc.TasksStatus(req.Status),
		Priority:    "medium",
		CompletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 1. Expect Initial Get
	mockRepo.On("Get", ctx, sqlc.GetTaskByIDParams{ID: taskID, UserID: userID}).Return(existingTask, nil).Once()

	// 2. Expect Group Access Check (Because GroupID is provided in the update)
	mockRepo.On("GetAccessibleGroup", ctx, sqlc.GetAccessibleGroupByIDParams{ID: groupID, UserID: userID}).
		Return(sqlc.AccessibleGroup{ID: groupID}, nil)

	// 3. Expect Update
	mockRepo.On("Update", ctx, mock.AnythingOfType("sqlc.UpdateTaskParams")).Return(nil)

	// 4. Expect Final Get
	mockRepo.On("Get", ctx, sqlc.GetTaskByIDParams{ID: taskID, UserID: userID}).Return(updatedTask, nil).Once()

	// Execute
	resp, err := svc.Update(ctx, userID, taskID, req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Updated Task", resp.Title)
	assert.Equal(t, "completed", resp.Status)
	assert.NotNil(t, resp.CompletedAt)

	mockRepo.AssertExpectations(t)
}

func TestTaskService_Delete(t *testing.T) {
	mockRepo, svc := setupTaskTestHelper(t)
	ctx := context.Background()
	userID := int64(1)
	taskID := int64(10)

	mockRepo.On("Delete", ctx, sqlc.DeleteTaskParams{ID: taskID, UserID: userID}).Return(nil)

	err := svc.Delete(ctx, userID, taskID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
