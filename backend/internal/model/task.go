package model

type TaskDTO struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	GroupID     *int64  `json:"group_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Priority    string  `json:"priority"`
	DueAt       *string `json:"due_at"`
	CompletedAt *string `json:"completed_at"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
