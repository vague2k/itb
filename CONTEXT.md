# Domain glossary

itb is a self-hosted, private budget tracker for one user. That drives deliberate choices throughout: no auth, no multi-tenancy, no monthly periods. Simplicity over generality.

## Terms

- **Category**: a labeled fund with a running balance, e.g. Groceries or House maintenance. Has a name. Never resets.
- **Transaction**: one income or expense against a Category. Amount is a signed integer in cents, positive for income and negative for expense. Has an optional note and a timestamp.
- **Balance**: the sum of a Category's transaction amounts. Never stored, always computed with a sqlc join.
- **Starting amount**: the first transaction for a Category, entered when it is created. Not a special column.

## Architecture in one paragraph

Go + chi server. `config.Config` reads env and opens SQLite (sqlc-generated queries). `handlers.Handler` is the composition root: it builds the services and the sub-handlers. Handlers parse HTTP, delegate to services, and render templ views. Mutations return htmx fragments so the category list updates in place. shadcn-templ primitives live in `views/shadcn`; hand-rolled compositions live in `views/components`.

## Error convention

- `services.UserError`: message shown verbatim to the user (toast).
- `services.InternalError`: wraps the cause, gets logged, user sees a generic message.

## Money convention

All amounts are `int64` cents. Parsing and formatting go through `internal/money`, nowhere else.
