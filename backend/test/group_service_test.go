package test

import (
	"context"
	"testing"
	"time"

	"todo-app/backend/internal/controller/dto"
	"todo-app/backend/internal/db/sqlc"
	"todo-app/backend/internal/service"
	"todo-app/backend/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupGroupTestHelper initializes the mock and service for group tests.
func setupGroupTestHelper(t *testing.T) (*mocks.GroupRepository, *service.GroupService) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo)
	return mockRepo, svc
}

func TestGroupService_Create(t *testing.T) {
	mockRepo, svc := setupGroupTestHelper(t)
	ctx := context.Background()
	userID := int64(1)

	req := dto.GroupRequest{
		Name:  "Family Trip",
		Color: "#ff0000",
	}

	groupID := int64(10)
	mockReturnedGroup := sqlc.TaskGroup{
		ID:        groupID,
		UserID:    userID,
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 1. Expect Create to be called
	mockRepo.On("Create", ctx, mock.AnythingOfType("sqlc.CreateGroupParams")).Return(groupID, nil)

	// 2. Expect Get to be called immediately after to fetch the DTO data
	mockRepo.On("Get", ctx, sqlc.GetGroupByIDParams{ID: groupID, UserID: userID}).Return(mockReturnedGroup, nil)

	// Execute
	resp, err := svc.Create(ctx, userID, req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, groupID, resp.ID)
	assert.Equal(t, "Family Trip", resp.Name)
	assert.Equal(t, "#ff0000", resp.Color)
	assert.Equal(t, "creator", resp.Role)

	mockRepo.AssertExpectations(t)
}

func TestGroupService_List(t *testing.T) {
	mockRepo, svc := setupGroupTestHelper(t)
	ctx := context.Background()
	userID := int64(1)

	mockGroups := []sqlc.AccessibleGroup{
		{ID: 10, UserID: userID, Name: "My Group", Color: "#fff", Role: "creator", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 11, UserID: int64(2), Name: "Shared Group", Color: "#000", Role: "member", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// 1. Expect List to return our mock array
	mockRepo.On("List", ctx, userID).Return(mockGroups, nil)

	// Execute
	resp, err := svc.List(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "My Group", resp[0].Name)
	assert.Equal(t, "creator", resp[0].Role)
	assert.Equal(t, "Shared Group", resp[1].Name)
	assert.Equal(t, "member", resp[1].Role)

	mockRepo.AssertExpectations(t)
}

func TestGroupService_Share(t *testing.T) {
	mockRepo, svc := setupGroupTestHelper(t)
	ctx := context.Background()
	creatorID := int64(1)
	groupID := int64(10)
	targetUserEmail := "friend@example.com"
	targetUserID := int64(2)

	req := dto.ShareGroupRequest{
		Email: targetUserEmail,
	}

	mockGroup := sqlc.TaskGroup{
		ID:     groupID,
		UserID: creatorID,
	}

	mockTargetUser := sqlc.User{
		ID:    targetUserID,
		Email: targetUserEmail,
	}

	// 1. Expect Get (verifies the group exists and the user is the creator)
	mockRepo.On("Get", ctx, sqlc.GetGroupByIDParams{ID: groupID, UserID: creatorID}).Return(mockGroup, nil)

	// 2. Expect UserByEmail (finds the user we are sharing with)
	mockRepo.On("UserByEmail", ctx, targetUserEmail).Return(mockTargetUser, nil)

	// 3. Expect Share (actually inserts the record)
	mockRepo.On("Share", ctx, sqlc.CreateGroupShareParams{GroupID: groupID, UserID: targetUserID}).Return(nil)

	// Execute
	err := svc.Share(ctx, creatorID, groupID, req)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGroupService_Update(t *testing.T) {
	mockRepo, svc := setupGroupTestHelper(t)
	ctx := context.Background()
	userID := int64(1)
	groupID := int64(10)

	req := dto.GroupUpdateRequest{
		Name:  "Updated Name",
		Color: "#00ff00",
	}

	// 1. Expect Update to be called
	mockRepo.On("Update", ctx, sqlc.UpdateGroupParams{
		ID:     groupID,
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}).Return(nil)

	// Execute
	err := svc.Update(ctx, userID, groupID, req)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGroupService_Delete(t *testing.T) {
	mockRepo, svc := setupGroupTestHelper(t)
	ctx := context.Background()
	userID := int64(1)
	groupID := int64(10)

	// 1. Expect Delete to be called
	mockRepo.On("Delete", ctx, sqlc.DeleteGroupParams{ID: groupID, UserID: userID}).Return(nil)

	// Execute
	err := svc.Delete(ctx, userID, groupID)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
