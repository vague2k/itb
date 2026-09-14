#!/usr/bin/env bash
set -euo pipefail

env="${1:-development}"

case "$env" in
  development)
    name="development"
    dotenv=".env"
    ;;
  prod | production)
    name="production"
    dotenv=".env.production"
    ;;
  staging)
    name="staging"
    dotenv=".env.staging"
    ;;
  *)
    echo "error: unknown environment '$env' (expected development, staging, or production)" >&2
    exit 1
    ;;
esac

if [[ ! -f "$dotenv" ]]; then
  echo "error: env file '$dotenv' does not exist (create it first)" >&2
  exit 1
fi

set -a
# shellcheck disable=SC1090
source "$dotenv"
set +a

GOOS= GOARCH= go tool templ generate
tailwindcss -i ./views/assets/css/globals.css -o ./views/assets/css/output.css --minify

out="dist/server-$name"
if [[ "${GOOS:-}" == "windows" ]]; then
  out="$out.exe"
fi

GOOS="${GOOS:-}" GOARCH="${GOARCH:-}" go build -ldflags "-X itb.ihatedoing.work/config.buildEnv=$name" -o "$out" ./cmd/server
echo "built $out (env=$name, embedded buildEnv=$name)"
