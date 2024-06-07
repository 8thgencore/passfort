-- Create secret_type_enum type
CREATE TYPE "secret_type_enum" AS ENUM ('password', 'text', 'file');

-- Create the secrets table
CREATE TABLE
    secrets (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        collection_id UUID REFERENCES collections (id) ON DELETE CASCADE,
        secret_type secret_type_enum NOT NULL,
        name VARCHAR NOT NULL,
        description VARCHAR,
        created_by UUID,
        updated_by UUID,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        linked_secret_id UUID
    );

ALTER TABLE secrets ADD CONSTRAINT fk_collection_id FOREIGN KEY (collection_id) REFERENCES collections (id);

-- Create password_secrets table
CREATE TABLE
    password_secrets (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        url VARCHAR NOT NULL,
        login VARCHAR NOT NULL,
        password BYTEA NOT NULL
    );

-- Create text_secrets table
CREATE TABLE
    text_secrets (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        text BYTEA NOT NULL
    );

-- Function to check secret links
CREATE OR REPLACE FUNCTION check_secret_link() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.secret_type = 'password' THEN
        PERFORM 1 FROM password_secrets WHERE id = NEW.linked_secret_id;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'Invalid linked_secret_id for password secret';
        END IF;
    ELSIF NEW.secret_type = 'text' THEN
        PERFORM 1 FROM text_secrets WHERE id = NEW.linked_secret_id;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'Invalid linked_secret_id for text secret';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for checking secret links
CREATE TRIGGER trigger_check_secret_link
BEFORE INSERT OR UPDATE ON secrets
FOR EACH ROW
EXECUTE FUNCTION check_secret_link();


-- Function to cascade delete linked secrets
CREATE OR REPLACE FUNCTION cascade_delete_linked_secret() RETURNS TRIGGER AS $$
BEGIN
    IF OLD.secret_type = 'password' THEN
        DELETE FROM password_secrets WHERE id = OLD.linked_secret_id;
    ELSIF OLD.secret_type = 'text' THEN
        DELETE FROM text_secrets WHERE id = OLD.linked_secret_id;
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Trigger for cascading delete linked secrets
CREATE TRIGGER trigger_cascade_delete_linked_secret
AFTER DELETE ON secrets
FOR EACH ROW
EXECUTE FUNCTION cascade_delete_linked_secret();