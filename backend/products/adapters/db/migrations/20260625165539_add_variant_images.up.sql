BEGIN;

CREATE TABLE products.variant_images (
	image_uuid    uuid         NOT NULL,
	variant_uuid  uuid         NOT NULL,
	url           varchar(1024) NOT NULL,
	position      integer      NOT NULL DEFAULT 0,
	created_at    TIMESTAMPTZ  NOT NULL,
	PRIMARY KEY (image_uuid),
	FOREIGN KEY (variant_uuid) REFERENCES products.variants (variant_uuid) ON DELETE CASCADE
);

CREATE INDEX idx_variant_images_variant_uuid ON products.variant_images (variant_uuid);

COMMIT;
