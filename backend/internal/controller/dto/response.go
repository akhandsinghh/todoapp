package dto

import "todo-app/backend/internal/model"

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TaskResponse struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	GroupID     *int64  `json:"group_id,omitempty"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Status      string  `json:"status"`
	Priority    string  `json:"priority"`
	DueAt       *string `json:"due_at,omitempty"`
	CompletedAt *string `json:"completed_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type TaskListResponse struct {
	Items []TaskResponse `json:"items"`
	Total int64          `json:"total"`
}

type GroupResponse struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ReminderResponse struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	TaskID    int64  `json:"task_id"`
	RemindAt  string `json:"remind_at"`
	Message   string `json:"message"`
	Sent      bool   `json:"sent"`
	CreatedAt string `json:"created_at"`
}

// Conversion helpers: domain (model) -> DTO
func ConvertUserDomainToResponse(u model.UserResponse) UserResponse {
	return UserResponse{ID: u.ID, Name: u.Name, Email: u.Email}
}

func ConvertTaskDomainToResponse(t model.TaskDTO) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		GroupID:     t.GroupID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Priority:    t.Priority,
		DueAt:       t.DueAt,
		CompletedAt: t.CompletedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func ConvertTaskListDomainToResponse(items []model.TaskDTO, total int64) TaskListResponse {
	out := make([]TaskResponse, 0, len(items))
	for _, it := range items {
		out = append(out, ConvertTaskDomainToResponse(it))
	}
	return TaskListResponse{Items: out, Total: total}
}

func ConvertGroupDomainToResponse(g model.GroupDTO) GroupResponse {
	return GroupResponse{
		ID:        g.ID,
		UserID:    g.UserID,
		Name:      g.Name,
		Color:     g.Color,
		Role:      g.Role,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}

func ConvertGroupListDomainToResponse(items []model.GroupDTO) []GroupResponse {
	out := make([]GroupResponse, 0, len(items))
	for _, it := range items {
		out = append(out, ConvertGroupDomainToResponse(it))
	}
	return out
}

func ConvertReminderDomainToResponse(r model.ReminderDTO) ReminderResponse {
	return ReminderResponse{
		ID:        r.ID,
		UserID:    r.UserID,
		TaskID:    r.TaskID,
		RemindAt:  r.RemindAt,
		Message:   r.Message,
		Sent:      r.Sent,
		CreatedAt: r.CreatedAt,
	}
}

func ConvertReminderListDomainToResponse(items []model.ReminderDTO) []ReminderResponse {
	out := make([]ReminderResponse, 0, len(items))
	for _, it := range items {
		out = append(out, ConvertReminderDomainToResponse(it))
	}
	return out
}
