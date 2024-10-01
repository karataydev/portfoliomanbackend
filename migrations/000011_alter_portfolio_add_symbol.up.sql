BEGIN;

ALTER TABLE portfolio
ADD COLUMN symbol VARCHAR(20) UNIQUE NOT NULL DEFAULT 'PF' || lpad(nextval('portfolio_id_seq'::regclass)::text, 8, '0');

CREATE INDEX IF NOT EXISTS idx_portfolio_symbol ON portfolio(symbol);


COMMIT;
