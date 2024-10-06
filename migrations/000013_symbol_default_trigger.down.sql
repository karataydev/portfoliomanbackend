BEGIN;

-- Drop the trigger
DROP TRIGGER IF EXISTS set_symbol_trigger ON portfolio;

-- Drop the function
DROP FUNCTION IF EXISTS set_symbol();

-- Restore the default value for the symbol column
ALTER TABLE portfolio
ALTER COLUMN symbol SET DEFAULT 'PF' || lpad(nextval('portfolio_id_seq'::regclass)::text, 8, '0');

COMMIT;
