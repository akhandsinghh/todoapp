# Database Documentation

## Purpose

The database layer stores users, task groups, tasks, and reminders in MySQL. The backend accesses MySQL through sqlc query methods, wrapped by repository files.

## Database Startup Flow

When the Go backend starts:

```text
cmd/server/main.go
  -> db.LoadEnv(".env")
  -> db.Connect(db.FromEnv())
  -> db.RunMigrations(database, "internal/db/migrations")
  -> sqlc.New(database)
  -> repositories receive sqlc Queries
```

The database must exist before backend startup:

```sql
CREATE DATABASE todo_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

The backend creates tables automatically by running the migration SQL files.

## Database File Purpose

### `backend/internal/db/connection.go`

Handles database configuration and connection setup.

Main responsibilities:

- Reads environment variables.
- Builds the MySQL DSN.
- Opens the database connection.
- Configures connection pooling.
- Pings MySQL.
- Runs migration SQL files.

### `backend/internal/db/migrations`

Migration files define the actual MySQL tables. They are executed by `RunMigrations` when the backend starts.

`001_users.sql`

Creates the `users` table.

Purpose:

- Stores registered users.
- Keeps name, email, password hash, created timestamp, and updated timestamp.
- Enforces unique email addresses.

Main columns:

```text
id
name
email
password_hash
created_at
updated_at
```

`002_task_groups.sql`

Creates the `task_groups` table.

Purpose:

- Stores user-owned task groups.
- Lets each user organize tasks by category.
- Deletes groups when the owning user is deleted.

Main columns:

```text
id
user_id
name
color
created_at
updated_at
```

`003_tasks.sql`

Creates the `tasks` table.

Purpose:

- Stores todos for each user.
- Optionally links a task to a group.
- Tracks status, priority, due date, and completion time.

Main columns:

```text
id
user_id
group_id
title
description
status
priority
due_at
completed_at
created_at
updated_at
```

`004_reminders.sql`

Creates the `reminders` table.

Purpose:

- Stores reminders linked to tasks.
- The scheduler finds unsent reminders when `remind_at` is due.

Main columns:

```text
id
user_id
task_id
remind_at
message
sent
created_at
updated_at
```

## Table Relationships

```text
users
  -> task_groups
  -> tasks
  -> reminders
```

Detailed relationships:

```text
users.id -> task_groups.user_id
users.id -> tasks.user_id
task_groups.id -> tasks.group_id
users.id -> reminders.user_id
tasks.id -> reminders.task_id
```

Delete behavior:

```text
Deleting a user deletes that user's groups, tasks, and reminders.
Deleting a task deletes its reminders.
Deleting a group sets related task group_id to NULL.
```

## Query Files

Query files are the sqlc source files. They define named SQL queries that sqlc can turn into Go methods.

`backend/internal/db/queries/user.sql`

Contains:

```text
CreateUser
GetUserByEmail
GetUserByID
```

Used by:

```text
repository/user_repository.go
service/auth_service.go
```

`backend/internal/db/queries/group.sql`

Contains:

```text
CreateGroup
ListGroupsByUser
GetGroupByID
UpdateGroup
DeleteGroup
```

Used by:

```text
repository/group_repository.go
service/group_service.go
```

`backend/internal/db/queries/task.sql`

Contains:

```text
CreateTask
ListTasksByUser
GetTaskByID
UpdateTask
DeleteTask
```

Used by:

```text
repository/task_repository.go
service/task_service.go
```

`backend/internal/db/queries/reminder.sql`

Contains:

```text
CreateReminder
ListRemindersByUser
ListDueReminders
MarkReminderSent
DeleteReminder
```

Used by:

```text
repository/reminder_repository.go
service/reminder_service.go
scheduler/reminder_scheduler.go
```

## sqlc Generated Layer

The generated-style files live in:

```text
backend/internal/db/sqlc
```

`db.go`

Defines the `DBTX` interface and `Queries` struct. This lets sqlc methods work with `*sql.DB` or a transaction.

`models.go`

Defines Go structs that represent database rows:

```text
User
TaskGroup
Task
Reminder
```

`querier.go`

Defines the `Querier` interface plus parameter structs used by generated query methods.

`user.sql.go`

Implements Go methods for user SQL queries.

`group.sql.go`

Implements Go methods for group SQL queries.

`task.sql.go`

Implements Go methods for task SQL queries.

`reminder.sql.go`

Implements Go methods for reminder SQL queries.

## Database Access Flow

Example: loading tasks.

```text
Dashboard.jsx
  -> GET /api/tasks
  -> task_controller.go
  -> task_service.go
  -> task_repository.go
  -> sqlc.ListTasksByUser()
  -> MySQL tasks table
```

Example: registering a user.

```text
Register.jsx
  -> POST /api/auth/register
  -> auth_controller.go
  -> auth_service.go
  -> user_repository.go
  -> sqlc.CreateUser()
  -> MySQL users table
```

## Regenerating sqlc Code

If sqlc is installed:

```bash
cd backend
sqlc generate
```

sqlc reads:

```text
backend/sqlc.yaml
backend/internal/db/migrations
backend/internal/db/queries
```

and writes generated Go files to:

```text
backend/internal/db/sqlc
```

## Important Notes

- MySQL is the only database used.
- The backend owns database access.
- The frontend never connects directly to MySQL.
- Repositories are the only backend layer that should call sqlc methods directly.
- Services should call repositories, not raw SQL.

