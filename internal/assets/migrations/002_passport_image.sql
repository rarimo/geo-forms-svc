-- +migrate Up
ALTER TABLE forms ADD COLUMN passport_image TEXT;

-- +migrate Down
ALTER TABLE forms DROP COLUMN passport_image;
