# Todo App Flow

This document explains where control goes across the whole todo app, from the React frontend to the Gin backend, sqlc query layer, and MySQL database.

## Complete Startup Flow

```text
Developer starts MySQL
  -> Developer creates/uses database todo_app
  -> Developer starts backend with go run ./cmd/server
  -> backend/cmd/server/main.go
  -> internal/db/connection.go loads .env and opens MySQL connection
  -> internal/db/connection.go runs internal/db/migrations/*.sql
  -> internal/db/sqlc creates Queries
  -> internal/repository files are created with sqlc Queries
  -> internal/service files are created with repositories
  -> internal/controller files are created with services
  -> internal/scheduler/reminder_scheduler.go starts in the background
  -> internal/routes/routes.go creates the Gin router and registers routes
  -> Gin starts listening on APP_PORT, default http://localhost:8080

Developer starts frontend with npm start
  -> frontend/webpack.config.js serves the React bundle
  -> frontend/public/index.html loads the bundle
  -> frontend/src/index.js mounts React
  -> frontend/src/App.jsx sets up AuthProvider and routes
  -> Browser shows login/register/dashboard depending on auth state
```

## Frontend Application Flow

```text
frontend/src/index.js
  -> renders App
  -> frontend/src/App.jsx
  -> wraps the app in AuthProvider
  -> creates BrowserRouter routes
```

Routes:

```text
/login
  -> frontend/src/pages/Login.jsx

/register
  -> frontend/src/pages/Register.jsx

/
  -> frontend/src/routes/PrivateRoute.jsx
  -> if AuthContext is loading, show loading state
  -> if user exists, render frontend/src/pages/Dashboard.jsx
  -> if user is missing, redirect to /login
```

## Auth Check On Page Load

```text
frontend/src/App.jsx
  -> frontend/src/context/AuthContext.jsx
  -> checks localStorage for token
  -> if no token, loading ends and user stays null
  -> if token exists, calls authApi.me()
  -> frontend/src/api/authApi.js
  -> frontend/src/api/axios.js attaches Authorization: Bearer <token>
  -> GET http://localhost:8080/api/auth/me
  -> backend Gin router
  -> internal/middleware/cors.go
  -> internal/middleware/logging.go
  -> internal/middleware/recovery.go
  -> internal/middleware/auth.go verifies JWT and sets userID in Gin context
  -> internal/controller/auth_controller.go Me
  -> internal/service/auth_service.go Me
  -> internal/repository/user_repository.go
  -> internal/db/sqlc/user.sql.go
  -> MySQL users table
  -> response returns user
  -> AuthContext stores user in React state
  -> PrivateRoute allows Dashboard
```

## Register Flow

```text
User submits Register form
  -> frontend/src/pages/Register.jsx
  -> AuthContext.register()
  -> frontend/src/api/authApi.js register()
  -> frontend/src/api/axios.js sends POST /auth/register
  -> POST http://localhost:8080/api/auth/register
  -> internal/routes/routes.go public auth route
  -> internal/controller/auth_controller.go Register
  -> internal/util/response.go decodes JSON through Gin
  -> internal/service/auth_service.go Register
  -> internal/util/validator.go validates input
  -> internal/util/password.go hashes password
  -> internal/repository/user_repository.go creates user
  -> internal/db/sqlc/user.sql.go runs INSERT
  -> MySQL users table
  -> internal/util/jwt.go creates JWT
  -> response returns token and user
  -> AuthContext saves token in localStorage and stores user
  -> frontend redirects user into Dashboard
```

## Login Flow

```text
User submits Login form
  -> frontend/src/pages/Login.jsx
  -> AuthContext.login()
  -> frontend/src/api/authApi.js login()
  -> frontend/src/api/axios.js sends POST /auth/login
  -> POST http://localhost:8080/api/auth/login
  -> internal/routes/routes.go public auth route
  -> internal/controller/auth_controller.go Login
  -> internal/service/auth_service.go Login
  -> internal/repository/user_repository.go loads user by email
  -> internal/db/sqlc/user.sql.go queries MySQL
  -> internal/util/password.go checks password
  -> internal/util/jwt.go creates JWT
  -> response returns token and user
  -> AuthContext saves token in localStorage and stores user
  -> frontend redirects user into Dashboard
```

## Dashboard Load Flow

```text
PrivateRoute allows Dashboard
  -> frontend/src/pages/Dashboard.jsx renders
  -> Dashboard calls loadGroups()
  -> frontend/src/api/groupApi.js listGroups()
  -> frontend/src/api/axios.js attaches JWT
  -> GET /api/groups
  -> Gin middleware chain
  -> Auth middleware sets userID
  -> internal/controller/group_controller.go List
  -> internal/service/group_service.go List
  -> internal/repository/group_repository.go
  -> internal/db/sqlc/group.sql.go
  -> MySQL task_groups table
  -> groups returned to Dashboard

Dashboard also calls loadTasks()
  -> frontend/src/api/taskApi.js listTasks({ status, group_id })
  -> GET /api/tasks
  -> Gin middleware chain
  -> internal/controller/task_controller.go List
  -> internal/util/pagination.go reads page and limit
  -> internal/service/task_service.go List
  -> internal/repository/task_repository.go
  -> internal/db/sqlc/task.sql.go
  -> MySQL tasks table
  -> tasks returned to Dashboard
  -> Dashboard passes tasks to TaskList
```

Dashboard component flow:

```text
frontend/src/pages/Dashboard.jsx
  -> frontend/src/components/Navbar.jsx
  -> frontend/src/components/GroupSidebar.jsx
  -> frontend/src/components/TaskForm.jsx
  -> frontend/src/components/TaskList.jsx
  -> frontend/src/components/TaskCard.jsx
```

