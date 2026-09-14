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

# Build the server for an environment: development (default), staging, or production
build env="development":
    ./build.sh {{ env }}

# Migrate the database
migrate command="up":
    go run ./cmd/server migrate {{ command }}

# Deploy to an environment: staging or production
deploy env:
    ./deployments/deploy.sh {{ env }}

# Run all tests
test:
    go test ./...
