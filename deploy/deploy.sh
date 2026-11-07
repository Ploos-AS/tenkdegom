#!/usr/bin/env bash
set -euo pipefail
# Manual deployment after provisioning a restricted SSH deploy account.
# Required: DEPLOY_HOST, DEPLOY_USER, DEPLOY_PATH. Optional: DEPLOY_PORT.
: "${DEPLOY_HOST:?Set DEPLOY_HOST}"
: "${DEPLOY_USER:?Set DEPLOY_USER}"
: "${DEPLOY_PATH:?Set DEPLOY_PATH}"
DEPLOY_PORT="${DEPLOY_PORT:-22}"
[[ "${DEPLOY_PATH}" == /* && "${DEPLOY_PATH}" != "/" ]] || { echo "DEPLOY_PATH must be an absolute non-root path" >&2; exit 1; }
[[ -f dist/index.html ]] || { echo "Build first: npm run build" >&2; exit 1; }
release="$(date -u +%Y%m%dT%H%M%SZ)-${GITHUB_SHA:-manual}"
release="${release:0:64}"
ssh -p "${DEPLOY_PORT}" -- "${DEPLOY_USER}@${DEPLOY_HOST}" "mkdir -p '${DEPLOY_PATH}/releases/${release}'"
rsync -az --delete -e "ssh -p ${DEPLOY_PORT}" dist/ "${DEPLOY_USER}@${DEPLOY_HOST}:${DEPLOY_PATH}/releases/${release}/"
ssh -p "${DEPLOY_PORT}" -- "${DEPLOY_USER}@${DEPLOY_HOST}" "ln -sfn '${DEPLOY_PATH}/releases/${release}' '${DEPLOY_PATH}/current.new' && mv -Tf '${DEPLOY_PATH}/current.new' '${DEPLOY_PATH}/current'"
echo "Published release: ${release}"
