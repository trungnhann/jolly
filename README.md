# jolly

Base scaffold for a Go e-commerce backend following a modular DDD-inspired structure.

## Current Modules

- `orders`
- `payments`
- `inventory`
- `users`
- `products` (Product catalog with Categories, Brands, and Variants support)

## Run

Recommended local setup (Docker-first):
```bash
task up
```

Run the Go process directly on the host (requires a DB connection string):
```bash
export POSTGRES_URL='postgres://jolly:jolly@localhost:5433/jolly?sslmode=disable'
task run
```

## Test

```bash
task test
```

## Example Requests

Health check:

```bash
curl -i http://localhost:8080/health
```

Create a user:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "name": "Customer One",
    "role": "customer"
  }'
```

Create an order (use an existing `user_uuid` as `customer_id`):

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "00000000-0000-0000-0000-000000000000",
    "currency": "USD",
    "items": [
      {
        "sku": "SKU-RED-SHOE-42",
        "quantity": 1,
        "unit_price_cents": 12999
      }
    ]
  }'
```

## Documentation

- `docs/architecture-reference.md`
- `docs/external-events.md`
- `docs/operations.md`

## Next Steps

- add persistence + migrations for `payments`
- add HTTP contracts for modules without public APIs yet
- add external events / async workflows (outbox + broker)
- add more modules such as `customers` and `carts`
