-- name: AddEmailAccount :one
INSERT INTO "mail_accounts" (
    "ownerId", "email", "name", "username", "password"
)
VALUES ($1, $2, $3, $4, $5) RETURNING "id";

-- name: GetMailAccountByEmail :one
SELECT m.* FROM "mail_accounts" m
LEFT JOIN "mail_share" s ON m."id" = s."account"
WHERE m."email" = $1 
LIMIT 1;

-- name: GetMailAccountById :one
SELECT m.* FROM "mail_accounts" m
LEFT JOIN "mail_share" s ON m."id" = s."account"
WHERE m."id" = $1 
LIMIT 1;

-- name: ListUserMailAccounts :many
SELECT DISTINCT "mail_accounts".* FROM "mail_accounts"
LEFT JOIN "mail_share" ON "mail_accounts"."id" = "mail_share"."account"
WHERE "mail_accounts"."ownerId" = $1 OR "mail_share"."userId" = $1
ORDER BY "mail_accounts"."id";

-- name: CountUserMailAccounts :one
SELECT COUNT(DISTINCT "mail_accounts"."id")
FROM "mail_accounts"
LEFT JOIN "mail_share" ON "mail_accounts"."id" = "mail_share"."account"
WHERE "mail_accounts"."ownerId" = $1 OR "mail_share"."userId" = $1;

-- name: RemoveMailAccount :exec
DELETE FROM "mail_accounts" 
WHERE "id" = $1;

-- name: AddShare :exec
INSERT INTO "mail_share" (
    "userId", "account", "permission"
)
VALUES ($1, $2, $3);

-- name: RemoveShare :exec
DELETE FROM "mail_share"
WHERE "userId" = $1 AND "account" = $2;

-- name: ListAccountShares :many
SELECT "userId" FROM "mail_share" WHERE "account" = $1
ORDER BY "userId";
