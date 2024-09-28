BEGIN;

-- Delete from allocation table
DELETE FROM allocation WHERE portfolio_id = 1;

-- Delete from asset table
DELETE FROM asset WHERE id > 0;

-- Delete from portfolio table
DELETE FROM portfolio WHERE id = 1;

COMMIT;
