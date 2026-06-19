package service

import (
	"context"
	"database/sql"
	"errors"
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
		return model.AuthResponse{}, errors.New("name, email and password are required")
	}
	if _, err := s.users.ByEmail(ctx, req.Email); err == nil {
		return model.AuthResponse{}, errors.New("email already registered")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return model.AuthResponse{}, err
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return model.AuthResponse{}, err
	}
	id, err := s.users.Create(ctx, req.Name, req.Email, hash)
	if err != nil {
		return model.AuthResponse{}, err
	}
	token, err := util.SignToken(id, req.Email, s.secret)
	if err != nil {
		return model.AuthResponse{}, err
	}
	return model.AuthResponse{Token: token, User: model.UserResponse{ID: id, Name: req.Name, Email: req.Email}}, nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error) {
	if !util.Required(req.Email, req.Password) {
		return model.AuthResponse{}, errors.New("email and password are required")
	}
	u, err := s.users.ByEmail(ctx, req.Email)
	if err != nil {
		return model.AuthResponse{}, errors.New("invalid credentials")
	}
	if !util.CheckPassword(req.Password, u.PasswordHash) {
		return model.AuthResponse{}, errors.New("invalid credentials")
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
		return model.UserResponse{}, err
	}
	return model.UserResponse{ID: u.ID, Name: u.Name, Email: u.Email}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID int64, req model.ChangePasswordRequest) error {
	if !util.Required(req.OldPassword, req.NewPassword, req.ConfirmPassword) {
		return errors.New("old password, new password and confirm password are required")
	}
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("new passwords do not match")
	}
	u, err := s.users.ByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if !util.CheckPassword(req.OldPassword, u.PasswordHash) {
		return errors.New("old password is incorrect")
	}
	hash, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	return s.users.UpdatePassword(ctx, userID, hash)
}

func (s *AuthService) ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) error {
	if !util.Required(req.Email, req.NewPassword, req.ConfirmPassword) {
		return errors.New("email, new password and confirm password are required")
	}
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	u, err := s.users.ByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("user not found")
	}
	hash, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	return s.users.UpdatePassword(ctx, u.ID, hash)
}
