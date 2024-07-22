-- +migrate Up

ALTER TABLE forms
    ADD COLUMN image_url TEXT;
ALTER TABLE forms 
    ALTER COLUMN image DROP NOT NULL;
ALTER TABLE forms
    ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc');


-- +migrate Down

ALTER TABLE forms
    DROP COLUMN image_url;
ALTER TABLE forms
    DROP COLUMN updated_at;
