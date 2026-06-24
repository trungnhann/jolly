# jolly

Base scaffold for a Go e-commerce backend following a modular DDD-inspired structure.

## Current Modules

- `orders`
- `payments`
- `inventory`

## Run

```bash
go run ./backend/cmd
```

## Test

```bash
go test ./...
```

## Example Request

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cus_123",
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

- add persistence adapters
- add sql migrations
- add HTTP contracts per module
- split `orders` into command/query if needed
- add more modules such as `catalog`, `customers`, and `carts`
