-- name: GetStockForUpdate :one
SELECT * FROM inventory.stocks WHERE sku = $1 FOR UPDATE;

-- name: UpdateStockQuantity :exec
UPDATE inventory.stocks SET quantity = $2, updated_at = NOW() WHERE stock_uuid = $1;

-- name: CreateReservation :exec
INSERT INTO inventory.reservations (reservation_uuid, order_uuid, created_at)
VALUES ($1, $2, $3);

-- name: CreateReservationItem :exec
INSERT INTO inventory.reservation_items (reservation_uuid, stock_uuid, quantity)
VALUES ($1, $2, $3);

-- name: GetReservationForOrder :one
SELECT * FROM inventory.reservations WHERE order_uuid = $1 LIMIT 1;

-- name: UpsertStock :one
INSERT INTO inventory.stocks (stock_uuid, sku, quantity, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT (sku) DO UPDATE
SET quantity = EXCLUDED.quantity, updated_at = NOW()
RETURNING *;
