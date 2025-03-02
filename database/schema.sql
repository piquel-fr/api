CREATE TABLE "users"(
    "id" SERIAL NOT NULL,
    "username" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "created" TIMESTAMPTZ NOT NULL
);
ALTER TABLE
    "users" ADD PRIMARY KEY("id");
