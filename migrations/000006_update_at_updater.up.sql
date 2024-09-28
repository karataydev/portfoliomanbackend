BEGIN;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_portfolio_updated_at
BEFORE UPDATE ON portfolio
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_allocation_updated_at
BEFORE UPDATE ON allocation
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMIT;
