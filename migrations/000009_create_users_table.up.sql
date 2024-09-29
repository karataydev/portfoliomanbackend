BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    google_id VARCHAR(255),
    profile_picture_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

INSERT INTO users (first_name, last_name, email, google_id, profile_picture_url) VALUES
('Emre', 'karatay', 'marufemre99@gmail.com', '123456789', 'https://example.com/john_doe.jpg');

-- update portfolio
ALTER TABLE portfolio
ADD CONSTRAINT fk_portfolio_user
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;

COMMIT;
