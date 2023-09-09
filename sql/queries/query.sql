-- name: CreateEvent :exec
INSERT INTO
    event (
        id,
        title,
        description,
        timezone,
        created_by,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateSchedule :exec
INSERT INTO
    schedule (
        id,
        event_id,
        start_time,
        "duration",
        is_full_day,
        recurring_interval,
        recurring_type
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateInvitation :exec
INSERT INTO
    invitation (id, event_id, user_id, token, status)
VALUES
    ($1, $2, $3, $4, $5);

-- name: DeleteEvent :exec
DELETE FROM
    event
WHERE
    id = $1;

-- name: UpdateEvent :exec
UPDATE
    event
SET
    title = $1,
    description = $2,
    timezone = $3,
    updated_at = $4;

-- name: UpsertSchedule :exec
INSERT INTO
    schedule (
        id,
        event_id,
        start_time,
        "duration",
        is_full_day,
        recurring_interval,
        recurring_type
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id, event_id) DO
UPDATE
SET
    start_time = $3,
    "duration" = $4,
    recurring_interval = $5,
    recurring_type = $6;

-- name: UpsertInvitation :exec
INSERT INTO
    invitation (id, event_id, user_id, token, status)
VALUES
    ($1, $2, $3, $4, $5) ON CONFLICT (id, event_id) DO
UPDATE
SET
    user_id = $3,
    token = $4,
    status = $5;

-- name: FindEventByID :one
SELECT
    *
FROM
    event
WHERE
    id = $1
LIMIT
    1;

-- name: FindSchedulesByEventID :many
SELECT
    *
FROM
    schedule
WHERE
    event_id = $1;

-- name: FindInvitationsByEventID :many
SELECT
    *
FROM
    invitation
WHERE
    event_id = $1;