-- name: AddSession :one
INSERT INTO "user_sessions" ("userId", "tokenHash", "userAgent", "ipAdress", "expiresAt")
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetSessionFromHash :one
SELECT * FROM "user_sessions" WHERE "tokenHash" = $1;

-- name: UpdateSession :exec
UPDATE "user_sessions" SET "tokenHash" = $2, "expiresAt" = $3 WHERE "userId" = $1;

-- name: DeleteSessionById :exec
DELETE FROM "user_sessions" WHERE "userId" = $1 AND "id" = $2;

-- name: DeleteSessionByHash :exec
DELETE FROM "user_sessions" WHERE "tokenHash" = $1;

-- name: GetUserSessions :many
SELECT * FROM "user_sessions" WHERE "userId" = $1 ORDER BY "id" ASC;

-- name: ClearUserSessions :exec
DELETE FROM "user_sessions" WHERE "userId" = $1;
