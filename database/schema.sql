CREATE TABLE "users"(
    "id" SERIAL NOT NULL,
    "username" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "group" TEXT NOT NULL,
    "created" TIMESTAMPTZ NOT NULL
);
ALTER TABLE
    "users" ADD PRIMARY KEY("id");
CREATE TABLE "groups"(
    "name" TEXT NOT NULL,
    "color" TEXT NOT NULL,
    "displayname" TEXT NULL,
    "weight" BIGINT NOT NULL
);
ALTER TABLE
    "groups" ADD PRIMARY KEY("name");
CREATE TABLE "permissions"(
    "permission" TEXT NOT NULL,
    "group" TEXT NOT NULL
);
ALTER TABLE
    "permissions" ADD PRIMARY KEY("permission");
ALTER TABLE
    "permissions" ADD CONSTRAINT "permissions_group_foreign" FOREIGN KEY("group") REFERENCES "groups"("name");
ALTER TABLE
    "users" ADD CONSTRAINT "users_group_foreign" FOREIGN KEY("group") REFERENCES "groups"("name");

