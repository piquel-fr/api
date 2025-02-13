-- name: GetGroupInfo :one
SELECT * FROM "groups" WHERE "name" = $1;

-- name: GetUserGroup :one
SELECT "group" FROM "users" WHERE id = $1;

-- name: GetInheritedGroups :many
SELECT "name" FROM "groups"
    WHERE "weight" <= (
        SELECT "weight"
        FROM "users"
        INNER JOIN "groups" ON "group" = "groups"."name"
        WHERE "id" = $1
    );
