# TODO

Each item is a complete slice. It carries the migration, the sqlc query, the service, the handler, the view, and the tests it needs, so you can pick one and land it without touching another first. Where two items touch the same table, they add their columns in separate migrations. Apply the migrations in order.

This list stays short on purpose. No reports dashboard, no import, no recurring entries. Those wait until the basics feel boring.

- [ ] Global quick-add
  - A button on every page opens a modal with a category combobox, an amount, and an income/expense toggle.
  - It posts to the existing transaction create endpoint.
  - Quick-add never asks for a date. The entry is stamped when you submit and You can fix the date later from the edit form.
  - The last used category becomes the default.
  - No migration.
  - Tests: the create endpoint mutates data, so add a handler test that posts from the modal and checks the updated fragment comes back.

- [ ] Rename a category
  - Add an update query, a service method that trims the name and rejects a blank one, a handler, and an inline rename control on the category page.
  - No migration.
  - Tests: service test for a normal rename and for a blank name.

- [ ] Edit and delete transactions
  - Add update and delete queries, service methods, handlers, and edit/delete controls on each ledger row.
  - You can change the amount, the note, and whether the row is income or expense.
  - Delete is permanent and asks for confirmation.
  - No migration. Income and expense stay derived from the sign of the amount.
  - Tests: service tests for updating each field, deleting a row, and rejecting a zero amount.

- [ ] Set and filter by transaction date
  - Add `occurred_at` to transactions and keep `created_at` as the record of when you typed the entry.
  - The create and edit forms get a date picker you can leave empty, which means now.
  - The category page gets a date-range picker that narrows the ledger and shows the range income, expense, and net next to the all-time balance.
  - Skip adjustment rows in those totals if the kind column from the adjustments item is already there.
  - Migration: add `occurred_at`, backfill it from `created_at`.
  - Queries: list a category's transactions filtered by range, plus a totals query. Regenerate sqlc.
  - Tests: service tests for the default timestamp, an explicit date, a range filter, and the totals.

- [ ] Balance adjustments
  - Add a `kind` column that marks each row as income, expense, or adjustment.
  - A dedicated "Adjust balance" action asks for the balance you actually see and writes the difference as an adjustment row.
  - Adjustments move the balance and are left out of the period income and expense totals.
  - Migration: add `kind`, backfill income or expense from the sign.
  - Queries: create takes kind; period totals skip adjustment rows.
  - Tests: service tests that an adjustment changes the balance but not the period income or expense, and that the action computes the right delta.

- [ ] Archive categories
  - Add `archived` to categories. Archive replaces delete.
  - An archived category leaves the main list for a dimmed archived section, can be restored, and rejects new transactions.
  - The quick-add picker and every other category list skip archived rows.
  - Migration: add `archived`, default 0.
  - Queries: archive, unarchive, list active, list archived. The transaction create service rejects an archived category.
  - Tests: service tests for archive, unarchive, and a write against an archived category.
  - This item retires the delete button and the `DeleteCategory` path.

- [ ] Export CSV and JSON
  - Two download endpoints and two buttons.
  - CSV gives one row per transaction with the category name.
  - JSON keeps categories, transactions, dates, kinds, and archive state.
  - No migration.
  - Tests: service tests for both serializations, including a transaction with a custom date and an adjustment.
