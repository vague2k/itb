set quiet

default:
    just --list

dev:
    rm -f ./views/assets/css/output.css
    just --dotenv-filename=.env watch-and-generate

[parallel]
watch-and-generate: tailwind-watch templ

tailwind-watch:
    tailwindcss -i ./views/assets/css/globals.css -o ./views/assets/css/output.css --watch

templ:
    go tool templ generate --watch --proxyport="{{ env('PORT', '') }}" --cmd="go run ./cmd/server" --open-browser=false

# Build templ, Tailwind, and the server binary
build:
    go tool templ generate
    tailwindcss -i ./views/assets/css/globals.css -o ./views/assets/css/output.css --minify
    go build -o ./tmp/server ./cmd/server

# Migrate the database
migrate command="up":
    go run ./cmd/server migrate {{ command }}

# Run all tests
test:
    go test ./...
