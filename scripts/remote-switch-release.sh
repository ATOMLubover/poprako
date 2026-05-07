#!/bin/sh
set -eu

IMAGE_TAG=${IMAGE_TAG:?IMAGE_TAG is required}
DEPLOY_ROOT=${DEPLOY_ROOT:-/opt/poprako-s}
MAIN_IMAGE=${MAIN_IMAGE:-poprako-s-main}
DATABASE_IMAGE=${DATABASE_IMAGE:-poprako-s-database}

RELEASE_DIR="${DEPLOY_ROOT}/releases/${IMAGE_TAG}"
SHARED_DIR="${DEPLOY_ROOT}/shared"
ENV_FILE="${SHARED_DIR}/.env"
COMPOSE_FILE="${SHARED_DIR}/docker/compose.prod.yml"
MIGRATION_BUNDLE="${RELEASE_DIR}/poprako-s-migrations-${IMAGE_TAG}.tar.gz"
MAIN_BUNDLE="${RELEASE_DIR}/${MAIN_IMAGE}-${IMAGE_TAG}.tar.gz"
DATABASE_BUNDLE="${RELEASE_DIR}/${DATABASE_IMAGE}-${IMAGE_TAG}.tar.gz"

require_env_key() {
    key=$1

    if ! grep -Eq "^${key}=.+$" "$ENV_FILE"; then
        echo "Missing required env key '${key}' in ${ENV_FILE}" >&2
        exit 1
    fi
}

[ -f "$ENV_FILE" ] || {
    echo "Missing ${ENV_FILE}. Create it first with runtime secrets." >&2
    exit 1
}

[ -f "$MAIN_BUNDLE" ] || {
    echo "Missing ${MAIN_BUNDLE}" >&2
    exit 1
}

[ -f "$DATABASE_BUNDLE" ] || {
    echo "Missing ${DATABASE_BUNDLE}" >&2
    exit 1
}

[ -f "$MIGRATION_BUNDLE" ] || {
    echo "Missing ${MIGRATION_BUNDLE}" >&2
    exit 1
}

mkdir -p "$SHARED_DIR"

docker load -i "$DATABASE_BUNDLE"
docker load -i "$MAIN_BUNDLE"
tar -xzf "$MIGRATION_BUNDLE" -C "$SHARED_DIR"

TMP_ENV="${ENV_FILE}.tmp"
grep -v '^IMAGE_TAG=' "$ENV_FILE" > "$TMP_ENV" || true
printf 'IMAGE_TAG=%s\n' "$IMAGE_TAG" >> "$TMP_ENV"
mv "$TMP_ENV" "$ENV_FILE"
chmod 600 "$ENV_FILE"

require_env_key APP_ENV
require_env_key POSTGRES_PASSWORD
require_env_key DATABASE_USER
require_env_key DATABASE_NAME
require_env_key DATABASE_PASSWORD
require_env_key DATABASE_HOST
require_env_key DATABASE_PORT
require_env_key JWT_SECRET
require_env_key JWT_EXPIRATION_HOURS
require_env_key OSS_PLATFORM
require_env_key R2_ACCOUNT_ID
require_env_key R2_ACCESS_KEY_ID
require_env_key R2_SECRET_ACCESS_KEY
require_env_key R2_BUCKET_NAME
require_env_key R2_REGION

docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" stop prod-main-server || true
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d --wait --force-recreate prod-postgres
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" run --rm prod-db-migrate
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d --force-recreate prod-main-server
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps
