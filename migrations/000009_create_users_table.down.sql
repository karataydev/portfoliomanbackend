BEGIN;

ALTER TABLE portfolio
DROP CONSTRAINT IF EXISTS fk_portfolio_user;

-- Drop the trigger
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_google_id;

-- Drop the table
DROP TABLE IF EXISTS users;

COMMIT;
