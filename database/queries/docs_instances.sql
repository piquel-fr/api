-- name: AddDocsInstance :one
INSERT INTO "docs_instances" (
    "ownerId", "name", "public", "repoOwner", "repoName", "repoRef", "root"
)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING "id";

-- name: UpdateDocsInstance :exec
UPDATE "docs_instances" SET
    "name" = $2,
    "public" = $3,
    "repoOwner" = $4,
    "repoName" = $5,
    "repoRef" = $6,
    "root" = $7
WHERE "id" = $1;

-- name: GetDocsInstanceByName :one
SELECT * FROM "docs_instances" WHERE "name" = $1;

-- name: GetDocsInstanceById :one
SELECT * FROM "docs_instances" WHERE "id" = $1;

-- name: ListDocsInstances :many
SELECT * FROM "docs_instances" LIMIT $1 OFFSET $2;

-- name: ListUserDocsInstances :many
SELECT * FROM "docs_instances" WHERE "ownerId" = $1 LIMIT $2 OFFSET $3;

-- name: CountUserDocsInstances :one
SELECT COUNT(*) FROM "docs_instances" WHERE "ownerId" = $1;

-- name: RemoveDocsInstance :exec
DELETE FROM "docs_instances" WHERE "id" = $1;
