# Architecture Reference

## Overview

`jolly` is a Go backend built as a modular monolith.

The project uses DDD as a set of boundary and ownership rules, not as ceremony.
The goal is to keep business modules explicit, let the codebase grow without
collapsing into one giant service layer, and still stay practical for a small
team.

At the current stage, the system should be understood as:

- one deployable application
- one Go process
- one Postgres database
- multiple module-owned schemas where useful
- synchronous in-process module collaboration by default
- a possible path to events and async processing later, when the need is real

This is not a microservice architecture. Internal module contracts are in-process
application boundaries, not network APIs.

## High-Level Overview

- Entry point: `backend/cmd/main.go`
- Composition root: `backend/svc.go`
- Module registry and lifecycle: `backend/common/module/module.go`
- Contracts registry: `backend/common/module/contracts/contracts.go`
- Shared HTTP router and middleware: `backend/common/http`
- Shared platform primitives: `backend/common`
- Business modules:
  - `backend/orders`
  - `backend/users`
  - `backend/payments`
  - `backend/inventory`

## Current Bounded Contexts

### orders

Owns:

- order placement
- order persistence
- order line items
- order state transitions
- order lifecycle timestamps

Current aggregate candidate:

- `orders/domain.Order`

`orders` is the source of truth for order state.

### users

Owns:

- user identity
- user role information
- user creation and retrieval

`users` is the source of truth for customer existence in the current model.

### payments

Owns:

- payment authorization behavior
- payment-related module contract used by `orders`

At the moment, this module is an in-process capability and does not persist its
own state yet.

### inventory

Owns:

- stock reservation behavior
- inventory-related module contract used by `orders`

This module is persistent and has its own `inventory.stocks` DB schema.

### products

Owns:

- product catalog cataloging, brand mapping, and hierarchical categories
- variants list and pricing definitions
- public module contracts

## Repository Shape

```text
jolly/
├── backend/
│   ├── cmd/                 # executable entrypoint
│   ├── common/              # platform and narrow shared-kernel primitives
│   ├── inventory/           # bounded context
│   ├── orders/              # bounded context
│   ├── payments/            # bounded context
│   ├── users/               # bounded context
│   └── svc.go               # composition root
├── docs/
├── docker-compose.yaml
└── Taskfile.yml
```

## Module Shape

Use this internal shape when the module needs it:

```text
<module>/
├── adapters/
│   ├── db/
│   ├── external/
│   ├── queue/
│   └── storage/
├── api/
│   ├── http/
│   └── module/
├── app/
│   ├── command/
│   └── query/
├── domain/
└── module.go
```

Notes:

- not every module needs every folder on day one
- add structure when complexity justifies it
- `queue/` is only needed when async processing is actually introduced
- `api/module` is for in-process contracts, not public HTTP transport

## Runtime Wiring

The application is composed in `backend.New(...)` in `backend/svc.go`.

Current boot sequence:

1. create the shared Postgres pool in `backend/cmd/main.go`
2. create the shared contracts registry
3. create the module registry
4. register modules in explicit startup order:
   - `payments`
   - `inventory`
   - `users`
   - `orders`
5. create the shared Echo router
6. run module lifecycle methods:
   - `Init(ctx)`
   - `RegisterContracts(ctx, contracts)`
   - `RegisterHttp(ctx, e)`
7. verify all contracts are present
8. start the HTTP server

The current init order matters. For example, `users` is initialized before
`orders`, which is important because `orders` can depend on `users` schema and
foreign keys.

## Module Lifecycle

Every module implements `common/module.Module`:

- `Name()`
- `Init(ctx)`
- `RegisterContracts(ctx, contracts)`
- `RegisterHttp(ctx, e)`

Responsibilities:

- `Init`:
  - build repositories and services
  - run migrations when the module owns DB schema
  - wire internal handlers
- `RegisterContracts`:
  - publish the module's in-process contract to other modules
