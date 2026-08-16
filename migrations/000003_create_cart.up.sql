CREATE TYPE status_type AS ENUM('active', 'checked_out', 'abandoned');
CREATE TABLE IF NOT EXISTS cart(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    status status_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_cart_active_user ON cart(user_id) WHERE status = 'active';