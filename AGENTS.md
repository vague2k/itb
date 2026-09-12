# AGENTS.md

## Style Guide

### General Principles
- Aggressively follow idomatic Go practices and conventions.
- Keep logic in one function unless it is composable or reusable.
- Do not extract single-use helpers preemptively. Inline at the call site unless the helper is reused, hides a genuinely complex boundary, or has a clear independent name that improves the caller.
- Keep helpers in the same package they support, in a dedicated utils file.
- Do not over-abstract simple expressions into many single-use helpers; extract only when it names a real concept like requireConfig or readMetadata.
- Prefer applying these rules until they fall out of gofmt or lint expectations.

Reduce total variable count by inlining when a value is only used once.

```go
// Good
output, err := services.Parse(filepath.Join(dir, "input.txt"))

// Bad
inputFilePath := filepath.Join(dir, "input.txt")
output, err := services.Parse(inputFilePath)
```

### File Structure

General File structure should follow this hierarchy from the top down.

- imports
- constants/vars
- neccessary structs for methods declared in the file
- main file struct
- main struct constructor
- public struct functions
- private struct functions
- private file util functions

```go
// Good
package store

import "strings"

const max = 3

type item struct { name string }

type Store struct { items []item }

func NewStore() *Store { return &Store{} }

func (s *Store) Get(name string) item { return s.find(name) }

func (s *Store) find(name string) item {
	for _, it := range s.items {
		if it.name == clean(name) {
			return it
		}
	}
	return item{}
}

func clean(s string) string { return strings.ToLower(s) }
```

### Destructuring

Avoid unnecessary destructuring. Use dot notation to preserve context.

```go
// Good
struct.a
struct.b

// Bad
a, b := struct.a, struct.b
```

### Functions

Inline small one liner returns and struct getters

```go
// Good
func (s *MyStruct) Value() { return s.value }

// Bad
func (s *MyStruct) Value() { 
    return s.value
}
```

### Error handling

- Never use `_` to discard errors. Always check errors.

Inline `error` checking if a function only returns an error.

```go
// Good
if err := hasError(); err != nil {
    return err
}

// Bad
err := hasError()
if err != nil {
    return err
}
```

Returned `error` should be wrapped to preserve context.

```go
// Good
if err := parse(); err != nil {
    return fmt.ErrorF("could not parse thing: %w", err)
}

// Bad
if err := parse(); err != nil {
    return err
}
```

### Control Flow

Avoid else statements. Prefer early returns.

```go
// Good
func doThing(condition bool) int {
    if condition { return 1 }
    return 2
}

// Bad
func doThing(condition bool) int {
    if condition { 
        return 1
    } else {
        return 2
    }
}
```

### Complex logic

- Keep helpers in the same package they support, in a dedicated utils file.
- Do not over-abstract simple expressions into many single-use helpers; extract only when it names a real concept like requireConfig or readMetadata.
- Do not add comments randomly. Seriously consider comments for non-obvious constraints and surprising behavior, not for obvious assignments or control flow.

When a function has several validation branches or supporting details, make the main function read as the happy path and move supporting details into small helpers below it.

```go
// Good
func loadThing(input any) {
  config := requireConfig(input)
  metadata := readMetadata(input)
  return createThing(config, metadata)
}

func requireConfig(input any) {
  ...
}
```

## Testing

- Test actual implementation, do not duplicate logic into tests
- Critical paths get a "happy path" and a "negative" case for resilience.
- Handler tests cover only endpoints that mutate data. Do not test handlers that only render HTML.
- Services are tested against a freshly migrated SQLite database from `internal/testdb`; no mocks and no network calls.
- Money is tested directly in `internal/money` with round-trip and invalid-input cases.

## Commands

```sh
just dev              # tailwind watch + templ watch with hot reload
just build            # templ generate + tailwind + go build
just migrate up       # apply migrations (down, status also valid)
just test             # go test ./...
go tool sqlc generate # regenerate query code and models
```