- `RegisterHttp`:
  - register HTTP endpoints into the shared router when needed

This keeps each bounded context responsible for its own wiring while the
composition root controls startup order.

## Layer Responsibilities

### domain

Contains:

- entities and aggregate roots
- value objects
- invariants
- state transitions
- repository interfaces when the domain needs persistence abstraction
- rehydration helpers for persisted aggregates

Must not depend on:

- HTTP
- Echo
- sqlc
- pgx
- another module's adapters

The domain is where business truth lives.

### app

Contains:

- command handlers
- query handlers
- use-case orchestration
- coordination across repositories and module contracts

May depend on:

- its own domain
- repository interfaces
- small dependency interfaces
- internal module contracts

The app layer coordinates the flow. It should not become the place where all
domain rules silently move.

### adapters

Contains:

- Postgres repositories
- sqlc-generated DB access
- external service clients
- storage implementations
- future event publishers and subscribers

Adapters translate infrastructure concerns into module-friendly interfaces and
models.

### api

Contains:

- HTTP transport
- request and response DTO mapping
- edge validation
- internal contract DTOs under `api/module`

The API layer should not contain core business decisions.

## DDD Guidance

### Aggregate Ownership

The main aggregate in the current codebase is `orders.Order`.

It should own:

- customer reference
- line items
- total amount
- status transitions
- `placed_at_utc`
- `created_at`
- `updated_at`

As the order flow grows, lifecycle rules should stay on the aggregate when they
represent business truth, not leak into app handlers as string assignments.

### Application Orchestration

The app layer should:

- validate and map incoming commands
- load or construct aggregates
- call repositories
- call cross-module contracts
- persist state changes

The app layer should not invent business rules that belong to the aggregate.

### Current Order Placement Rule

The most important architectural rule in the current codebase is:

1. validate the request
2. construct the `Order` aggregate
3. persist the order first
4. call downstream capabilities such as `inventory` and `payments`
5. move the order through explicit states after each step

Current early states are:

- `pending`
- `inventory_reserved`
- `payment_authorized`
- `failed`

Why this matters:

- the order exists before risky side effects happen
- failure leaves a traceable record in the source of truth
- compensation becomes possible later
- observability improves immediately

Asynchronous decoupling is implemented using RabbitMQ and Watermill, allowing us to publish cross-module events (e.g., `orders.order_paid` event triggers stock level decrements) asynchronously.

## Module Communication

Preferred order of collaboration:

1. same-module collaboration through app + domain
2. cross-module synchronous collaboration through internal contracts
3. cross-module asynchronous collaboration through events

Rules:

- do not import another module's adapters
- do not import another module's internal domain for convenience
- keep contract DTOs explicit
- treat `api/module` as an internal app boundary
- do not pretend internal contracts are public APIs

### Contracts Registry

The shared contracts bundle lives in `backend/common/module/contracts/contracts.go`.

Current contracts:

- `Orders`
- `Payments`
- `Inventory`
- `Users`

Modules publish themselves into this registry during `RegisterContracts`, and
startup fails fast if any required contract is missing.

## HTTP Layer

- Router: Echo
- Shared router created once in `backend/common/http.NewEcho()`
- Shared error model in `backend/common/errors.go`
- Shared Echo error handler in `backend/common/errors_echo.go`
- Shared middleware in `backend/common/http/middlewares.go`

Current shared route:

- `GET /health`

Module HTTP APIs are registered into the same router during startup.

### Request Lifecycle

Typical request flow:

1. request hits the shared Echo router
2. middleware runs:
   - timeout
   - panic recovery
   - correlation ID propagation
   - request/response logging
3. module HTTP handler decodes the request
4. HTTP DTOs are mapped into app commands or queries
5. app layer executes the use case
6. domain and repositories enforce business and persistence rules
7. errors bubble up to the shared Echo error handler
8. a stable JSON error response is returned

### Error Model

Public errors should be explicit and stable:

