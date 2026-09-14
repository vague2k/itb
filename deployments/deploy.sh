#!/usr/bin/env bash
set -euo pipefail

# Usage:
#   ./deployments/deploy.sh production -> (:12942) itb
#   ./deployments/deploy.sh staging    -> (:12943) itb-staging
#
# Overridable via env:
#   SSH_TARGET   (default root@itb)
#   PROD_PATH    (default /root/itb)
#   STAGING_PATH (default /root/itb-staging)

cd "$(dirname "$0")/.."

SSH_TARGET="${SSH_TARGET:-root@itb}"
PROD_PATH="${PROD_PATH:-/root/itb}"
STAGING_PATH="${STAGING_PATH:-/root/itb-staging}"

MODE="${1:-}"
case "$MODE" in
    production)
        ENV_NAME="production"
        DOTENV=".env.production"
        DEPLOY_PATH="$PROD_PATH"
        SERVICE="itb"
        ;;
    staging)
        ENV_NAME="staging"
        DOTENV=".env.staging"
        DEPLOY_PATH="$STAGING_PATH"
        SERVICE="itb-staging"
        ;;
    *)
        echo "Usage: $0 [production|staging]" >&2
        exit 1
        ;;
esac

if [[ ! -f "$DOTENV" ]]; then
    echo "error: env file '$DOTENV' does not exist (create it first)" >&2
    exit 1
fi

echo "==> Building ($ENV_NAME)"
./build.sh "$ENV_NAME"

echo "==> Syncing binary and env to $SSH_TARGET:$DEPLOY_PATH"
ssh "$SSH_TARGET" "mkdir -p '$DEPLOY_PATH'"
rsync -a --info=progress2 "dist/server-$ENV_NAME" "$SSH_TARGET:$DEPLOY_PATH/server"
rsync -a --info=progress2 "$DOTENV" "$SSH_TARGET:$DEPLOY_PATH/"

# Run migrations before restarting so a failed migration aborts the deploy
# (set -e) and never restarts the service into a half-migrated state.
echo "==> Running migrations"
ssh "$SSH_TARGET" "cd '$DEPLOY_PATH' && ./server migrate up"

echo "==> Restarting $SERVICE"
ssh "$SSH_TARGET" "systemctl restart '$SERVICE'"

echo "==> Status"
ssh "$SSH_TARGET" "systemctl --no-pager -l status '$SERVICE' | head -20"
