package model

type ReminderDTO struct {
	ID        int64
	UserID    int64
	TaskID    int64
	RemindAt  string
	Message   string
	Sent      bool
	CreatedAt string
}
