BEGIN;

-- First, drop the index
DROP INDEX IF EXISTS idx_transaction_allocation_id;

-- Then drop the table
DROP TABLE IF EXISTS transaction;

COMMIT;
