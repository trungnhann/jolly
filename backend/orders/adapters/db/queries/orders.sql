-- name: SaveOrder :exec
INSERT INTO orders.orders (
	order_uuid,
	customer_id,
	currency,
	status,
	total_cents,
	placed_at_utc,
	created_at,
	updated_at
)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8)
;

-- name: SaveOrderLineItem :exec
INSERT INTO orders.order_line_items (
	line_item_uuid,
	order_uuid,
	sku,
	quantity,
	unit_price_cents
)
VALUES
	($1, $2, $3, $4, $5)
;

-- name: UpdateOrderStatus :exec
UPDATE orders.orders
SET
	status = $2,
	updated_at = $3
WHERE
	order_uuid = $1
;

-- name: GetOrder :one
SELECT
	*
FROM
	orders.orders
WHERE
	order_uuid = $1
LIMIT 1;

-- name: GetOrderLineItems :many
SELECT
	*
FROM
	orders.order_line_items
WHERE
	order_uuid = $1
ORDER BY
	line_item_uuid;
