ALTER TABLE orders.orders
  ADD COLUMN created_at TIMESTAMPTZ;

ALTER TABLE orders.orders
  ADD COLUMN updated_at TIMESTAMPTZ;

UPDATE orders.orders
SET
  created_at = placed_at_utc,
  updated_at = placed_at_utc
WHERE created_at IS NULL OR updated_at IS NULL;

ALTER TABLE orders.orders
  ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE orders.orders
  ALTER COLUMN updated_at SET NOT NULL;
