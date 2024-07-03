-- +migrate Up

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
    created_at   TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc')
);

-- +migrate Down

DROP TABLE IF EXISTS forms;
