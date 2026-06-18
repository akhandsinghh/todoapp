# Frontend Documentation

## Purpose

The frontend is a React application written in JavaScript and bundled by Webpack. It provides screens for login, registration, task management, groups, task filtering, completion toggling, deletion, and reminders.

Node is only used as the frontend tool runtime because Webpack, npm, and React build tooling run on Node. The backend remains Go-only.

## How The Frontend Starts

Frontend source starts from:

```text
frontend/src/index.js
```

Startup flow:

```text
npm start
  -> webpack-dev-server
  -> public/index.html
  -> src/index.js
  -> src/App.jsx
  -> AuthProvider
  -> BrowserRouter
  -> Login/Register/Dashboard routes
```

Default frontend URL:

```text
http://localhost:3000
```

## Frontend Control Flow

### Initial Page Load

```text
public/index.html
  -> creates <div id="root">
  -> Webpack injects bundled JavaScript
  -> src/index.js mounts React
  -> App.jsx configures routes
  -> AuthContext checks localStorage token
  -> PrivateRoute allows Dashboard or redirects to Login
```

### Login Flow

```text
Login.jsx
  -> user submits email/password
  -> AuthContext.login()
  -> api/authApi.js
  -> api/axios.js
  -> POST http://localhost:8080/api/auth/login
  -> token saved in localStorage
  -> user is redirected to Dashboard
```

### Task Creation Flow

```text
Dashboard.jsx
  -> TaskForm.jsx
  -> taskApi.createTask()
  -> axios client adds Authorization header
  -> Go backend /api/tasks
  -> Dashboard reloads task list
  -> TaskList.jsx renders TaskCard.jsx items
```

### Group Flow

```text
Dashboard.jsx
  -> GroupSidebar.jsx
  -> groupApi.createGroup()
  -> Go backend /api/groups
  -> Dashboard reloads groups
  -> active group filter reloads tasks
```

## Frontend File Purpose

### Root Frontend Files

`frontend/package.json`

Defines frontend dependencies and scripts. `npm start` runs Webpack dev server. `npm run build` creates a production bundle.

`frontend/webpack.config.js`

Webpack configuration. It sets the React entry file, output bundle, Babel loader, CSS loader, HTML template, dev server, and API proxy.

`frontend/babel.config.js`

Babel configuration for compiling modern JavaScript and JSX.

### `public`

`frontend/public/index.html`

HTML shell loaded by the browser. React mounts into `<div id="root"></div>`.

### `src`

`frontend/src/index.js`

React entrypoint. It finds `#root` in `index.html` and renders `<App />`.

`frontend/src/App.jsx`

Top-level React component. It wires `AuthProvider`, `BrowserRouter`, public routes, and protected dashboard routing.

### `src/api`

API modules keep HTTP calls separate from UI components.

`axios.js`

Creates the Axios client with base URL `http://localhost:8080/api`. It also attaches the JWT token from `localStorage` to every request.

`authApi.js`

Contains login, register, and current-user API calls.

`taskApi.js`

Contains task API calls and reminder API calls used from task screens.

`groupApi.js`

Contains group API calls.

### `src/context`

`AuthContext.jsx`

Central authentication state. It stores the current user, checks existing tokens on page load, exposes `login`, `register`, and `logout`, and saves/removes the JWT token in `localStorage`.

### `src/routes`

`PrivateRoute.jsx`

Protects dashboard routes. If auth is still loading, it shows a loading message. If no user is logged in, it redirects to `/login`.

### `src/pages`

Pages are route-level screens.

`Login.jsx`

Login form. Calls `AuthContext.login()` and redirects to `/` after successful login.

`Register.jsx`

Registration form. Calls `AuthContext.register()` and redirects to `/` after successful registration.

`Dashboard.jsx`

Main authenticated application screen. Loads groups and tasks, manages filters, handles task creation, deletion, completion toggling, reminder creation, and group creation.

### `src/components`

Components are reusable UI blocks used by pages.

`Navbar.jsx`

Top bar showing the app name, logged-in user name, and logout button.

`GroupSidebar.jsx`

Sidebar for selecting all tasks or a group. Also contains the new-group form and color swatches.

`TaskForm.jsx`

Task creation form. Collects title, notes, priority, group, and due date.

`TaskList.jsx`

Receives tasks and renders a list of `TaskCard` components. Shows an empty state when there are no tasks.

`TaskCard.jsx`

Displays one task. Allows completion toggling, reminder creation, and deletion.

### `src/styles`

`global.css`

Global reset, base typography, form styles, buttons, auth screens, alerts, and loading states.

`dashboard.css`

Dashboard layout, navbar, sidebar, toolbar, task form, task list, task cards, reminder row, responsive behavior.

## Running Frontend

```bash
cd frontend
npm install
npm start
```

Production build:

```bash
npm run build
```

## Frontend To Backend Communication

The frontend calls the Go backend through Axios:

```text
React component
  -> api module
  -> axios.js
  -> Go backend route
```

Authenticated calls include:

```text
Authorization: Bearer <token>
```

The token is received from login/register and stored in:

```text
localStorage["token"]
```

