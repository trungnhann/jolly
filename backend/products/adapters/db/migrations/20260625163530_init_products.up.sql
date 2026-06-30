BEGIN;

CREATE SCHEMA IF NOT EXISTS products;

CREATE TABLE products.products (
	product_uuid  uuid         NOT NULL,
	name          varchar(255) NOT NULL,
	description   text         NOT NULL,
	status        varchar(32)  NOT NULL,
	created_at    TIMESTAMPTZ  NOT NULL,
	updated_at    TIMESTAMPTZ  NOT NULL,
	PRIMARY KEY (product_uuid)
);

CREATE TABLE products.variants (
	variant_uuid  uuid         NOT NULL,
	product_uuid  uuid         NOT NULL,
	sku           varchar(255) NOT NULL,
	name          varchar(255) NOT NULL,
	price_cents   bigint       NOT NULL,
	created_at    TIMESTAMPTZ  NOT NULL,
	updated_at    TIMESTAMPTZ  NOT NULL,
	PRIMARY KEY (variant_uuid),
	CONSTRAINT variants_sku_unique UNIQUE (sku),
	FOREIGN KEY (product_uuid) REFERENCES products.products (product_uuid) ON DELETE CASCADE
);

CREATE INDEX idx_variants_product_uuid ON products.variants (product_uuid);

COMMIT;
