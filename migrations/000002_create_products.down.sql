CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price float8 NOT NULL CHECK ( price > 0 ),
    stock_quantity INT NOT NULL,
    is_active BOOLEAN NOT NULL
)