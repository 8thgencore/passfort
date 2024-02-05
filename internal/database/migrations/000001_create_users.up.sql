-- Create users_role_enum type
CREATE TYPE "users_role_enum" AS ENUM ('admin', 'user');

-- Create users table
CREATE TABLE
  "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR not null,
    "email" VARCHAR not null,
    "password" VARCHAR not null,
    "role" users_role_enum default 'user',
    "created_at" timestamptz not null default (now ()),
    "updated_at" timestamptz not null default (now ())
    UNIQUE (email)
  );
