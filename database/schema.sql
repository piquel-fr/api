CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY NOT NULL,
    "username" TEXT NOT NULL UNIQUE,
    "name" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "user_sessions" (
    "id" SERIAL PRIMARY KEY NOT NULL,
    "userId" INTEGER REFERENCES "users" ("id") NOT NULL,
    "tokenHash" VARCHAR(255) NOT NULL,
    "userAgent" TEXT NOT NULL,
    "ipAdress" VARCHAR(45) NOT NULL,
    "expiresAt" TIMESTAMPTZ NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
