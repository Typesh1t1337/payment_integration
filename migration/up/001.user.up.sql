

CREATE TABLE IF NOT EXISTS users(
   id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   name        VARCHAR(255) NOT NULL,
   email       VARCHAR(255) NOT NULL UNIQUE,
   created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE OR REPLACE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();



