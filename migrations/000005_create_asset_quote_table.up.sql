BEGIN;

CREATE TABLE IF NOT EXISTS asset_quote (
    id BIGSERIAL PRIMARY KEY,
    asset_id BIGINT NOT NULL,
    quote DOUBLE PRECISION NOT NULL,
    quote_time TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (asset_id) REFERENCES asset(id)
);

CREATE INDEX IF NOT EXISTS idx_asset_quote_asset_id ON asset_quote(asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_quote_quote_time ON asset_quote(quote_time);

COMMIT;
