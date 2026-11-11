# Stride Backend API Reference

> Base URL: `http://<host>:<port>/api/v1`
> Health probe:`GET /healthz` (returns `{"status":"ok"}`)

This document describes every HTTP endpoint, the request/response shapes, and the JSON envelope.

---

## 1. Conventions

### 1.1 Envelope

Every JSON response follows the same shape:

```jsonc
// success
{ "data": <payload> }

// error
{ "error": { "code": "<machine_code>", "message": "<human readable>" } }
```

`<machine_code>` is one of:

| Code                | HTTP | Meaning                                          |
| ------------------- | ---- | ------------------------------------------------ |
| `invalid_argument`  | 400  | Invalid input format / validation error.         |
| `unauthenticated`   | 401  | Missing or invalid token.                        |
| `permission_denied` | 403  | Authenticated but not allowed for this resource. |
| `not_found`         | 404  | Resource does not exist for this user.           |
| `internal`          | 500  | Unhandled error.                                 |

### 1.2 Authentication

All routes under `/api/v1` (except `/auth/register` and `/auth/login`) require
a JWT bearer token:

```
Authorization: Bearer <token>
```

`token` is the JWT returned by `/auth/register` and `/auth/login`. It encodes
`user_id` and `username` and expires after `JWT_EXPIRES_HOURS` hours
(default 72).

### 1.3 IDs

All IDs are UUIDv4 strings. Composite primary keys (e.g. `goal_metric_rel`)
include `evaluated_at` so the same metric can be evaluated multiple times.

### 1.4 Timestamps

All timestamps are RFC 3339 (`2024-01-02T15:04:05Z`).

---

## 2. Auth

### POST `/auth/register`

```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "super-secret"
}
```

`201 Created`:

```jsonc
{
  "data": {
    "user": {
      "id": "<uuid>",
      "username": "alice",
      "email": "alice@example.com",
    },
    "token": "<jwt>",
  },
}
```

Errors: `400` for malformed input, `409` if username or email is taken.

### POST `/auth/login`

```json
{ "identifier": "alice OR alice@example.com", "password": "super-secret" }
```

`200 OK`: same envelope shape as register. `401` if credentials are wrong.

### GET `/auth/me`

`200 OK`:

```json
{
  "data": { "id": "<uuid>", "username": "alice", "email": "alice@example.com" }
}
```

---

## 3. Goals

### POST `/goals`

```jsonc
{
  "title": "Read 24 books",
  "description": "12 novels, 12 non-fiction", // optional
  "expected_start_time": "2024-01-01T00:00:00Z", // optional
  "expected_end_time": "2024-12-31T23:59:59Z", // optional
  "parent_goal_id": "<uuid>", // optional
}
```

`201 Created` returns the created goal.

### PATCH `/goals/:id`

All fields are optional. To clear `parent_goal_id`, send `"parent_goal_id": null`.

```jsonc
{
  "title": "...",
  "description": "...",
  "expected_start_time": "...",
  "expected_end_time": "...",
  "parent_goal_id": null,
}
```

Returns the updated goal. `400` if you try to make a goal its own parent.

### GET `/goals`

Query parameters:

| Name              | Type | Meaning                                         |
| ----------------- | ---- | ----------------------------------------------- |
| `completed`       | bool | Filter by `finished_at IS NOT NULL` if `true`.  |
| `deleted`         | bool | Only soft-deleted goals if `true`.              |
| `include_deleted` | bool | Include soft-deleted goals alongside live ones. |
| `parent_goal_id`  | UUID | Only goals with this parent.                    |

Returns `{ "data": [Goal, ...] }`.

### GET `/goals/:id`

Returns one goal or 404. Soft-deleted goals are returned if `include_deleted=true`
in the query string.

### DELETE `/goals/:id`

Soft delete. `204 No Content`.

### POST `/goals/:id/restore`

Restore a soft-deleted goal. Returns the restored goal.

### DELETE `/goals/:id/permanent`

Permanently delete a goal (rows in `goals` and dependent join tables go away).
`204 No Content`.

### POST `/goals/:id/complete`

Stamps `finished_at = now()`. Returns the goal.

### POST `/goals/:id/uncomplete`

Clears `finished_at`. Returns the goal.

## 4. Tasks

Same lifecycle as goals (`POST/PATCH/GET/DELETE/restore/permanent/complete/uncomplete`), plus the following.

### Body (POST `/tasks` and PATCH `/tasks/:id`)

```jsonc
{
  "title": "Wire the database",
  "description": "...",
  "expected_start_time": "...",
  "expected_end_time": "...",
  "goal_id": "<uuid>", // null/omit to detach
  "parent_task_id": "<uuid>", // null/omit to detach

  // patch-only convenience flags
  "clear_goal": true,
  "clear_parent": true,
}
```

`clear_goal` / `clear_parent` exist because JSON cannot distinguish between
"the client sent `null`" and "the client omitted the field" in many languages.
Sending them as `true` sets the corresponding FK to `NULL`.

### Query parameters on GET `/tasks`

| Name              | Type | Meaning                                         |
| ----------------- | ---- | ----------------------------------------------- |
| `goal_id`         | UUID | Only tasks belonging to this goal.              |
| `no_goal`         | bool | Only tasks **without** a goal.                  |
| `parent_task_id`  | UUID | Only sub-tasks of this parent.                  |
| `completed`       | bool | Filter by completion state.                     |
| `deleted`         | bool | Only soft-deleted tasks.                        |
| `include_deleted` | bool | Include soft-deleted tasks alongside live ones. |

---

## 5. Tags

There are two parallel families: **goal tags** and **task tags**. The endpoints
are symmetric — only the prefix differs.

### Catalogue (per user)

| Method | Path             | Body                |
| ------ | ---------------- | ------------------- |
| POST   | `/goal-tags`     | `{ "name": "..." }` |
| GET    | `/goal-tags`     | —                   |
| DELETE | `/goal-tags/:id` | —                   |

(Substitute `task-tags` for the task variant.)

---
