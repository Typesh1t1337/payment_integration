-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS users(
   id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   name        VARCHAR(255) NOT NULL,
   email       VARCHAR(255) NOT NULL UNIQUE,
   password    TEXT NOT NULL,
   created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE OR REPLACE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE products (
      id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      name        VARCHAR(255) NOT NULL,
      price       NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
      created_at  TIMESTAMP NOT NULL DEFAULT now(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE OR REPLACE TRIGGER products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS orders(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id)  ON DELETE RESTRICT,
    status TEXT NOT NULL,
    locked_until TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE OR REPLACE TRIGGER orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();


CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_one_created_per_user
    ON orders (user_id)
    WHERE status = 'created';


CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_order_status ON orders(status);

CREATE TABLE IF NOT EXISTS order_items(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity integer NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE UNIQUE INDEX IF NOT EXISTS idx_order_items_order_product ON order_items(order_id, product_id);

CREATE OR REPLACE TRIGGER order_items_updated_at
    BEFORE UPDATE ON order_items
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE invoices (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    order_id      UUID REFERENCES orders(id) ON DELETE SET NULL,
    status       VARCHAR(20) NOT NULL,
    total_amount NUMERIC(12, 2) NOT NULL CHECK (total_amount >= 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    paid_at TIMESTAMPTZ
);


CREATE INDEX IF NOT EXISTS idx_invoices_user_id ON invoices(user_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);

CREATE OR REPLACE TRIGGER invoices_updated_at
    BEFORE UPDATE ON invoices
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE UNIQUE INDEX IF NOT EXISTS idx_invoices_one_paid_per_order
    ON invoices(order_id)
    WHERE status = 'paid';
-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS invoices;
DROP FUNCTION IF EXISTS set_updated_at();