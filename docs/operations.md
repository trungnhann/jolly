# Operations

## Overview

This document describes how to run and operate `jolly` in local development.

Current local setup is Docker-first:

- app runs in Docker
- Postgres runs in Docker with a named volume
- Redis runs in Docker
- MinIO runs in Docker
- module migrations run automatically on app startup

## Local Stack

Services started by `docker-compose.yaml`:

- `app`
- `postgres`
- `redis`
- `minio`
- `minio-init`

Host ports:

- app: `http://localhost:8080`
- Postgres: `localhost:5433`
- Redis: `localhost:6380`
- MinIO API: `http://localhost:9000`
- MinIO Console: `http://localhost:9001`

The non-default Postgres and Redis host ports are intentional to avoid
conflicts with services installed directly on the host machine.

## Common Tasks

Current `Taskfile.yml` tasks:

- `task up`: build and start the Docker stack
- `task down`: stop the Docker stack
- `task restart`: restart the Docker stack
- `task logs`: tail app logs
- `task ps`: show running containers
- `task run`: run the backend directly on the host
- `task test`: run all Go tests
- `task gen`: regenerate generated code
- `task fmt`: run `gofmt`
- `task tidy`: run `go mod tidy`
- `task check`: run `fmt`, `tidy`, and `test`

Migration helper tasks:

- `task mig:new MODULE=orders NAME=add_customer_fk`
- `task mig:new-orders NAME=add_customer_fk`
- `task mig:new-users NAME=add_created_at`

## Local Development

### Start the stack

Recommended command:

```bash
task up
```

Useful follow-up commands:

```bash
task ps
task logs
```

### Stop the stack

```bash
task down
```

### Restart after code or config changes

```bash
task restart
```

## Running the App Outside Docker

If you want to run the Go process directly on the host, the app still needs a
database connection string.

Required environment variable:

- `POSTGRES_URL` or `DATABASE_URL`

Example for the current local Docker Postgres:

```bash
export POSTGRES_URL='postgres://jolly:jolly@localhost:5433/jolly?sslmode=disable'
task run
```

The app listens on:

```text
:8080
```

## Health Checks

Basic liveness endpoint:

```bash
curl http://localhost:8080/health
```

Expected behavior:

- HTTP `200 OK`

## Database and Persistence

Current local Postgres configuration:

- user: `jolly`
- password: `jolly`
- database: `jolly`
- host port: `5433`

Postgres data is persisted in the Docker named volume:

- `postgres-data`

MinIO data is persisted in:

- `minio-data`

This means:

- `task down` does not remove data
- `docker compose down -v` removes volumes and deletes local persisted data

### Reset all local data

```bash
docker compose down -v
task up
```

### Reset only Postgres data

```bash
docker compose down
docker volume rm jolly_postgres-data
task up
```

The actual volume name may vary by compose project name. Check with:

```bash
docker volume ls | grep jolly
```

## Redis

Current local Redis configuration:

- host: `localhost`
- port: `6380`

Inside the Docker network, the app talks to Redis using:

```text
redis:6379
```

## MinIO

Current local MinIO configuration:

- endpoint: `http://localhost:9000`
- console: `http://localhost:9001`
- access key: `jolly`
- secret key: `jolly123`
- default bucket created by `minio-init`: `jolly`

This local setup is convenience-focused and is not intended to mirror production
security settings.

## Migrations

Modules that own DB schema embed their migrations and run them during startup.

Migration bootstrap behavior:

- schema is created if missing
- migration state is stored in `<schema>.schema_migrations`
- migrations are loaded from embedded files, not from the live filesystem

Current DB-backed modules with migrations:

- `users`
- `orders`

Current migration directories:

- `backend/users/adapters/db/migrations/`
- `backend/orders/adapters/db/migrations/`

### Creating a new migration

Recommended pattern:

```bash
task mig:new MODULE=orders NAME=add_status_index
```

This creates a timestamp-based `.up.sql` file such as:

```text
backend/orders/adapters/db/migrations/20260624123456_add_status_index.up.sql
```

### Important migration note

Because migrations are embedded into the binary, after adding or renaming a
migration file you must rebuild or restart the app:

```bash
task restart
```

## Generated Code

Generated code currently includes:

- OpenAPI server and client code
- sqlc database code

Regenerate after changing SQL or OpenAPI definitions:

```bash
task gen
```

## Troubleshooting

### App cannot connect to Postgres

Check:

- Docker stack is up: `task ps`
- Postgres is healthy in compose output
- the DSN uses port `5433` on the host, not `5432`

### Migration file is not applied

Check:

- file name uses the expected format: `<timestamp>_<name>.up.sql`
- the new file is in the correct module migration directory
- the app was restarted after the file was added

### Database name mismatch

If Postgres was initialized with an older database name in an existing volume,
changing `POSTGRES_DB` alone will not recreate the database. In that case:

- create the database manually, or
- reset the Postgres volume

### Port conflicts

Current local ports were chosen to avoid common conflicts:

- Postgres uses `5433`
- Redis uses `6380`

If these are already in use on your machine, change the host-side port mapping
in `docker-compose.yaml`.

## Reliability Notes

- keep write flows idempotent where possible
- log and propagate correlation IDs
- persist important business state before risky side effects
- prefer simple synchronous flows first
- add compensation, outbox, or async processing only when complexity justifies it
