BEGIN;

DROP TABLE IF EXISTS allocation;

DROP INDEX IF EXISTS idx_allocation_portfolio_id;
DROP INDEX IF EXISTS idx_allocation_asset_id;

COMMIT;
