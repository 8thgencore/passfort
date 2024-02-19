-- Create collections table
CREATE TABLE "collections" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(), 
    "name" varchar NOT NULL,
    "description" varchar,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- Create users_collections table
CREATE TABLE "users_collections" (
    "user_id" uuid REFERENCES users(id),
    "collection_id" uuid REFERENCES collections(id),
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, collection_id)
);

-- Create indexes
CREATE INDEX users_collections_user_id ON users_collections(user_id);

-- Add foreign key constraints
ALTER TABLE users_collections ADD CONSTRAINT fk_users_collections_user_id
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE users_collections ADD CONSTRAINT fk_users_collections_collection_id
    FOREIGN KEY (collection_id) REFERENCES collections(id);
