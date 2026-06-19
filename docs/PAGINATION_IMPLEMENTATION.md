# Pagination and Sorting Implementation

## What changed

### Backend
- `backend/internal/db/queries/task.sql`
  - Extended `ListTasksByUser` query with dynamic sorting fields and order.
  - Added `CountTasksByUser` query to return the total task count for pagination metadata.

- `backend/internal/controller/task_controller.go`
  - Updated `TaskController.List` to read new query params:
    - `sort_by`
    - `sort_order`
    - `page`
    - `limit`
  - Returned a paginated response object with `tasks`, `total`, `page`, and `limit`.

- `backend/internal/service/task_service.go`
  - Added validation for sort values.
  - Added a count call for total matching tasks.
  - Passed pagination and sort parameters to the repository.

- `backend/internal/repository/task_repository.go`
  - Updated `List` to use the generated `ListTasksByUserParams` type.
  - Added `Count` for `CountTasksByUser`.

- `backend/internal/model/response.go`
  - Added `TaskListResponse` to shape the paginated JSON response.

### Frontend
- `frontend/src/pages/Dashboard.jsx`
  - Added state for pagination and sorting:
    - `page`
    - `limit`
    - `sortBy`
    - `sortOrder`
    - `totalTasks`
  - Updated `loadTasks()` to pass pagination and sorting query params to the API.
  - Added UI controls for:
    - sort field
    - sort order
    - page size
    - page navigation buttons
  - Updated `useEffect` hooks so tasks are reloaded only when:
    - the active group changes
    - the status filter changes
    - the page changes
    - the sort field changes
    - the sort order changes
    - the page size changes
  - Updated mutation handlers (`create`, `update`, `delete`) to reload only after the user performs a task mutation.

## How frequent fetches were reduced

- The prior implementation fetched tasks on every mutation using full reload logic, and the UI state could be updated in a way that caused extra refreshes.
- With pagination and sorting, the component now keeps explicit state for page, limit, sort, and order.
- `loadTasks()` is only called when:
  - filters or sorting/page controls change
  - a task is created, updated, or deleted
- This means the app does not refetch tasks on every render; it refetches only on real user actions.
- For task mutations, the code now reloads the current page after the mutation so the paged data stays consistent with the server.

## Notes

- `frontend/src/api/taskApi.js` already supported passing query params, so no API helper change was needed for pagination.
- The new backend response shape is a paginated object, but `Dashboard.jsx` also handles raw array responses safely.

## Result

- Users can now select page size and sort order.
- Tasks are displayed page by page.
- The app shows total tasks and page navigation.
- Fetching is limited to filter/sort/page changes and task mutations, instead of redundant loads.
