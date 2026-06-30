package test

import (
	"context"
	"database/sql"
	"testing"

	"todo-app/backend/internal/controller/dto"
	"todo-app/backend/internal/db/sqlc"
	"todo-app/backend/internal/service"
	"todo-app/backend/internal/util"
	"todo-app/backend/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupTestHelper is a quick utility to set up the mock and service for each test.
func setupTestHelper(t *testing.T) (*mocks.UserRepository, *service.AuthService, string) {
	mockRepo := mocks.NewUserRepository(t)
	secret := "test-secret-key"
	svc := service.NewAuthService(mockRepo, secret)
	return mockRepo, svc, secret
}

func TestAuthService_Register(t *testing.T) {
	mockRepo, svc, _ := setupTestHelper(t)
	ctx := context.Background()

	req := dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// 1. Expect checking if user exists to return sql.ErrNoRows (user does not exist)
	mockRepo.On("ByEmail", ctx, req.Email).Return(sqlc.User{}, sql.ErrNoRows)

	// 2. Expect Create to be called. We use mock.Anything for the hash because bcrypt generates a unique hash every time.
	mockRepo.On("Create", ctx, req.Name, req.Email, mock.AnythingOfType("string")).Return(int64(1), nil)

	resp, err := svc.Register(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)

	// Access the User field directly as dto.UserResponse
	userResp := resp.User
	assert.Equal(t, int64(1), userResp.ID)
	assert.Equal(t, "John Doe", userResp.Name)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo, svc, _ := setupTestHelper(t)
	ctx := context.Background()

	password := "password123"
	// Generate a real hash so util.CheckPassword passes inside the Login function
	hashedPassword, _ := util.HashPassword(password)

	req := dto.LoginRequest{
		Email:    "john@example.com",
		Password: password,
	}

	mockUser := sqlc.User{
		ID:           1,
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: hashedPassword,
	}

	// 1. Expect ByEmail to find the user
	mockRepo.On("ByEmail", ctx, req.Email).Return(mockUser, nil)

	resp, err := svc.Login(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)

	// Access the User field directly as dto.UserResponse
	userResp := resp.User
	assert.Equal(t, int64(1), userResp.ID)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Me(t *testing.T) {
	mockRepo, svc, _ := setupTestHelper(t)
	ctx := context.Background()
	userID := int64(1)

	mockUser := sqlc.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// 1. Expect ByID to return the user
	mockRepo.On("ByID", ctx, userID).Return(mockUser, nil)

	resp, err := svc.Me(ctx, userID)

	// Assertions (return type model.UserResponse already)
	assert.NoError(t, err)
	assert.Equal(t, userID, resp.ID)
	assert.Equal(t, "John Doe", resp.Name)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_ChangePassword(t *testing.T) {
	mockRepo, svc, _ := setupTestHelper(t)
	ctx := context.Background()
	userID := int64(1)

	oldPassword := "oldpass123"
	newPassword := "newpass456"
	oldHashedPassword, _ := util.HashPassword(oldPassword)

	req := dto.ChangePasswordRequest{
		OldPassword:     oldPassword,
		NewPassword:     newPassword,
		ConfirmPassword: newPassword,
	}

	mockUser := sqlc.User{
		ID:           userID,
		PasswordHash: oldHashedPassword,
	}

	// 1. Expect ByID to find the user
	mockRepo.On("ByID", ctx, userID).Return(mockUser, nil)

	// 2. Expect UpdatePassword to be called with the new hash
	mockRepo.On("UpdatePassword", ctx, userID, mock.AnythingOfType("string")).Return(nil)

	err := svc.ChangePassword(ctx, userID, req)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_ForgotPassword(t *testing.T) {
	mockRepo, svc, _ := setupTestHelper(t)
	ctx := context.Background()

	newPassword := "newpass456"

	req := dto.ForgotPasswordRequest{
		Email:           "john@example.com",
		NewPassword:     newPassword,
		ConfirmPassword: newPassword,
	}

	mockUser := sqlc.User{
		ID:    1,
		Email: "john@example.com",
	}

	// 1. Expect ByEmail to find the user
	mockRepo.On("ByEmail", ctx, req.Email).Return(mockUser, nil)

	// 2. Expect UpdatePassword to be called
	mockRepo.On("UpdatePassword", ctx, mockUser.ID, mock.AnythingOfType("string")).Return(nil)

	err := svc.ForgotPassword(ctx, req)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
