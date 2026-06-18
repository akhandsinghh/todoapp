# Backend Documentation

## Purpose

The backend is a Go HTTP API built with the Gin framework. It handles authentication, validates requests, applies todo business rules, reads and writes MySQL data through the sqlc query layer, and returns JSON responses to the React frontend.

There is no Node backend in this project.

## How The Backend Starts

The backend starts from:

```text
backend/cmd/server/main.go
```

Startup flow:

```text
main.go
  -> loads backend/.env
  -> opens MySQL connection
  -> runs SQL migrations
  -> creates sqlc Queries
  -> creates repositories
  -> creates services
  -> creates controllers
  -> creates a Gin router
  -> registers routes with CORS, logging, and recovery middleware
  -> starts HTTP server on APP_PORT
```

Default backend URL:

```text
http://localhost:8080
```

## Request Flow

For a protected task request like `POST /api/tasks`, control moves like this:

```text
Browser / React frontend
  -> HTTP request with Authorization header
  -> middleware.CORS
  -> middleware.Logging
  -> middleware.Recovery
  -> middleware.Auth
  -> routes/routes.go
  -> controller/task_controller.go
  -> service/task_service.go
  -> repository/task_repository.go
  -> db/sqlc/task.sql.go
  -> MySQL
```

Response flow goes back in reverse:

```text
MySQL
  -> sqlc row/model
  -> repository
  -> service DTO
  -> controller JSON response
  -> React frontend
```

## Backend File Purpose

### Root Backend Files

`backend/go.mod`

Defines the Go module name and backend dependency list. The backend uses Gin and the MySQL driver.

`backend/go.sum`

Stores dependency checksums after `go mod tidy` downloads modules.

`backend/.env`

Local backend environment configuration. It contains server port, JWT secret, MySQL credentials, database name, and allowed frontend origin.

`backend/.env.example`

Example environment file for other developers.

`backend/sqlc.yaml`

sqlc configuration. It tells sqlc where migrations and query files live, and where generated Go code should be placed.

### `cmd/server`

`backend/cmd/server/main.go`

Application entrypoint. It builds the whole backend dependency graph, runs migrations, starts the reminder scheduler, registers routes, and starts the HTTP server.

### `internal/controller`

Controllers are Gin-facing files. They decode JSON, read path/query parameters, call services, and write JSON responses.

`auth_controller.go`

Handles register, login, and current-user endpoints.

`task_controller.go`

Handles listing, creating, updating, completing, and deleting tasks.

`group_controller.go`

Handles task group CRUD endpoints.

`reminder_controller.go`

Handles reminder creation, listing, and deletion endpoints.

### `internal/service`

Services contain business logic. Controllers should stay thin, while services decide validation, defaults, status changes, ownership checks, and DTO conversion.

`auth_service.go`

Registers users, logs users in, checks passwords, creates JWT tokens, and loads the current user.

`task_service.go`

Validates task input, parses due dates, applies default priority/status logic, toggles completion timestamps, and converts database tasks to API DTOs.

`group_service.go`

Validates group names, applies default colors, and converts database groups to API DTOs.

`reminder_service.go`

Validates reminder input, ensures the task belongs to the current user, creates reminders, lists reminders, and supports scheduler operations.

### `internal/repository`

Repositories isolate database access. They call the sqlc-generated query methods and hide sqlc details from services.

`user_repository.go`

Creates users and loads users by email or ID.

`task_repository.go`

Creates, lists, loads, updates, and deletes tasks.

`group_repository.go`

Creates, lists, loads, updates, and deletes task groups.

`reminder_repository.go`

Creates, lists, finds due reminders, marks reminders sent, and deletes reminders.

### `internal/middleware`

Middleware runs before or around route handlers.

`auth.go`

Checks `Authorization: Bearer <token>`, verifies the JWT, and places `userID` into the Gin context.

`cors.go`

Allows the React frontend origin to call the Go API during development.

`logging.go`

Logs HTTP method, path, and request duration.

`recovery.go`

Recovers from panics and returns a JSON 500 response instead of crashing the process.

### `internal/model`

Models define request and response shapes used by controllers and services.

`request.go`

Request payload structs for auth, groups, tasks, and reminders.

`response.go`

Shared API response structs such as auth response and message response.

`user.go`

Safe user response shape. It never exposes password hashes.

`task.go`

Task response DTO used by the API.

`group.go`

Group response DTO used by the API.

`reminder.go`

Reminder response DTO used by the API.

### `internal/routes`

`routes.go`

Creates the Gin router, registers all HTTP paths, and maps each method/path to the correct controller function. It also applies auth middleware to protected routes.

### `internal/scheduler`

`reminder_scheduler.go`

Runs a background ticker every minute. It finds due reminders, logs them, and marks them as sent.

### `internal/db`

`connection.go`

Loads `.env`, builds the MySQL DSN, opens the database connection, configures connection pooling, pings MySQL, and runs migration files.

See `docs/DATABASE.md` for migrations, query files, and sqlc details.

### `internal/util`

Reusable helpers.

`jwt.go`

Creates and verifies HMAC-signed JWT tokens.

`password.go`

Hashes and checks passwords.

`validator.go`

Small validation helpers for required fields, task priority, and task status.

`response.go`

Writes JSON success and error responses.

`pagination.go`

Parses `page` and `limit` query parameters into limit/offset values.

## Backend Routes

Public:

```text
POST /api/auth/register
POST /api/auth/login
GET  /health
```

Protected:

```text
GET    /api/auth/me
GET    /api/groups
POST   /api/groups
PUT    /api/groups/{id}
DELETE /api/groups/{id}
GET    /api/tasks
POST   /api/tasks
PUT    /api/tasks/{id}
DELETE /api/tasks/{id}
GET    /api/reminders
POST   /api/reminders
DELETE /api/reminders/{id}
```

## Running Backend

```bash
cd backend
go mod tidy
go run ./cmd/server
```

Before running, create the MySQL database:

```sql
CREATE DATABASE todo_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```
