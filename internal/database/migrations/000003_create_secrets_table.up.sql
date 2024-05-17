-- Create secret_type_enum type
CREATE TYPE "secret_type_enum" AS ENUM ('password', 'text', 'file');

-- Create the secrets table
CREATE TABLE
    secrets (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        collection_id UUID REFERENCES collections (id) ON DELETE CASCADE,
        secret_type secret_type_enum NOT NULL,
        name varchar NOT NULL,
        description varchar,
        created_by UUID REFERENCES users (id),
        updated_by UUID REFERENCES users (id),
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );