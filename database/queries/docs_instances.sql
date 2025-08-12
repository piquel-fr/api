-- name: AddDocsInstance :one
INSERT INTO "docs_instances" (
    "ownerId", "name", "public", "repoOwner", "repoName", "repoRef",
    "root", "pathPrefix", "highlightStyle", "fullPage", "useTailwind"
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING "id";

-- name: UpdateDocsInstance :exec
UPDATE "docs_instances" SET
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

-- name: GetDocsInstanceByName :one
SELECT * FROM "docs_instances" WHERE "name" = $1;

-- name: GetDocsInstanceById :one
SELECT * FROM "docs_instances" WHERE "id" = $1;

-- name: ListUserDocsInstances :many
SELECT * FROM "docs_instances" WHERE "ownerId" = $1 LIMIT $2 OFFSET $3;

-- name: CountUserDocsInstances :one
SELECT COUNT(*) FROM "docs_instances" WHERE "ownerId" = $1;

-- name: RemoveDocsInstance :exec
DELETE FROM "docs_instances" WHERE "id" = $1;
