

CREATE TABLE IF NOT EXISTS orders(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id)  ON DELETE RESTRICT,
    status TEXT NOT NULL,
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