## Create Group Flow

```text
User creates a group in GroupSidebar
  -> frontend/src/components/GroupSidebar.jsx
  -> Dashboard onCreate handler
  -> frontend/src/api/groupApi.js createGroup()
  -> POST /api/groups
  -> Gin middleware chain
  -> internal/controller/group_controller.go Create
  -> internal/service/group_service.go Create
  -> internal/repository/group_repository.go
  -> internal/db/sqlc/group.sql.go
  -> MySQL task_groups table
  -> created group returned
  -> Dashboard reloads group list
```

## Create Task Flow

```text
User submits TaskForm
  -> frontend/src/components/TaskForm.jsx
  -> Dashboard onCreate handler
  -> frontend/src/api/taskApi.js createTask()
  -> POST /api/tasks
  -> Gin middleware chain
  -> internal/controller/task_controller.go Create
  -> internal/service/task_service.go Create
  -> validates title, priority, status, optional group, optional due date
  -> internal/repository/task_repository.go
  -> internal/db/sqlc/task.sql.go
  -> MySQL tasks table
  -> created task returned
  -> Dashboard reloads task list
```

## Update Or Complete Task Flow

```text
User toggles or edits a task
  -> frontend/src/components/TaskCard.jsx
  -> frontend/src/components/TaskList.jsx
  -> Dashboard onToggle/onUpdate handler
  -> frontend/src/api/taskApi.js updateTask()
  -> PUT /api/tasks/:id
  -> Gin middleware chain
  -> internal/controller/task_controller.go Update
  -> reads id from Gin route param
  -> internal/service/task_service.go Update
  -> checks user owns task
  -> applies completed_at when status becomes completed
  -> internal/repository/task_repository.go
  -> internal/db/sqlc/task.sql.go
  -> MySQL tasks table
  -> updated task returned
  -> Dashboard reloads task list
```

## Delete Task Flow

```text
User deletes a task
  -> frontend/src/components/TaskCard.jsx
  -> frontend/src/components/TaskList.jsx
  -> Dashboard onDelete handler
  -> frontend/src/api/taskApi.js deleteTask()
  -> DELETE /api/tasks/:id
  -> Gin middleware chain
  -> internal/controller/task_controller.go Delete
  -> internal/service/task_service.go Delete
  -> internal/repository/task_repository.go
  -> internal/db/sqlc/task.sql.go
  -> MySQL tasks table
  -> response returns message
  -> Dashboard reloads task list
```

## Reminder Flow

```text
User creates a reminder for a task
  -> frontend/src/components/TaskCard.jsx
  -> frontend/src/components/TaskList.jsx
  -> Dashboard onReminder handler
  -> frontend/src/api/taskApi.js createReminder()
  -> POST /api/reminders
  -> Gin middleware chain
  -> internal/controller/reminder_controller.go Create
  -> internal/service/reminder_service.go Create
  -> checks task belongs to current user
  -> internal/repository/reminder_repository.go
  -> internal/db/sqlc/reminder.sql.go
  -> MySQL reminders table
  -> response returns reminder
```

Background reminder scheduler:

```text
backend/cmd/server/main.go
  -> internal/scheduler/reminder_scheduler.go Start
  -> every minute calls reminder service
  -> internal/service/reminder_service.go ProcessDue
  -> internal/repository/reminder_repository.go finds due reminders
  -> internal/db/sqlc/reminder.sql.go queries MySQL
  -> scheduler logs due reminders
  -> repository marks reminders as sent
```

## Backend Request Flow Pattern

Most protected backend requests follow this shape:

```text
HTTP request from React
  -> internal/routes/routes.go Gin route
  -> internal/middleware/cors.go
  -> internal/middleware/logging.go
  -> internal/middleware/recovery.go
  -> internal/middleware/auth.go for protected routes
  -> internal/controller/*_controller.go
  -> internal/service/*_service.go
  -> internal/repository/*_repository.go
  -> internal/db/sqlc/*.sql.go
  -> MySQL
```

Most responses return in reverse:

```text
MySQL row
  -> sqlc model
  -> repository result
  -> service response DTO
  -> controller JSON response
  -> Axios response
  -> React state update
  -> UI rerender
```

## Important Files By Layer

Frontend:

```text
frontend/src/index.js                React mount point
frontend/src/App.jsx                 App routes and AuthProvider
frontend/src/context/AuthContext.jsx Auth state, login, register, logout, current user
frontend/src/routes/PrivateRoute.jsx Protected dashboard route
frontend/src/api/axios.js            Axios base URL and JWT header
frontend/src/api/authApi.js          Auth API calls
frontend/src/api/groupApi.js         Group API calls
frontend/src/api/taskApi.js          Task and reminder API calls
frontend/src/pages/Login.jsx         Login screen
frontend/src/pages/Register.jsx      Register screen
frontend/src/pages/Dashboard.jsx     Main task workspace
frontend/src/components/*.jsx        UI pieces used by Dashboard
```

Backend:

```text
backend/cmd/server/main.go                 Backend composition root
backend/internal/routes/routes.go          Gin router and route groups
backend/internal/middleware/*.go           CORS, logging, recovery, JWT auth
backend/internal/controller/*.go           Gin handlers
backend/internal/service/*.go              Business logic
backend/internal/repository/*.go           Database access wrapper
backend/internal/db/sqlc/*.go              sqlc-generated query layer
backend/internal/db/queries/*.sql          sqlc source queries
backend/internal/db/migrations/*.sql       MySQL schema migrations
backend/internal/scheduler/*.go            Background reminder processing
backend/internal/util/*.go                 JWT, password, validation, response helpers
```

