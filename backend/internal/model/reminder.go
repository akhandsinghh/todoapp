package model

type ReminderDTO struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	TaskID    int64  `json:"task_id"`
	RemindAt  string `json:"remind_at"`
	Message   string `json:"message"`
	Sent      bool   `json:"sent"`
	CreatedAt string `json:"created_at"`
}
