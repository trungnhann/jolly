BEGIN;

CREATE TABLE products.categories (
    category_uuid uuid NOT NULL,
    parent_category_uuid uuid,
    name varchar(255) NOT NULL,
    slug varchar(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (category_uuid),
    CONSTRAINT categories_slug_unique UNIQUE (slug),
    CONSTRAINT fk_categories_parent FOREIGN KEY (parent_category_uuid) REFERENCES products.categories (category_uuid) ON DELETE SET NULL
);

CREATE TABLE products.brands (
    brand_uuid uuid NOT NULL,
    name varchar(255) NOT NULL,
    slug varchar(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (brand_uuid),
    CONSTRAINT brands_slug_unique UNIQUE (slug)
);

ALTER TABLE products.products ADD COLUMN category_uuid uuid;
ALTER TABLE products.products ADD COLUMN brand_uuid uuid;

ALTER TABLE products.products 
    ADD CONSTRAINT fk_products_category FOREIGN KEY (category_uuid) REFERENCES products.categories (category_uuid) ON DELETE SET NULL;

ALTER TABLE products.products 
    ADD CONSTRAINT fk_products_brand FOREIGN KEY (brand_uuid) REFERENCES products.brands (brand_uuid) ON DELETE SET NULL;

COMMIT;
