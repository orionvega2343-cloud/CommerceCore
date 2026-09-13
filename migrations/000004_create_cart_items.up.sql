CREATE TABLE IF NOT EXISTS cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL REFERENCES cart(id),
    product_id INT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL CHECK ( quantity > 0 ),
    price_snapshot INT NOT NULL
)