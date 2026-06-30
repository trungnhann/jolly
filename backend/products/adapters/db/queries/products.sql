-- name: CreateProduct :exec
INSERT INTO products.products (
	product_uuid,
	name,
	description,
	status,
	created_at,
	updated_at
)
VALUES
	($1, $2, $3, $4, $5, $6)
;

-- name: GetProduct :one
SELECT
	*
FROM
	products.products
WHERE
	product_uuid = $1
LIMIT 1;

-- name: ListProducts :many
SELECT
	*
FROM
	products.products
ORDER BY
	created_at DESC;

-- name: UpdateProduct :exec
UPDATE products.products
SET
	name = $2,
	description = $3,
	status = $4,
	updated_at = $5
WHERE
	product_uuid = $1;

-- name: DeleteProduct :exec
DELETE FROM products.products
WHERE
	product_uuid = $1;

-- name: CreateVariant :exec
INSERT INTO products.variants (
	variant_uuid,
	product_uuid,
	sku,
	name,
	price_cents,
	created_at,
	updated_at
)
VALUES
	($1, $2, $3, $4, $5, $6, $7)
;

-- name: UpdateVariant :exec
UPDATE products.variants
SET
	sku = $2,
	name = $3,
	price_cents = $4,
	updated_at = $5
WHERE
	variant_uuid = $1;

-- name: DeleteVariant :exec
DELETE FROM products.variants
WHERE
	variant_uuid = $1;

-- name: GetVariantsForProduct :many
SELECT
	*
FROM
	products.variants
WHERE
	product_uuid = $1
ORDER BY
	created_at ASC;

-- name: GetVariantBySKU :one
SELECT
	*
FROM
	products.variants
WHERE
	sku = $1
LIMIT 1;

-- name: DeleteVariantsForProduct :exec
DELETE FROM products.variants
WHERE
	product_uuid = $1;

-- name: CreateVariantImage :exec
INSERT INTO products.variant_images (
	image_uuid,
	variant_uuid,
	url,
	position,
	created_at
)
VALUES
	($1, $2, $3, $4, $5)
;

-- name: GetImagesForVariant :many
SELECT
	*
FROM
	products.variant_images
WHERE
	variant_uuid = $1
ORDER BY
	position ASC;

-- name: DeleteImagesForVariant :exec
DELETE FROM products.variant_images
WHERE
	variant_uuid = $1;

-- name: DeleteVariantImage :exec
DELETE FROM products.variant_images
WHERE
	image_uuid = $1;
