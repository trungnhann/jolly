BEGIN;

CREATE SCHEMA IF NOT EXISTS orders;

CREATE TABLE orders.orders
(
	order_uuid    uuid         NOT NULL,
	customer_id   varchar(255) NOT NULL,
	currency      varchar(3)   NOT NULL,
	status        varchar(32)  NOT NULL,
	total_cents   bigint       NOT NULL,
	placed_at_utc TIMESTAMPTZ  NOT NULL,
	PRIMARY KEY (order_uuid)
);

CREATE TABLE orders.order_line_items
(
	line_item_uuid    uuid         NOT NULL,
	order_uuid        uuid         NOT NULL,
	sku               varchar(255) NOT NULL,
	quantity          integer      NOT NULL,
	unit_price_cents  bigint       NOT NULL,
	PRIMARY KEY (line_item_uuid),
	FOREIGN KEY (order_uuid) REFERENCES orders.orders (order_uuid) ON DELETE CASCADE
);

CREATE INDEX idx_order_line_items_order_uuid
	ON orders.order_line_items (order_uuid);

COMMIT;
