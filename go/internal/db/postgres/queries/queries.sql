-- name: GetRoom :one
SELECT 
    "id", "theme", "created_at", "updated_at"
FROM room
WHERE "id" = $1;

-- name: GetRooms :many
SELECT 
    "id", "theme", "created_at", "updated_at"
FROM room;

-- name: CreateRoom :one
INSERT INTO room 
  ("theme") VALUES
  ($1)
RETURNING "id", "theme", "created_at", "updated_at";

-- name: GetMessage :one
SELECT
    "id", "room_id", "text", "reaction_count", "answered", "created_at", "updated_at"
FROM message
WHERE "id" = $1;

-- name: GetRoomMessages :many
SELECT
    "id", "room_id", "text", "reaction_count", "answered", "created_at", "updated_at"
FROM message
WHERE "room_id" = $1;

-- name: CreateMessage :one
INSERT INTO message 
  ("room_id", "text")
  VALUES ($1, $2)
RETURNING "id", "room_id", "text", "reaction_count", "answered", "created_at", "updated_at";

-- name: ReactToMessage :one
UPDATE message
SET
    "reaction_count" = "reaction_count" + 1
WHERE "id" = $1
RETURNING "reaction_count";

-- name: RemoveReactionFromMessage :one
UPDATE message
SET
    "reaction_count" = "reaction_count" - 1
WHERE "id" = $1
RETURNING "reaction_count";

-- name: MarkMessageAsAnswered :exec
UPDATE message
SET
    "answered" = true
WHERE "id" = $1;
