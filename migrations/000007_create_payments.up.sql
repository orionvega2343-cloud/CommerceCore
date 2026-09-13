CREATE TYPE payment_status AS ENUM('pending', 'succeeded', 'failed');
CREATE TABLE IF NOT EXISTS payments(
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id),
    amount INT NOT NULL CHECK ( amount > 0 ),
    status payment_status NOT NULL,
    method TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)