CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY NOT NULL,
    "username" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "documentation" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "ownerId" SERIAL REFERENCES "users" ("id") NOT NULL,
  
  "repoOwner" TEXT NOT NULL,
  "repoName" TEXT NOT NULL,
  "repoRef" TEXT NOT NULL,
  
  "highlightStyle" VARCHAR(127) NOT NULL DEFAULT "tokyonight",
  "root" VARCHAR(127) NOT NULL DEFAULT "/",
  "fullPage" BOOLEAN NOT NULL DEFAULT FALSE,
  "useTailwind" BOOLEAN NOT NULL DEFAULT FALSE,
  "public" BOOLEAN NOT NULL DEFAULT FALSE
);
