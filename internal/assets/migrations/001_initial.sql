-- +migrate Up
CREATE OR REPLACE FUNCTION trigger_set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
AS $$ BEGIN NEW.updated_at = (NOW() AT TIME ZONE 'utc'); RETURN NEW; END; $$;

CREATE TABLE IF NOT EXISTS forms
(
    id           uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    nullifier    TEXT NOT NULL,
    status       TEXT NOT NULL,
    name         TEXT NOT NULL,
    surname      TEXT NOT NULL,
    id_num       TEXT NOT NULL,
    birthday     TEXT NOT NULL,
    citizen      TEXT NOT NULL,
    visited      TEXT NOT NULL,
    purpose      TEXT NOT NULL,
    country      TEXT NOT NULL,
    city         TEXT NOT NULL,
    address      TEXT NOT NULL,
    postal       TEXT NOT NULL,
    phone        TEXT NOT NULL,
    email        TEXT NOT NULL,
    image        TEXT NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at   TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE INDEX IF NOT EXISTS forms_status_index ON forms (status);
CREATE INDEX IF NOT EXISTS forms_nullifier_index ON forms (nullifier);

DROP TRIGGER IF EXISTS set_updated_at ON forms;
CREATE TRIGGER set_updated_at
    BEFORE UPDATE
    ON forms
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

-- +migrate Down

DROP TABLE IF EXISTS forms;

DROP FUNCTION IF EXISTS trigger_set_updated_at();
