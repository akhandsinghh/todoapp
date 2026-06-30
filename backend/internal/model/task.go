package model

type TaskDTO struct {
	ID          int64
	UserID      int64
	GroupID     *int64
	Title       string
	Description string
	Status      string
	Priority    string
	DueAt       *string
	CompletedAt *string
	CreatedAt   string
	UpdatedAt   string
}
