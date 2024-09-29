BEGIN;

CREATE TABLE IF NOT EXISTS transaction (
    id BIGSERIAL PRIMARY KEY,
    side SMALLINT NOT NULL CHECK (side IN (0, 1)),
    quantity DECIMAL(18, 8) NOT NULL,
    price DECIMAL(18, 8) NOT NULL,
    allocation_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (allocation_id) REFERENCES allocation(id)
);

CREATE INDEX IF NOT EXISTS idx_transaction_allocation_id ON transaction(allocation_id);

COMMIT;
