package service

import (
	"context"
	"database/sql"
	"errors"
	apperr "todo-app/backend/internal/errors"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/repository"
	"todo-app/backend/internal/util"
)

type AuthService struct {
	users  *repository.UserRepository
	secret string
}

func NewAuthService(users *repository.UserRepository, secret string) *AuthService {
	return &AuthService{users: users, secret: secret}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (model.AuthResponse, error) {
	if !util.Required(req.Name, req.Email, req.Password) {
		return model.AuthResponse{}, apperr.BadRequest("name, email and password are required")
	}
	if _, err := s.users.ByEmail(ctx, req.Email); err == nil {
		return model.AuthResponse{}, apperr.Conflict("email already registered")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return model.AuthResponse{}, err
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return model.AuthResponse{}, err
	}
	id, err := s.users.Create(ctx, req.Name, req.Email, hash)
	if err != nil {
		return model.AuthResponse{}, apperr.Internal("failed to create user")
	}
	token, err := util.SignToken(id, req.Email, s.secret)
	if err != nil {
		return model.AuthResponse{}, err
	}
	return model.AuthResponse{Token: token, User: model.UserResponse{ID: id, Name: req.Name, Email: req.Email}}, nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error) {
	if !util.Required(req.Email, req.Password) {
		return model.AuthResponse{}, apperr.BadRequest("email and password are required")
	}
	u, err := s.users.ByEmail(ctx, req.Email)
	if err != nil {
		return model.AuthResponse{}, apperr.Unauthorized("invalid credentials")
	}
	if !util.CheckPassword(req.Password, u.PasswordHash) {
		return model.AuthResponse{}, apperr.Unauthorized("invalid credentials")
	}
	token, err := util.SignToken(u.ID, u.Email, s.secret)
	if err != nil {
		return model.AuthResponse{}, err
	}
	return model.AuthResponse{Token: token, User: model.UserResponse{ID: u.ID, Name: u.Name, Email: u.Email}}, nil
}

func (s *AuthService) Me(ctx context.Context, userID int64) (model.UserResponse, error) {
	u, err := s.users.ByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserResponse{}, apperr.NotFound("user not found")
		}
		return model.UserResponse{}, apperr.Internal("failed to fetch user")
	}
	return model.UserResponse{ID: u.ID, Name: u.Name, Email: u.Email}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID int64, req model.ChangePasswordRequest) error {
	if !util.Required(req.OldPassword, req.NewPassword, req.ConfirmPassword) {
		return apperr.BadRequest("old password, new password and confirm password are required")
	}
	if req.NewPassword != req.ConfirmPassword {
		return apperr.BadRequest("new passwords do not match")
	}
	u, err := s.users.ByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperr.NotFound("user not found")
		}
		return apperr.Internal("failed to fetch user")
	}
	if !util.CheckPassword(req.OldPassword, u.PasswordHash) {
		return apperr.BadRequest("old password is incorrect")
	}
	hash, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return apperr.Internal("failed to hash password")
	}
	err = s.users.UpdatePassword(ctx, userID, hash)
	if err != nil {
		return apperr.Internal("failed to update password")
	}
	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) error {
	if !util.Required(req.Email, req.NewPassword, req.ConfirmPassword) {
		return apperr.BadRequest("email, new password and confirm password are required")
	}
	if req.NewPassword != req.ConfirmPassword {
		return apperr.BadRequest("passwords do not match")
	}
	u, err := s.users.ByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperr.NotFound("user not found")
		}
		return apperr.Internal("failed to fetch user")
	}
	hash, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return apperr.Internal("failed to hash password")
	}
	err = s.users.UpdatePassword(ctx, u.ID, hash)
	if err != nil {
		return apperr.Internal("failed to update password")
	}
	return nil
}
