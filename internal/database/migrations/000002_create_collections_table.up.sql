-- Create collections table
CREATE TABLE
    collections (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name VARCHAR NOT NULL,
        description VARCHAR,
        created_by UUID,
        updated_by UUID,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

-- Create users_collections table
CREATE TABLE
    users_collections (
        user_id UUID UNIQUE,
        collection_id UUID,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        PRIMARY KEY (user_id, collection_id)
    );

-- Create indexes
CREATE INDEX users_collections_user_id ON users_collections (user_id);

-- Add foreign key constraints
ALTER TABLE users_collections ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE users_collections ADD CONSTRAINT fk_collection_id FOREIGN KEY (collection_id) REFERENCES collections (id);