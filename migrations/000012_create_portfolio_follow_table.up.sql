BEGIN;

CREATE TABLE IF NOT EXISTS portfolio_follow (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    portfolio_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_portfolio_follow_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_portfolio_follow_portfolio
        FOREIGN KEY (portfolio_id)
        REFERENCES portfolio(id)
        ON DELETE CASCADE,
    CONSTRAINT uq_portfolio_follow_user_portfolio
        UNIQUE (user_id, portfolio_id)
);

-- Add indexes to improve query performance
CREATE INDEX IF NOT EXISTS idx_portfolio_follow_user_id ON portfolio_follow(user_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_follow_portfolio_id ON portfolio_follow(portfolio_id);

-- Add a comment to the table
COMMENT ON TABLE portfolio_follow IS 'Stores relationships between users and the portfolios they follow';

COMMIT;