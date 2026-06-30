package dto

type RegisterRequest struct {
	Name     string `json:"name" binding:"required" validate:"required"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required"`
}

type GroupRequest struct {
	Name  string `json:"name" binding:"required" validate:"required"`
	Color string `json:"color" binding:"omitempty" validate:"omitempty"`
}

type GroupUpdateRequest struct {
	Name  string `json:"name" binding:"omitempty" validate:"omitempty"`
	Color string `json:"color" binding:"omitempty" validate:"omitempty"`
}

type ShareGroupRequest struct {
	Email string `json:"email" binding:"required,email" validate:"required,email"`
}

type TaskRequest struct {
	GroupID     *int64 `json:"group_id" binding:"omitempty" validate:"omitempty"`
	Title       string `json:"title" binding:"required" validate:"required"`
	Description string `json:"description" binding:"omitempty" validate:"omitempty"`
	Status      string `json:"status" binding:"omitempty,oneof=pending completed" validate:"omitempty,oneof=pending completed"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low medium high" validate:"omitempty,oneof=low medium high"`
	DueAt       string `json:"due_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type TaskUpdateRequest struct {
	GroupID     *int64 `json:"group_id" binding:"omitempty" validate:"omitempty"`
	Title       string `json:"title" binding:"omitempty" validate:"omitempty"`
	Description string `json:"description" binding:"omitempty" validate:"omitempty"`
	Status      string `json:"status" binding:"omitempty,oneof=pending completed" validate:"omitempty,oneof=pending completed"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low medium high" validate:"omitempty,oneof=low medium high"`
	DueAt       string `json:"due_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type ReminderRequest struct {
	TaskID   int64  `json:"task_id" binding:"required" validate:"required"`
	RemindAt string `json:"remind_at" binding:"required,datetime=2006-01-02T15:04:05Z07:00" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Message  string `json:"message" binding:"omitempty" validate:"omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword" validate:"required,eqfield=NewPassword"`
}

type ForgotPasswordRequest struct {
	Email           string `json:"email" binding:"required,email" validate:"required,email"`
	NewPassword     string `json:"new_password" binding:"required,min=6" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword" validate:"required,eqfield=NewPassword"`
}
