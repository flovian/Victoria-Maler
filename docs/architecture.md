# Architecture

Victoria Maler is a Go web application with a server-rendered frontend and a
small JSON API. It stores data in SQLite and uses Bitcoin Core as an immutable
verification layer.

## Component overview

```
cmd/
  server/   HTTP server entrypoint (wiring + static assets)
  migrate/  standalone migration runner
  seed/     demo data seeder
internal/
  config/       env-aware configuration (.env + environment variables)
  database/     SQLite connection + SQL migration runner
  models/       domain structs matching the database schema
  repositories/ data access (SQLite implementations)
  services/     business logic (auth, campaigns, donations, reports, evidence)
  bitcoin/      OP_RETURN framing, bitcoind JSON-RPC, anchoring + verification
  handlers/     HTTP handlers and page renderer
  middleware/   request logging, JWT auth, admin guards
  routes/       route table
migrations/     versioned SQL schema files
web/
  templates/    layouts + Go html/template pages
  static/       CSS, JavaScript, images and uploads
scripts/        run / migrate / seed helpers
```

## Request flow

1. `main.go` loads config, opens SQLite, runs migrations, constructs
   repositories/services/bitcoin clients and builds the `HandlerSet`.
2. `routes.New` registers HTML pages and JSON API routes.
3. `middleware.Logger` logs every request. Protected routes are wrapped with
   `RequireAuth`, which validates a JWT from the `token` cookie or
   `Authorization: Bearer` header.
4. Handlers call services, which use repositories against SQLite.
5. Cleanup reports and evidence are SHA-256 hashed and anchored through
   `bitcoin.AnchoringService`, then recorded in `bitcoin_records`.

## Key design decisions

* **Server-rendered pages** — Go's `html/template` renders pages from a shared
  base layout; the renderer exposes helper funcs such as `progressPercent`.
* **JSON API envelope** — every response uses `{ success, message, data }`
  (`internal/utils/response.go`).
* **Repository interfaces + SQLite implementations** — services depend on
  interfaces so they can be tested against an in-memory database.
* **Graceful Bitcoin degradation** — if no bitcoind node is reachable, anchors
  are stored with a `sim-` txid and `simulated` status so the entire platform
  still works; verification clearly labels these records.

## Environment

See `.env` for the full set of keys: `DB_PATH`, `PORT`, `ENV`, `JWT_SECRET`,
`UPLOAD_DIR` and the `BITCOIN_*` settings.
