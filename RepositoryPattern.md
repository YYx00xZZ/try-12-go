# Repository Pattern Notes

The project already isolates persistence behind the `UserRepository` interface so that the rest of the code never depends on a specific database driver.

- **Contract** — `internal/repository/user.go` defines the `UserRepository` interface (`List(ctx context.Context) ([]User, error)`). Handlers and other callers accept that interface only.

- **Postgres implementation** — `internal/repository/postgres/user.go` provides a `UserRepository` backed by `*sql.DB`. All SQL and Postgres-specific details live here.

- **Composition root** — `cmd/server/main.go` wires the implementation by calling `postgresrepo.NewUserRepository(pg)` before handing it to `handler.NewUserHandler`. That single location controls which database implementation is used.

To switch databases:

1. Add another package (for example `internal/repository/sqlite`) that implements `UserRepository` with the alternative driver.
2. Provide or reuse a connection helper suited to that database.
3. Change the wiring in `cmd/server/main.go` (or in configuration/factory logic) to instantiate the new repository implementation.

Because the handler layer only references the interface, none of the HTTP code or higher-level logic needs to change when you swap implementations.

## Switching via configuration

- Set `DB_BACKEND=postgres` or `DB_BACKEND=mongo` in `.env` (or the runtime environment).
- Postgres uses the existing DSN variables (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DATABASE_URL`).
- MongoDB uses `MONGO_URI`, `MONGO_DB`, and optionally `MONGO_COLLECTION` to choose the target collection.
- The server reads `DB_BACKEND` at startup and instantiates the matching repository without requiring handler changes.

# Development
Seed mongodb:
```
docker compose exec mongo mongosh mydb --eval '
db.users.deleteMany({});
db.users.insertMany([
  {id: 1, name: "Alice"},
  {id: 2, name: "Bob"},
  {id: 3, name: "Charlie"}
]);
'
```
