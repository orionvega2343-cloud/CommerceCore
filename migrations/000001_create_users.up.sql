CREATE TYPE select_role AS ENUM('customer', 'admin');
CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    role select_role NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)