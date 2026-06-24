ALTER TABLE orders.orders
  ALTER COLUMN customer_id TYPE uuid USING customer_id::uuid;

ALTER TABLE orders.orders
  ADD CONSTRAINT orders_orders_customer_id_fkey
  FOREIGN KEY (customer_id) REFERENCES users.users(user_uuid);

CREATE INDEX idx_orders_customer_id
  ON orders.orders (customer_id);