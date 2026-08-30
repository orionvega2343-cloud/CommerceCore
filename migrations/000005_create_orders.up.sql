CREATE TYPE order_status AS ENUM('created', 'paid', 'shipped', 'cancelled', 'completed');
CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    cart_id INT NOT NULL REFERENCES cart(id),
    status order_status NOT NULL,
    total_amount INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);