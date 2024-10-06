BEGIN;

CREATE OR REPLACE FUNCTION set_symbol()
RETURNS TRIGGER AS $$
BEGIN
  NEW.symbol := 'PF' || lpad(NEW.id::text, 8, '0');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_symbol_trigger
BEFORE INSERT ON portfolio
FOR EACH ROW
EXECUTE FUNCTION set_symbol();

ALTER TABLE portfolio ALTER COLUMN symbol DROP DEFAULT;

COMMIT;
