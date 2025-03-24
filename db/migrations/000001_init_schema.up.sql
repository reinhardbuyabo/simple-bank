CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    owner VARCHAR NOT NULL,
    balance BIGINT NOT NULL,
    currency VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE entries (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id),
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE transfers (
    id BIGSERIAL PRIMARY KEY,
    from_account_id BIGINT NOT NULL REFERENCES accounts(id),
    to_account_id BIGINT NOT NULL REFERENCES accounts(id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_entries_account_id ON entries(account_id);
CREATE INDEX idx_transfers_from_account_id ON transfers(from_account_id);
CREATE INDEX idx_transfers_to_account_id ON transfers(to_account_id);
CREATE INDEX idx_transfers_from_to_account_id ON transfers(from_account_id, to_account_id);

COMMENT ON COLUMN "entries"."amount" IS 'can be positive or negative';
COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';