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

// setupTaskTestHelper initializes the mock and service for task tests.
func setupTaskTestHelper(t *testing.T) (*mocks.TaskRepository, *service.TaskService) {
	mockRepo := mocks.NewTaskRepository(t)
	svc := service.NewTaskService(mockRepo)
	return mockRepo, svc
}

func TestTaskService_Create(t *testing.T) {
	mockRepo, svc := setupTaskTestHelper(t)
	ctx := context.Background()
	userID := int64(1)

	// We use a nil group ID to bypass the ensureGroupAccess check for this test
	req := dto.TaskRequest{
		Title:       "Buy Groceries",
		Description: "Milk, Eggs, Bread",
		Priority:    "high",
		DueAt:       time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	taskID := int64(10)
	mockReturnedTask := sqlc.Task{
		ID:          taskID,
		UserID:      userID,
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: true},
		Status:      "pending",
		Priority:    sqlc.TasksPriority(req.Priority),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 1. Expect Create to be called with any valid sqlc.CreateTaskParams
	mockRepo.On("Create", ctx, mock.AnythingOfType("sqlc.CreateTaskParams")).Return(taskID, nil)

	// 2. Expect Get to be called immediately after to return the full DTO
	mockRepo.On("Get", ctx, mock.AnythingOfType("sqlc.GetTaskByIDParams")).Return(mockReturnedTask, nil)

	// Execute
	resp, err := svc.Create(ctx, userID, req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, taskID, resp.ID)
	assert.Equal(t, "Buy Groceries", resp.Title)
	assert.Equal(t, "pending", resp.Status)
	assert.Equal(t, "high", resp.Priority)

	mockRepo.AssertExpectations(t)
}

func TestTaskService_List(t *testing.T) {
	mockRepo, svc := setupTaskTestHelper(t)
	ctx := context.Background()
	userID := int64(1)

	mockTasks := []sqlc.Task{
		{ID: 1, UserID: userID, Title: "Task 1", Status: "pending", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, UserID: userID, Title: "Task 2", Status: "completed", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// 1. Expect Count to return 2 total tasks
	mockRepo.On("Count", ctx, mock.AnythingOfType("sqlc.CountTasksByUserParams")).Return(int64(2), nil)

	// 2. Expect List to return our mock array
	mockRepo.On("List", ctx, mock.AnythingOfType("sqlc.ListTasksByUserParams")).Return(mockTasks, nil)

	// Execute (status="", groupID=0, ungrouped=false, limit=10, offset=0, sortBy="due_at", sortOrder="desc")
	resp, err := svc.List(ctx, userID, "", 0, false, 10, 0, "due_at", "desc")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Items, 2)
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
		Title:    "Updated Title",
		Status:   "completed",
		Priority: "low",
		GroupID:  &groupID,
	}

	currentTask := sqlc.Task{
		ID:        taskID,
		UserID:    userID,
		Title:     "Old Title",
		Status:    "pending",
		Priority:  "medium",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedTask := sqlc.Task{
		ID:          taskID,
		UserID:      userID,
		GroupID:     sql.NullInt64{Int64: groupID, Valid: true},
		Title:       req.Title,
		Status:      sqlc.TasksStatus(req.Status),
		Priority:    sqlc.TasksPriority(req.Priority),
		CompletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		CreatedAt:   currentTask.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	// 1. Expect Get (initial fetch to check if task exists)
	mockRepo.On("Get", ctx, sqlc.GetTaskByIDParams{ID: taskID, UserID: userID}).Return(currentTask, nil).Once()

	// 2. Expect ensureGroupAccess (GetAccessibleGroup) since groupID is provided
	mockRepo.On("GetAccessibleGroup", ctx, mock.AnythingOfType("sqlc.GetAccessibleGroupByIDParams")).Return(sqlc.AccessibleGroup{}, nil)
	// 3. Expect Update
	mockRepo.On("Update", ctx, mock.AnythingOfType("sqlc.UpdateTaskParams")).Return(nil)

	// 4. Expect Get (final fetch to return updated DTO)
	mockRepo.On("Get", ctx, sqlc.GetTaskByIDParams{ID: taskID, UserID: userID}).Return(updatedTask, nil).Once()

	// Execute
	resp, err := svc.Update(ctx, userID, taskID, req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", resp.Title)
	assert.Equal(t, "completed", resp.Status)
	assert.NotNil(t, resp.CompletedAt) // completed_at should now be populated

	mockRepo.AssertExpectations(t)
}

func TestTaskService_Delete(t *testing.T) {
	mockRepo, svc := setupTaskTestHelper(t)
	ctx := context.Background()
	userID := int64(1)
	taskID := int64(10)

	// 1. Expect Delete to be called and succeed
	mockRepo.On("Delete", ctx, sqlc.DeleteTaskParams{ID: taskID, UserID: userID}).Return(nil)

	// Execute
	err := svc.Delete(ctx, userID, taskID)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
