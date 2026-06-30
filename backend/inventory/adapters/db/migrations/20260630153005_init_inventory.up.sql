BEGIN;

CREATE SCHEMA IF NOT EXISTS inventory;

CREATE TABLE inventory.stocks (
    stock_uuid UUID NOT NULL,
    sku VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (stock_uuid),
    CONSTRAINT stocks_sku_unique UNIQUE (sku),
    CONSTRAINT quantity_non_negative CHECK (quantity >= 0)
);

CREATE TABLE inventory.reservations (
    reservation_uuid UUID NOT NULL,
    order_uuid UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (reservation_uuid),
    CONSTRAINT reservations_order_uuid_unique UNIQUE (order_uuid)
);

CREATE TABLE inventory.reservation_items (
    reservation_uuid UUID NOT NULL,
    stock_uuid UUID NOT NULL,
    quantity INTEGER NOT NULL,
    PRIMARY KEY (reservation_uuid, stock_uuid),
    FOREIGN KEY (reservation_uuid) REFERENCES inventory.reservations (reservation_uuid) ON DELETE CASCADE,
    FOREIGN KEY (stock_uuid) REFERENCES inventory.stocks (stock_uuid) ON DELETE CASCADE,
    CONSTRAINT quantity_positive CHECK (quantity > 0)
);

CREATE INDEX idx_reservation_items_reservation_uuid ON inventory.reservation_items (reservation_uuid);
CREATE INDEX idx_reservations_order_uuid ON inventory.reservations (order_uuid);

-- Seed stock levels for testing (generating a static v4 UUID for the seeded SKU-123)
INSERT INTO inventory.stocks (stock_uuid, sku, quantity)
VALUES ('c1a5b8d2-4e9f-4b0a-8a1c-9d8e7f6a5b4c', 'SKU-123', 100)
ON CONFLICT (sku) DO NOTHING;

COMMIT;
