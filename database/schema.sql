CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY NOT NULL,
    "username" TEXT NOT NULL UNIQUE,
    "name" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "mail_accounts" (
    "id" SERIAL PRIMARY KEY NOT NULL,
    "ownerId" SERIAL REFERENCES "users" ("id") NOT NULL,
    "email" TEXT NOT NULL UNIQUE,
    "name" TEXT NOT NULL,
    "username" TEXT NOT NULL,
    "password" TEXT NOT NULL
);

CREATE TABLE "mail_share" (
    "userId" INTEGER REFERENCES "users" ("id") NOT NULL,
    "account" INTEGER REFERENCES "mail_accounts" ("id") NOT NULL,
    "permission" TEXT NOT NULL,
    UNIQUE ("userId", "account")
);
