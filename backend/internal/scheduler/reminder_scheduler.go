package scheduler

import (
    "context"
    "log"
    "time"
    "todo-app/backend/internal/service"
)

type ReminderScheduler struct { reminders *service.ReminderService }
func NewReminderScheduler(reminders *service.ReminderService) *ReminderScheduler { return &ReminderScheduler{reminders:reminders} }
func (s *ReminderScheduler) Start(ctx context.Context) {
    ticker := time.NewTicker(time.Minute)
    go func() {
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done(): return
            case <-ticker.C: s.process(ctx)
            }
        }
    }()
}
func (s *ReminderScheduler) process(ctx context.Context) {
    items, err := s.reminders.Due(ctx, 100); if err != nil { log.Printf("reminder scheduler error: %v", err); return }
    for _, r := range items { log.Printf("reminder due: user=%d task=%d message=%s", r.UserID, r.TaskID, r.Message.String); _ = s.reminders.MarkSent(ctx, r.ID) }
}
