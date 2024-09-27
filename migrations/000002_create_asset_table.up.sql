BEGIN;

CREATE TABLE IF NOT EXISTS asset (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

CREATE INDEX IF NOT EXISTS idx_asset_name ON asset(name);
CREATE INDEX IF NOT EXISTS idx_asset_symbol ON asset(symbol);

COMMIT;
