#!/bin/sh
set -eu

IMAGE_TAG=${IMAGE_TAG:-$(git rev-parse --short=12 HEAD)}
DEPLOY_ROOT=${DEPLOY_ROOT:-/opt/poprako-s}
SERVER_USER=${SERVER_USER:?SERVER_USER is required}
SERVER_HOST=${SERVER_HOST:?SERVER_HOST is required}
MAIN_IMAGE=${MAIN_IMAGE:-poprako-s-main}
DIST_DIR=${DIST_DIR:-dist}
TARGET_PLATFORM=${TARGET_PLATFORM:-linux/amd64}

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT_DIR=$(CDPATH= cd -- "${SCRIPT_DIR}/.." && pwd)
RELEASE_DIR="${DEPLOY_ROOT}/releases/${IMAGE_TAG}"
REMOTE_BIN_DIR="${DEPLOY_ROOT}/shared/bin"
UPLOAD_ENV_SCRIPT="${SCRIPT_DIR}/upload-shared-env.sh"

cd "$ROOT_DIR"

mkdir -p "$DIST_DIR"

docker build --platform "${TARGET_PLATFORM}" -f docker/poprako-s-main/Dockerfile -t "${MAIN_IMAGE}:${IMAGE_TAG}" .

docker save "${MAIN_IMAGE}:${IMAGE_TAG}" | gzip >"${DIST_DIR}/${MAIN_IMAGE}-${IMAGE_TAG}.tar.gz"
tar -czf "${DIST_DIR}/poprako-s-migrations-${IMAGE_TAG}.tar.gz" migrations docker/prod-database-migrate.sh docker/compose.prod.yml

ssh "${SERVER_USER}@${SERVER_HOST}" "mkdir -p '${RELEASE_DIR}' '${REMOTE_BIN_DIR}'"
scp "${DIST_DIR}/${MAIN_IMAGE}-${IMAGE_TAG}.tar.gz" "${SERVER_USER}@${SERVER_HOST}:${RELEASE_DIR}/"
scp "${DIST_DIR}/poprako-s-migrations-${IMAGE_TAG}.tar.gz" "${SERVER_USER}@${SERVER_HOST}:${RELEASE_DIR}/"
scp "${SCRIPT_DIR}/remote-switch-release.sh" "${SERVER_USER}@${SERVER_HOST}:${REMOTE_BIN_DIR}/"
scp "${UPLOAD_ENV_SCRIPT}" "${SERVER_USER}@${SERVER_HOST}:${REMOTE_BIN_DIR}/"
ssh "${SERVER_USER}@${SERVER_HOST}" "chmod 755 '${REMOTE_BIN_DIR}/remote-switch-release.sh'"
ssh "${SERVER_USER}@${SERVER_HOST}" "chmod 755 '${REMOTE_BIN_DIR}/upload-shared-env.sh'"

"${UPLOAD_ENV_SCRIPT}"

printf '%s\n' "Release ${IMAGE_TAG} uploaded to ${SERVER_USER}@${SERVER_HOST}:${RELEASE_DIR}"
printf '%s\n' "Target platform: ${TARGET_PLATFORM}"
printf '%s\n' "Runtime env uploaded to ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_ROOT}/shared/.env"
printf '%s\n' "Run on server:"
printf '%s\n' "IMAGE_TAG=${IMAGE_TAG} DEPLOY_ROOT=${DEPLOY_ROOT} sh ${REMOTE_BIN_DIR}/remote-switch-release.sh"
