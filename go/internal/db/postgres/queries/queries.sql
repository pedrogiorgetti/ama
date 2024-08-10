-- name: GetRoom :one
SELECT 
    "id", "name", "created_at", "updated_at"
FROM room
WHERE "id" = $1;

-- name: GetRooms :many
SELECT 
    "id", "name", "created_at", "updated_at"
FROM room;

-- name: CreateRoom :one
INSERT INTO room 
  ("name") VALUES
  ($1)
RETURNING "id", "name", "created_at", "updated_at";

-- name: GetQuestion :one
SELECT
    "id", "room_id", "text", "reaction_count", "answered", "created_at", "updated_at"
FROM question
WHERE "id" = $1;

-- name: GetRoomQuestions :many
SELECT
    "id", "room_id", "text", "reaction_count", "answered", "created_at", "updated_at"
FROM question
WHERE "room_id" = $1;

-- name: CreateQuestion :one
INSERT INTO question 
  ("room_id", "text")
  VALUES ($1, $2)
RETURNING "id", "room_id", "text", "reaction_count", "answered", "created_at", "updated_at";

-- name: ReactToQuestion :one
UPDATE question
SET
    "reaction_count" = "reaction_count" + 1
WHERE "id" = $1
RETURNING "reaction_count";

-- name: RemoveReactionFromQuestion :one
UPDATE question
SET
    "reaction_count" = "reaction_count" - 1
WHERE "id" = $1
RETURNING "reaction_count";

-- name: MarkQuestionAsAnswered :exec
UPDATE question
SET
    "answered" = true
WHERE "id" = $1;
