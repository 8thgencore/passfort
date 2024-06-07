-- Create users_role_enum type
CREATE TYPE "users_role_enum" AS ENUM ('admin', 'user');

-- Create users table
CREATE TABLE
  users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    master_password VARCHAR,
    salt BYTEA,
    is_verified BOOLEAN DEFAULT false NOT NULL,
    role users_role_enum DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
    UNIQUE (email)
  );

CREATE INDEX users_email ON users (email);