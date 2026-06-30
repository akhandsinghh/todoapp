package service

import (
	"context"
	"database/sql"
	"errors"

	"todo-app/backend/internal/controller/dto"
	"todo-app/backend/internal/db/sqlc"
	apperr "todo-app/backend/internal/errors"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/util"
)

type UserRepository interface {
	ByEmail(ctx context.Context, email string) (sqlc.User, error)
	ByID(ctx context.Context, id int64) (sqlc.User, error)
	Create(ctx context.Context, name, email, hash string) (int64, error)
	UpdatePassword(ctx context.Context, id int64, hash string) error
}

type AuthService struct {
	users  UserRepository
	secret string
}

func NewAuthService(users UserRepository, secret string) *AuthService {
	return &AuthService{users: users, secret: secret}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error) {
	if _, err := s.users.ByEmail(ctx, req.Email); err == nil {
		return dto.AuthResponse{}, apperr.Conflict("email already registered")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return dto.AuthResponse{}, err
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	id, err := s.users.Create(ctx, req.Name, req.Email, hash)
	if err != nil {
		return dto.AuthResponse{}, apperr.Internal("failed to create user")
	}
	token, err := util.SignToken(id, req.Email, s.secret)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	return dto.AuthResponse{Token: token, User: dto.ConvertUserDomainToResponse(model.UserResponse{ID: id, Name: req.Name, Email: req.Email})}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	u, err := s.users.ByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, apperr.Unauthorized("invalid credentials")
	}
	if !util.CheckPassword(req.Password, u.PasswordHash) {
		return dto.AuthResponse{}, apperr.Unauthorized("invalid credentials")
	}
	token, err := util.SignToken(u.ID, u.Email, s.secret)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	return dto.AuthResponse{Token: token, User: dto.ConvertUserDomainToResponse(model.UserResponse{ID: u.ID, Name: u.Name, Email: u.Email})}, nil
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

func (s *AuthService) ChangePassword(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error {
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

func (s *AuthService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
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
