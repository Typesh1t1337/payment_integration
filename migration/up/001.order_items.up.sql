

CREATE TABLE IF NOT EXISTS order_items(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    product_id uuid NOT NULL REFRENECES products(id) ON DELETE RESTRICT,
    QUANTITY integer NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE UNIQUE INDEX IF NOT EXISTS ON order_items(order_id, product_id);

CREATE OR REPLACE TRIGGER order_items_updated_at
    BEFORE UPDATE ON order_items
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
