BEGIN;

ALTER TABLE asset_quote
DROP CONSTRAINT IF EXISTS idx_unique_asset_quote_asset_id_quote_time;

COMMIT;
