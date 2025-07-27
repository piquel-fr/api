-- name: AddDocumentation :one
INSERT INTO "documentation" (
    "ownerId", "name", "public", "repoOwner", "repoName", "repoRef",
    "root", "pathPrefix", "highlightStyle", "fullPage", "useTailwind"
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING "id";

-- name: UpdateDocumentation :exec
UPDATE "documentation" SET
    "name" = $2,
    "public" = $3,
    "repoOwner" = $4,
    "repoName" = $5,
    "repoRef" = $6,
    "root" = $7,
    "pathPrefix" = $8,
    "highlightStyle" = $9,
    "fullPage" = $10,
    "useTailwind" = $11
WHERE "id" = $1;

-- name: TransferDocumentation :exec
UPDATE "documentation" SET "ownerId" = $2 WHERE "id" = $1;

-- name: GetDocumentationByName :one
SELECT * FROM "documentation" WHERE "name" = $1;

-- name: GetDocumentationById :one
SELECT * FROM "documentation" WHERE "id" = $1;

-- name: ListUserDocsInstances :many
SELECT * FROM "documentation" WHERE "ownerId" = $1 LIMIT $2 OFFSET $3;

-- name: CountUserDocsInstances :one
SELECT COUNT(*) FROM "documentation" WHERE "ownerId" = $1;

-- name: RemoveDocumentation :exec
DELETE FROM "documentation" WHERE "id" = $1;
