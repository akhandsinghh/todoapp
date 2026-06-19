package model

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type GroupRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
type ShareGroupRequest struct {
	Email string `json:"email"`
}
type TaskRequest struct {
	GroupID     *int64 `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueAt       string `json:"due_at"`
}
type ReminderRequest struct {
	TaskID   int64  `json:"task_id"`
	RemindAt string `json:"remind_at"`
	Message  string `json:"message"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ForgotPasswordRequest struct {
	Email           string `json:"email"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
