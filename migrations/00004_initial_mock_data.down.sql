BEGIN;

-- Delete from allocation table
DELETE FROM allocation WHERE portfolio_id = 1;

-- Delete from asset table
DELETE FROM asset WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10);

-- Delete from portfolio table
DELETE FROM portfolio WHERE id = 1;

COMMIT;
