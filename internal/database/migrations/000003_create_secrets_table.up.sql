-- Create secret_type_enum type
CREATE TYPE "secret_type_enum" AS ENUM ('password', 'text', 'file');


-- Create the secrets table
CREATE TABLE IF NOT EXISTS secrets (
    id UUID PRIMARY KEY,
    collection_id UUID NOT NULL,
    secret_type secret_type_enum NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL
);
