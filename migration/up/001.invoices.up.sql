

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


CREATE INDEX IF NOT EXISTS ON invoices(user_id);
CREATE INDEX IF NOT EXISTS ON invoices(status);

CREATE OR REPLACE TRIGGER invoices_updated_at
    BEFORE UPDATE ON invoices
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE UNIQUE INDEX IF NOT EXISTS idx_invoices_one_paid_per_order
    ON invoices(order_id)
    WHERE status = 'paid';
