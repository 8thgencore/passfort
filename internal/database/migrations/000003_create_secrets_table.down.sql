-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_check_secret_link ON secrets;

DROP FUNCTION IF EXISTS check_secret_link;

-- Drop tables
DROP TABLE IF EXISTS secrets;

DROP TABLE IF EXISTS text_secrets;

DROP TABLE IF EXISTS password_secrets;

-- Drop enums
DROP TYPE IF EXISTS secret_type_enum;