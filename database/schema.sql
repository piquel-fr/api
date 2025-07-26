CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY NOT NULL,
    "username" TEXT NOT NULL UNIQUE,
    "name" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "documentation" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "ownerId" SERIAL REFERENCES "users" ("id") NOT NULL,
  "name" TEXT NOT NULL UNIQUE,
  "public" BOOLEAN NOT NULL DEFAULT FALSE,
  
  "repoOwner" TEXT NOT NULL,
  "repoName" TEXT NOT NULL,
  "repoRef" TEXT NOT NULL,
  
  "highlightStyle" VARCHAR(127) NOT NULL DEFAULT "tokyonight",
  "root" VARCHAR(127) NOT NULL DEFAULT "/",
  "fullPage" BOOLEAN NOT NULL DEFAULT FALSE,
  "useTailwind" BOOLEAN NOT NULL DEFAULT FALSE
);
