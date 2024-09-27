BEGIN;

-- Insert into portfolio table
INSERT INTO portfolio (id, user_id, name, description, created_at, updated_at)
VALUES (1, 1, 'Hisse Senedi Portfoyum', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert into asset table
INSERT INTO asset (id, name, symbol, description)
VALUES
    (1, 'Apple Inc.', 'AAPL', 'Technology company'),
    (2, 'S&P 500 Index', 'SP500', 'Stock market index'),
    (3, 'MTV', 'MTV', 'Media company'),
    (4, 'Euro', 'EUR', 'European currency'),
    (5, 'Hepsiburada', 'HEPS', 'E-commerce company'),
    (6, 'US Dollar', 'USD', 'United States currency'),
    (7, 'Gold', 'GOLD', 'Precious metal'),
    (8, 'Alphabet Inc.', 'GOOGL', 'Technology company'),
    (9, 'Meta Platforms', 'META', 'Technology company'),
    (10, 'NVIDIA Corporation', 'NVDA', 'Technology company');

-- Insert into allocation table
INSERT INTO allocation (portfolio_id, asset_id, target_percentage)
VALUES
    (1, 1, 20.00),
    (1, 2, 30.00),
    (1, 3, 10.00),
    (1, 4, 5.00),
    (1, 5, 5.00),
    (1, 6, 20.00),
    (1, 7, 2.00),
    (1, 8, 2.00),
    (1, 9, 2.00),
    (1, 10, 4.00);

COMMIT;
