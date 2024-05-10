-- Create users_role_enum type
CREATE TYPE "users_role_enum" AS ENUM ('admin', 'user');

-- Create users table
CREATE TABLE
  users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    master_password VARCHAR,
    is_verified BOOLEAN DEFAULT false NOT NULL,
    role users_role_enum DEFAULT 'user',
    created_at timestamptz NOT NULL DEFAULT (now ()),
    updated_at timestamptz NOT NULL DEFAULT (now ()),
    UNIQUE (email)
  );