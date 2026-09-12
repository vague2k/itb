# ITB (In The Budget) 

**ITB** is a single-user budget tracker. Make categories (Groceries, House maintenance), give each a starting amount, then add income and expenses as they happen. A category's balance is the sum of its transactions, so every number traces back to a real entry.

Self-hosted, but never reachable by the public. This is a private tool for one user, and that fact drives deliberate design choices throughout: no accounts, no auth, no multi-tenancy.

Domain glossary and architecture live in **CONTEXT.md**. AGENTS.md holds the agent rules.

## Tech stack

Go, chi, templ, sqlc with the modernc.org/sqlite driver (pure-Go, no cgo), goose migrations, htmx, Tailwind CLI, and shadcn-templ.

Money is stored as `int64` cents everywhere and only parsed or formatted through `internal/money`.

## Project structure

```
.
├── cmd/                # Application entrypoint
│   └── server/         # HTTP server binary (`serve` default; `migrate up|down|status`)
│
├── config/             # Config struct + Init() (sole env-var reader), opens DB
├── internal/
│   ├── money/          # Cent parsing and formatting (never float)
│   └── testdb/         # file-backed sqlite test DB helper
│
├── server/
│   ├── database/       # sqlc layer
│   │   ├── generated/  # sqlc-generated query code and params (build artifact)
│   │   ├── migrations/ # goose migrations
│   │   ├── models/     # sqlc-generated table models
│   │   └── queries/    # hand-written SQL for sqlc
│   ├── handlers/       # HTTP handlers + composition root
│   ├── services/       # Application logic used by handlers
│   ├── router.go       # Route definitions
│   └── server.go       # Server bootstrap, asset and component-script routes
│
└── views/              # Frontend (templ + static assets)
    ├── assets/         # CSS, JS (htmx, theme)
    ├── components/     # Hand-rolled compositions (category list, detail, theme switcher)
    ├── layouts/        # Page shells (Base)
    ├── pages/          # Pages (categories, category detail)
    └── shadcn/         # Vendored shadcn-templ primitives + script bundle
```

## Development

```sh
just dev      # Tailwind watch + templ generate --watch (restarts server on .templ change)
just build    # templ generate + minified Tailwind + go build -> ./tmp/server
just test     # go test ./...
just migrate up|down|status   # go run ./cmd/server migrate

go tool sqlc generate   # regenerate query code and models after editing queries/
```
