BEGIN;

ALTER TABLE asset_quote
DROP COLUMN IF EXISTS symbol;

DROP INDEX IF EXISTS idx_portfolio_symbol;

COMMIT;
