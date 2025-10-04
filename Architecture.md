# Architecture Overview

```mermaid
graph TD
  subgraph Client
    A[HTTP Client]
    B[Swagger UI]
  end

  subgraph Container Stack
    direction TB
    A -->|REST calls| S(Server)
    B -->|Docs via /docs| S

    subgraph GoApp
      S[Echo-based API]
      H[Handlers]
      R[Repository Interfaces]
    end

    S --> H
    H --> R

    subgraph Persistence
      direction LR
      RP[(Postgres Repo)]
      RM[(Mongo Repo)]
    end

    R -->|DB_BACKEND=postgres| RP
    R -->|DB_BACKEND=mongo| RM

    subgraph Databases
      P[(Postgres)]
      M[(MongoDB)]
    end

    RP --> P
    RM --> M
  end

  subgraph Tooling
    MIG[migrate binary]
    DOCS[Swag CLI]
  end

  MIG -->|migrations/| P
  DOCS -->|generates| DOCSFILES[docs/]
  DOCSFILES --> B
```

- The API runs inside Docker Compose alongside Postgres, MongoDB, and a one-shot migration container.
- Handlers depend on repository interfaces, allowing runtime selection between Postgres and Mongo implementations via `DB_BACKEND`.
- Swagger UI serves generated documentation while clients interact through the same Echo server.
- `migrate` applies SQL migrations only when Postgres is selected; Mongo uses ad-hoc seeding.
