-- +migrate Up

CREATE INDEX IF NOT EXISTS forms_status_index ON forms (status);
CREATE INDEX IF NOT EXISTS forms_nullifier_index ON forms (nullifier);

-- +migrate Down

