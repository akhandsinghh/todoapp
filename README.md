# Todo App

Full-stack todo application using only the requested stack:

- Frontend: React with JavaScript and Webpack
- Backend: Go with Gin
- Database: MySQL
- ORM/query layer: sqlc-style generated Go code under `backend/internal/db/sqlc`

## Project Structure

The app follows the requested architecture, with backend controllers, services, repositories, middleware, models, routes, scheduler, migrations, sqlc queries, and a React/Webpack frontend.

## Detailed Documentation

Separate documentation is available here:

- `docs/BACKEND.md` explains every backend folder/file and the Go API request flow.
- `docs/FRONTEND.md` explains every frontend folder/file and the React application flow.
- `docs/DATABASE.md` explains migrations, tables, sqlc query files, generated sqlc files, and database flow.

## Prerequisites

Install these first:

- Go 1.22+
- Node.js 18+
- MySQL 8+
- Optional: sqlc, if you want to regenerate the files in `backend/internal/db/sqlc`

## Database Setup

Create the MySQL database:

```sql
CREATE DATABASE todo_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Update `backend/.env` if your MySQL username, password, host, port, or database name are different.

The backend automatically runs the SQL files in `backend/internal/db/migrations` on startup.

## Start Backend

```bash
cd backend
go mod tidy
go run ./cmd/server
```

The backend starts on `http://localhost:8080` by default.

Health check:

```bash
curl http://localhost:8080/health
```

## Start Frontend

Open a second terminal:

```bash
cd frontend
npm install
npm start
```

The frontend starts on `http://localhost:3000`.

## Main API Routes

Public routes:

- `POST /api/auth/register`
- `POST /api/auth/login`

Authenticated routes require `Authorization: Bearer <token>`:

- `GET /api/auth/me`
- `GET /api/groups`
- `POST /api/groups`
- `PUT /api/groups/{id}`
- `DELETE /api/groups/{id}`
- `GET /api/tasks`
- `POST /api/tasks`
- `PUT /api/tasks/{id}`
- `DELETE /api/tasks/{id}`
- `GET /api/reminders`
- `POST /api/reminders`
- `DELETE /api/reminders/{id}`

## Regenerate sqlc Code

The `backend/sqlc.yaml` file is included. If sqlc is installed, run:

```bash
cd backend
sqlc generate
```

The repository currently includes generated-style files so the app is complete even before regeneration.

## Notes

- Dates from the frontend are sent as RFC3339 timestamps.
- The reminder scheduler runs every minute and logs due reminders, then marks them as sent.
- For development, CORS is configured for `http://localhost:3000` in `backend/.env`.
