-- name: AddDocumentation :one
INSERT INTO "documentation" (
    "ownerId", "name", "public", "repoOwner", "repoName", "repoRef",
    "highlightStyle", "root", "fullPage", "useTailwind"
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING "id";

-- name: UpdateDocumentation :exec
UPDATE "documentation" SET
  "name" = $2,
  "public" = $3,
  "repoOwner" = $4,
  "repoName" = $5,
  "repoRef" = $6,
  "highlightStyle" = $7,
  "root" = $8,
  "fullPage" = $9,
  "useTailwind" = $10
WHERE "id" = $1;

-- name: TransferDocumentation :exec
UPDATE "documentation" SET "ownerId" = $2 WHERE "id" = $1;

-- name: GetDocumentationByName :one
SELECT * FROM "documentation" WHERE "name" = $1;

-- name: GetDocumentationById :one
SELECT * FROM "documentation" WHERE "id" = $1;

-- name: GetUserDocumentations :many
SELECT * FROM "documentation" WHERE "ownerId" = $1;

-- name: RemoveDocumentation :exec
DELETE FROM "documentation" WHERE "id" = $1;
