# Group Sharing Feature

## Purpose

Group sharing allows a group creator to share a task group with another registered user by email. It does not use websockets. Shared users see the group in their sidebar and can work with tasks inside that group through normal HTTP API calls.

The UI labels groups like this:

```text
creator - group belongs to the logged-in user
shared  - group was shared with the logged-in user
```

## How It Works

The feature adds a join table named `group_shares`.

```text
task_groups
  -> group owner is stored in task_groups.user_id

group_shares
  -> stores which extra users can access a group
```

When a creator shares a group, the backend:

```text
GroupSidebar.jsx
  -> groupApi.shareGroup(groupId, { email })
  -> POST /api/groups/:id/share
  -> GroupController.Share
  -> GroupService.Share
  -> checks the current user owns the group
  -> finds the target user by email
  -> inserts group_id + user_id into group_shares
```

## Backend Behavior

### Group List

`GET /api/groups` now returns both owned and shared groups. Each group includes a `role` field:

```json
{
  "id": 1,
  "user_id": 4,
  "name": "Work",
  "color": "#4f46e5",
  "role": "creator"
}
```

For shared groups, `role` is `"shared"`.

### Sharing Rules

Only the creator can share a group. Shared users cannot re-share that group.

The backend rejects:

```text
empty email
sharing a group the current user does not own
sharing with a user that does not exist
sharing with yourself
```

### Task Access In Shared Groups

Task access was updated so users can access:

```text
their own ungrouped tasks
their own tasks
tasks inside groups they created
tasks inside groups shared with them
```

This means if a shared user creates a task inside the creator's group, the creator can also see that task.

Group update and delete remain creator-only.

## Frontend Behavior

The sidebar now shows:

```text
All tasks
Ungrouped
Group name + creator/shared label
Share by email form for creator groups only
```

Shared groups can still be selected and used for task filtering and task creation. The share form is hidden for groups where the logged-in user is not the creator.

## Files Changed

### Database

`backend/internal/db/migrations/005_group_shares.sql`

Added the `group_shares` table with foreign keys to `task_groups` and `users`.

`backend/internal/db/queries/group.sql`

Added queries for:

```text
CreateGroupShare
ListAccessibleGroups
GetAccessibleGroupByID
```

`backend/internal/db/queries/task.sql`

Updated task access queries so tasks can be read, updated, and deleted when they are inside accessible shared groups.

### Generated SQL Layer

`backend/internal/db/sqlc/group.sql.go`

Added generated-style methods and structs for accessible groups and group sharing.

`backend/internal/db/sqlc/task.sql.go`

Updated generated-style task query methods for shared-group access.

`backend/internal/db/sqlc/querier.go`

Added the new query methods to the `Querier` interface.

### Backend Models

`backend/internal/model/group.go`

Added `Role` to `GroupDTO`.

`backend/internal/model/request.go`

Added `ShareGroupRequest` with an `email` field.

### Backend Repository

`backend/internal/repository/group_repository.go`

Added repository methods for sharing a group, listing accessible groups, loading accessible group details, and finding a user by email.

`backend/internal/repository/task_repository.go`

Added access to `GetAccessibleGroupByID` so task creation/update can validate group access.

### Backend Service

`backend/internal/service/group_service.go`

Added group role conversion and `Share` business logic.

`backend/internal/service/task_service.go`

Added group-access validation for tasks assigned to a group.

### Backend Controller And Routes

`backend/internal/controller/group_controller.go`

Added `Share` handler for `POST /api/groups/:id/share`.

`backend/internal/routes/routes.go`

Registered the new share route:

```text
POST /api/groups/:id/share
```

### Frontend API

`frontend/src/api/groupApi.js`

Added:

```js
shareGroup(id, payload)
```

### Frontend UI

`frontend/src/components/GroupSidebar.jsx`

Added:

```text
creator/shared labels
Ungrouped sidebar button
owner-only share-by-email form
```

`frontend/src/pages/Dashboard.jsx`

Passed the share handler into `GroupSidebar` and reset the current page when changing group filters.

`frontend/src/styles/dashboard.css`

Added styles for group labels, share form, and grouped sidebar rows.

## What Was Not Changed

No websockets were added.

Group update and delete are still creator-only.

The dashboard's local task add/update/delete logic was kept intact.
