-- name: GetUserPermissions :many
SELECT "permission"
FROM "permissions"
WHERE "group" IN (
    SELECT "name"
    FROM "groups"
    WHERE "weight" <= (
        SELECT "weight"
        FROM "users"
        INNER JOIN "groups" ON "group" = "groups"."name"
        WHERE "id" = $1
    )
);
