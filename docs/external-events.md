# External Events

## Purpose

This document defines domain and integration events emitted across modules.

Current status:

- there is no event bus or message broker integration yet
- modules communicate synchronously via module contracts (`backend/common/module/contracts`)

This doc is kept so we can add events later without changing the documentation structure.

Events should be:

- explicit
- versioned
- idempotent to consume
- documented before broad adoption

## Event Envelope

Suggested envelope:

```json
{
  "event_name": "orders.order_placed",
  "event_version": 1,
  "event_id": "uuid",
  "occurred_at": "2026-06-11T10:00:00Z",
  "producer": "orders",
  "payload": {}
}
```

## Initial Event Catalog

### orders.order_placed (planned)

Produced by: `orders`

Used by:

- payments
- inventory
- fulfillment
- notifications

Payload:

- order_uuid
- customer_uuid
- total_amount
- currency
- ordered_at
- line_items_snapshot

### payments.payment_authorized (planned)

Produced by: `payments`

Used by:

- orders

Payload:

- order_uuid
- payment_uuid
- authorized_amount
- currency
- authorized_at

### payments.payment_captured (planned)

Produced by: `payments`

Used by:

- orders
- accounting

Payload:

- order_uuid
- payment_uuid
- captured_amount
- currency
- captured_at

### inventory.stock_reserved (planned)

Produced by: `inventory`

Used by:

- orders
- fulfillment

Payload:

- order_uuid
- reservation_uuid
- reserved_items

### inventory.stock_reservation_failed (planned)

Produced by: `inventory`

Used by:

- orders

Payload:

- order_uuid
- reason
- failed_items

### fulfillment.shipment_created (planned)

Produced by: `fulfillment`

Used by:

- orders
- notifications

Payload:

- order_uuid
- shipment_uuid
- created_at

### fulfillment.order_shipped (planned)

Produced by: `fulfillment`

Used by:

- orders
- notifications

Payload:

- order_uuid
- shipment_uuid
- shipped_at

### fulfillment.order_delivered (planned)

Produced by: `fulfillment`

Used by:

- orders
- notifications

Payload:

- order_uuid
- shipment_uuid
- delivered_at

## Event Design Rules

- Never publish internal DB models directly
- Use stable names and stable payload fields
- Add new fields backward-compatibly
- Prefer snapshots over live references for historical correctness
- Document retries and deduplication strategy per event
