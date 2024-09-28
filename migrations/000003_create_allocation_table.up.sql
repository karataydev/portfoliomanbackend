BEGIN;

CREATE TABLE IF NOT EXISTS allocation (
    id BIGSERIAL PRIMARY KEY,
    portfolio_id BIGINT NOT NULL,
    asset_id BIGINT NOT NULL,
    target_percentage DECIMAL(5,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_portfolio
        FOREIGN KEY(portfolio_id)
        REFERENCES portfolio(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_asset
        FOREIGN KEY(asset_id)
        REFERENCES asset(id)
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_allocation_portfolio_id ON allocation(portfolio_id);
CREATE INDEX IF NOT EXISTS idx_allocation_asset_id ON allocation(asset_id);

COMMIT;
