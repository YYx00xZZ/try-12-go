# try-12-go

REST API sample with switchable persistence backends (Postgres or MongoDB), containerised with Docker Compose and documented via Swagger UI.

## Prerequisites
- Docker and Docker Compose v2
- (Optional) `swag` CLI for regenerating API docs: `go run github.com/swaggo/swag/cmd/swag@v1.8.12`

## Running with Postgres (default)
1. Build the images so migrations and binaries are updated:
   ```bash
   docker compose build --no-cache migrate app
   ```
2. Start the stack:
   ```bash
   docker compose up
   ```
   Postgres becomes healthy, migrations apply automatically, then the API starts on [`http://localhost:8080`](http://localhost:8080).

## Switching to MongoDB
1. Set `DB_BACKEND=mongo` in `.env` (or export it in your shell).
2. Adjust the Mongo connection values if needed (`MONGO_URI`, `MONGO_DB`, `MONGO_COLLECTION`). The Compose file already points the app to the bundled `mongo` service via `mongodb://mongo:27017`.
3. Rebuild and start:
   ```bash
   docker compose build app
   docker compose up app mongo
   ```
   Compose still boots the Postgres container (to satisfy dependencies), but the migrate task exits immediately when `DB_BACKEND` is not `postgres`.
   Seed `users` documents in Mongo (fields `id` and `name`) to exercise the `/users` endpoint.

## API Documentation
- Swagger UI: [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)
- Specification artifacts live in `docs/` and are regenerated with:
  ```bash
  go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/server/main.go -o docs --parseInternal
  ```

## Configuration Reference
`.env` variables (loaded by Docker Compose and the app):

| Variable | Default | Description |
| --- | --- | --- |
| `DB_BACKEND` | `postgres` | Selects the persistence backend (`postgres` or `mongo`). |
| `DB_HOST` / `DB_PORT` / `DB_USER` / `DB_PASSWORD` / `DB_NAME` | see `.env` | Standard Postgres connection info. |
| `DATABASE_URL` | derived | Full Postgres DSN used by migrations and local runs. |
| `MONGO_URI` | `mongodb://localhost:27017` | Connection string for MongoDB. Overridden to `mongodb://mongo:27017` inside Compose. |
| `MONGO_DB` | `mydb` | Mongo database name housing the users collection. |
| `MONGO_COLLECTION` | `users` | Mongo collection queried by the API. |
| `MONGO_PORT` | `27017` | Published MongoDB port when using Docker Compose. |
| `PORT` | `8080` | HTTP port exposed by the API container. |

## Useful Commands
- Run migrations only (Postgres): `docker compose run --rm migrate`
- Follow app logs: `docker compose logs -f app`
- Tear down: `docker compose down -v`
