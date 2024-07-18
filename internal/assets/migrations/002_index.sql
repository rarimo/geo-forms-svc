-- +migrate Up

CREATE INDEX IF NOT EXISTS balances_status_index ON balances (status);
CREATE INDEX IF NOT EXISTS balances_nullifier_index ON balances (nullifier);

-- +migrate Down