- `message`: safe public message
- `slug`: machine-readable identifier
- `details`: optional field or entity details

Generic unexpected errors become `500 internal_server_error`. Business or input
errors should be mapped before they reach that fallback.

## Persistence and Migrations

Current persistence model:

- one Postgres database
- module-owned schemas where useful
- module-owned migrations
- module-owned `schema_migrations` table per schema

Current DB-backed modules:

- `users`
- `orders`

Current in-memory or non-persistent modules:

- `payments`

Migrations are embedded into the binary per module and executed during module
initialization.

### Cross-Schema Foreign Keys

Cross-schema foreign keys are acceptable in the current modular monolith model
when both modules:

- live in the same database
- deploy together
- share the same lifecycle expectations

The current `orders -> users` relationship fits this model.

If a module is likely to move to a separate database later, avoid introducing
hard foreign-key coupling too early.

## `common` Rules

`backend/common` is allowed, but it must stay narrow.

Good candidates:

- UUID helpers
- time and clock helpers
- generic DB transaction helpers
- generic HTTP middleware
- logging utilities
- migration bootstrap
- truly shared low-level types such as `Currency`

Bad candidates:

- module-specific business rules
- module-specific enums
- convenience helpers that bypass module boundaries
- cross-module shortcuts that leak ownership

If a type belongs to one bounded context, keep it in that module.

## Strong Types

Use strong typing where it reduces boundary mistakes:

- HTTP schemas can expose explicit value shapes
- sqlc can map columns to domain-meaningful Go types
- module-owned IDs should use module-owned wrappers where that adds clarity

The goal is compile-time safety, not type proliferation.

## Testing Guidance

Current repo tasks:

- `task test`
- `task gen`
- `task fmt`
- `task tidy`
- `task check`

Testing intent:

- unit tests for domain and app logic
- adapter tests for DB behavior
- higher-level flow checks when orchestration spans multiple modules

When behavior depends on persistence or migrations, tests should verify the real
repository behavior, not just mocks.

## Local Development

Common local tasks:

```bash
task up
task down
task restart
task logs
task run
task test
task gen
```

Migration helper:

```bash
task mig:new MODULE=orders NAME=add_customer_fk
```

The local Docker setup runs the app and its dependencies for development. The
current development model favors simple local reproducibility over production
parity for every infrastructure concern.

## Asynchronous Messaging & Message Broker

We use **Watermill** with **RabbitMQ** (AMQP) for cross-module asynchronous collaboration. This enables decoupled event-driven workflows where transactional safety or immediate consistency is not strictly required.

### Broker & Connection
- **Broker**: RabbitMQ running locally on port `5672` (or configured via `RABBITMQ_URL`).
- **Library**: `ThreeDotsLabs/watermill-amqp/v3` is utilized to implement standard Publisher/Subscriber mechanics.
- **Durable Queues**: All AMQP queues are created as durable, ensuring that messages are not lost during broker restarts.

### Asynchronous Flow Pattern
Modules implementing asynchronous events must declare a `queue/` adapter folder (e.g., `backend/inventory/adapters/queue/`) housing:
- **Publishers**: Components that emit domain/integration events (e.g., publishing `orders.order_paid` when payment is confirmed).
- **Consumers/Subscribers**: Consumers that handle incoming messages asynchronously (e.g., decrementing stocks in inventory upon receiving `orders.order_paid` event).

### Transactional Decoupling
Consistent with our database transaction guidelines, side effects (such as publishing to the message broker or sending email notifications) must be performed **after the database transaction block successfully commits**, avoiding duplicate publication in case of retry or rollback.

## Architectural Rules

- keep bounded context ownership explicit
- keep `common` small and boring
- persist important business state before risky side effects
- let aggregates own business invariants and state transitions
- use internal contracts for sync collaboration
- introduce events only when async decoupling is actually needed
- prefer simple maintainable flows before adding outbox or saga machinery
- update docs when architecture decisions change
