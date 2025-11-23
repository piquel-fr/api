-- name: AddEmailAccount :one
INSERT INTO "mail_accounts" (
    "ownerId", "email", "name", "username", "password"
)
VALUES ($1, $2, $3, $4, $5) RETURNING "id";

-- name: GetMailAccountByEmail :one
SELECT * FROM "mail_accounts" WHERE "email" = $1;

-- name: GetMailAccountById :one
SELECT * FROM "mail_accounts" WHERE "id" = $1;

-- name: ListMailAccounts :many
SELECT * FROM "mail_accounts" LIMIT $1 OFFSET $2;

-- name: ListUserMailAccounts :many
SELECT * FROM "mail_accounts" WHERE "ownerId" = $1 LIMIT $2 OFFSET $3;

-- name: CountUserMailAccounts :one
SELECT COUNT(*) FROM "mail_accounts" WHERE "ownerId" = $1;

-- name: RemoveMailAccount :exec
DELETE FROM "mail_accounts" WHERE "id" = $1;
