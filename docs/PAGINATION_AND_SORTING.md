# Pagination And Sorting Feature

## Purpose

Pagination lets the user choose how many tasks are shown on the main dashboard at one time. Sorting lets the user order tasks by priority or deadline date in ascending or descending order.

This was implemented on the main dashboard only, so it works for:

```text
All tasks
Ungrouped tasks
Any selected group
Pending filter
Completed filter
Shared groups
Creator groups
```

## User Flow

On the dashboard, the user can choose:

```text
Page size: 5, 10, 20, or 50
Sort by: Deadline date or Priority
Order: Ascending or Descending
```

Then:

```text
Load   - applies selected page size and sort options
Reload - reloads the current page with the current applied options
Prev   - goes to the previous page
Next   - goes to the next page
```

Example:

```text
20 total tasks
page size = 10

Page 1 shows tasks 1-10
Page 2 shows tasks 11-20
```

## Backend Implementation

The backend already had a pagination helper:

```text
backend/internal/util/pagination.go
```

It reads:

```text
page
limit
```

and converts them to:

```text
LIMIT
OFFSET
```

The task list endpoint was extended to also read:

```text
sort_by
sort_order
ungrouped
```

Supported sort values:

```text
sort_by=priority
sort_by=due_at
```

Supported order values:

```text
sort_order=asc
sort_order=desc
```

The backend normalizes invalid sort options. If `sort_by` is not supported, the query falls back to the default due-date/created-date ordering.

## API Response Shape

`GET /api/tasks` now returns an object:

```json
{
  "items": [],
  "total": 20
}
```

`items` contains only the current page of tasks.

`total` contains the total number of matching tasks before pagination. The frontend uses this to calculate how many pages exist.

## Sorting Details

### Deadline Date

Deadline sorting uses `due_at`.

Tasks without a deadline are still included. The SQL query uses a fallback date so null deadlines do not break sorting.

### Priority

Priority sorting uses this order for ascending:

```text
low
medium
high
```

For descending:

```text
high
medium
low
```

## Ungrouped Filtering

An `ungrouped` query flag was added so the backend can return only tasks where `group_id IS NULL`.

This was needed because client-only filtering would make pagination totals wrong. For example, if the server returned 10 mixed tasks and the frontend filtered out grouped tasks locally, the page might show fewer than 10 tasks and the total page count would be incorrect.

## Frontend Implementation

`Dashboard.jsx` now stores two sets of control state:

```text
controls - what the user has currently selected in the dropdowns
query    - what has actually been applied by clicking Load
```

This means changing a dropdown does not immediately fetch tasks. The task list reloads only when:

```text
Load is clicked
Reload is clicked
page changes
status changes
group filter changes
```

This keeps the existing dashboard behavior of avoiding unnecessary task-list fetches.

Local task handling was preserved:

```text
addLocalTask
updateLocalTask
deleteLocalTask
```

Delete still removes the task from local state after the API delete succeeds instead of forcing a full list fetch.

## Files Changed

### Backend Controller

`backend/internal/controller/task_controller.go`

Changed `List` to read:

```text
group_id
ungrouped
status
sort_by
sort_order
page
limit
```

Then it passes these values to the task service.

### Backend Service

`backend/internal/service/task_service.go`

Changed task listing to:

```text
normalize sort options
count total matching tasks
fetch only the requested page
return TaskListResponse
```

Also added group access validation for shared-group task behavior.

### Backend Repository

`backend/internal/repository/task_repository.go`

Added:

```text
Count
GetAccessibleGroup
```

`Count` supports pagination totals.

### Backend Models

`backend/internal/model/response.go`

Added:

```go
type TaskListResponse struct {
    Items []TaskDTO `json:"items"`
    Total int64     `json:"total"`
}
```

### SQL Query Files

`backend/internal/db/queries/task.sql`

Updated task listing query to support:

```text
pagination
total count
priority sorting
deadline sorting
ungrouped filtering
shared-group access
```

### Generated SQL Layer

`backend/internal/db/sqlc/task.sql.go`

Updated generated-style code for:

```text
ListTasksByUser
CountTasksByUser
GetTaskByID
UpdateTask
DeleteTask
```

`backend/internal/db/sqlc/querier.go`

Added `CountTasksByUser` to the query interface.

### Frontend API

`frontend/src/api/taskApi.js`

No function name changed. The existing `listTasks(params)` function is used with more query parameters.

### Frontend Dashboard

`frontend/src/pages/Dashboard.jsx`

Added:

```text
page state
totalTasks state
page size control
sort-by control
sort-order control
Load button
Reload button
Prev/Next buttons
ungrouped filter handling
support for { items, total } task responses
```

Kept:

```text
local add task behavior
local update task behavior
local delete task behavior
```

### Frontend Sidebar

`frontend/src/components/GroupSidebar.jsx`

Added the `Ungrouped` button so pagination can work for ungrouped tasks as a server-side filter.

### Frontend Styles

`frontend/src/styles/dashboard.css`

Added styles for:

```text
pagination toolbar
page controls
responsive pagination layout
```

## What Was Not Changed

No websocket logic was added.

The existing task API module function names were not changed.

The dashboard still avoids unnecessary task-list fetches by preserving the local add/update/delete logic.
