CREATE TYPE "users_role_enum" AS ENUM ('admin', 'user');

CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY, 
  "name" varchar not null, 
  "email" varchar not null, 
  "password" varchar not null, 
  "role" users_role_enum default 'user', 
  "created_at" timestamptz not null default (now()), 
  "updated_at" timestamptz not null default (now())
);

CREATE UNIQUE INDEX "email" on "users" ("email");
