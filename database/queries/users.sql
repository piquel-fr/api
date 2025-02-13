-- name: AddUser :exec
INSERT INTO "users" ("username", "name", "image", "email", "group", "created")
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUserByUsername :one
SELECT * FROM "users" WHERE "username" = $1;

-- name: GetUserById :one
SELECT * FROM "users" WHERE "id" = $1;

-- name: GetUserByEmail :one
SELECT * FROM "users" WHERE "email" = $1;

-- name: UpdateUser :exec
UPDATE "users" SET "username" = $2, "name" = $3, "image" = $4 WHERE "id" = $1;
