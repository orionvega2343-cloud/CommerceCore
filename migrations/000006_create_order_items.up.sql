CREATE  TABLE IF NOT EXISTS order_items(
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES order(id),
    product_id INT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL CHECK ( quantity > 0 ),
    price_per_unit INT NOT NULL CHECK ( price_per_unit > 0 )
)