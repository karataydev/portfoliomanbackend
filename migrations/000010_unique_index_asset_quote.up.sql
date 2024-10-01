BEGIN;

ALTER TABLE asset_quote
ADD CONSTRAINT idx_unique_asset_quote_asset_id_quote_time UNIQUE (asset_id, quote_time);

COMMIT;
