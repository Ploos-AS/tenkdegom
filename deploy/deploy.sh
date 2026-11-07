#!/usr/bin/env bash
set -euo pipefail

: "${DEPLOY_HOST:?Set DEPLOY_HOST}"
: "${DEPLOY_USER:?Set DEPLOY_USER}"
: "${DEPLOY_PATH:?Set DEPLOY_PATH}"
DEPLOY_PORT="${DEPLOY_PORT:-22}"

[[ "${DEPLOY_PATH}" == "/srv/tenkdegom" ]] || { echo "Unexpected DEPLOY_PATH" >&2; exit 1; }
[[ "${DEPLOY_USER}" == "tenkdegom-deploy" ]] || { echo "Unexpected DEPLOY_USER" >&2; exit 1; }
[[ "${DEPLOY_HOST}" =~ ^[a-zA-Z0-9.:-]+$ ]] || { echo "Invalid DEPLOY_HOST" >&2; exit 1; }
[[ "${DEPLOY_PORT}" =~ ^[0-9]+$ ]] || { echo "Invalid DEPLOY_PORT" >&2; exit 1; }
[[ -f dist/index.html ]] || { echo "Build first: npm run build" >&2; exit 1; }

release="$(date -u +%Y%m%dT%H%M%SZ)-${GITHUB_SHA:-manual}"
release="${release:0:64}"
[[ "$release" =~ ^[a-zA-Z0-9_-]+$ ]] || { echo "Invalid release identifier" >&2; exit 1; }

ssh -p "$DEPLOY_PORT" -- "$DEPLOY_USER@$DEPLOY_HOST" "mkdir -p '$DEPLOY_PATH/releases/$release'"
rsync -az --delete -e "ssh -p $DEPLOY_PORT" dist/ "$DEPLOY_USER@$DEPLOY_HOST:$DEPLOY_PATH/releases/$release/"
ssh -p "$DEPLOY_PORT" -- "$DEPLOY_USER@$DEPLOY_HOST" "test -f '$DEPLOY_PATH/releases/$release/index.html' && ln -sfn '$DEPLOY_PATH/releases/$release' '$DEPLOY_PATH/current.new' && mv -Tf '$DEPLOY_PATH/current.new' '$DEPLOY_PATH/current'"
echo "Published release: $release"
