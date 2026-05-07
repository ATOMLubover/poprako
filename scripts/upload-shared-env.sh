#!/bin/sh
set -eu

SOURCE_ENV_FILE=${SOURCE_ENV_FILE:-.env}
DEPLOY_ROOT=${DEPLOY_ROOT:-/opt/poprako-s}
SERVER_USER=${SERVER_USER:?SERVER_USER is required}
SERVER_HOST=${SERVER_HOST:?SERVER_HOST is required}
IMAGE_TAG=${IMAGE_TAG:-$(git rev-parse --short=12 HEAD)}

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT_DIR=$(CDPATH= cd -- "${SCRIPT_DIR}/.." && pwd)
REMOTE_ENV_FILE="${DEPLOY_ROOT}/shared/.env"

case "$SOURCE_ENV_FILE" in
    /*) SOURCE_ENV_PATH="$SOURCE_ENV_FILE" ;;
    *) SOURCE_ENV_PATH="${ROOT_DIR}/${SOURCE_ENV_FILE}" ;;
esac

[ -f "$SOURCE_ENV_PATH" ] || {
    echo "Missing source env file: ${SOURCE_ENV_PATH}" >&2
    exit 1
}

read_env() {
    key=$1
    default=${2-}
    required=${3:-0}
    value=""
    found=0

    while IFS= read -r line || [ -n "$line" ]; do
        case "$line" in
            ''|'#'*)
                continue
                ;;
            "$key="*)
                value=${line#*=}
                found=1
                ;;
        esac
    done < "$SOURCE_ENV_PATH"

    if [ "$found" -eq 0 ]; then
        if [ "$required" -eq 1 ]; then
            echo "Missing required key '${key}' in ${SOURCE_ENV_PATH}" >&2
            exit 1
        fi

        printf '%s' "$default"
        return 0
    fi

    case "$value" in
        \"*\")
            value=${value#\"}
            value=${value%\"}
            ;;
        \'*\')
            value=${value#\'}
            value=${value%\'}
            ;;
    esac

    if [ -z "$value" ] && [ "$required" -eq 1 ]; then
        echo "Required key '${key}' is empty in ${SOURCE_ENV_PATH}" >&2
        exit 1
    fi

    if [ -z "$value" ]; then
        printf '%s' "$default"
        return 0
    fi

    printf '%s' "$value"
}

append_env() {
    key=$1
    value=$2

    printf '%s=%s\n' "$key" "$value" >> "$TMP_ENV"
}

postgresPassword=$(read_env POSTGRES_PASSWORD "" 1)
databasePassword=$(read_env DATABASE_PASSWORD "$postgresPassword")
jwtSecret=$(read_env JWT_SECRET "" 1)
jwtExpirationHours=$(read_env JWT_EXPIRATION_HOURS "336")
ossPlatform=$(read_env OSS_PLATFORM "" 1)

if [ "$ossPlatform" != "r2" ]; then
    echo "Unsupported OSS_PLATFORM '${ossPlatform}'. Current production build supports only 'r2'." >&2
    exit 1
fi

r2AccountId=$(read_env R2_ACCOUNT_ID "" 1)
r2AccessKeyId=$(read_env R2_ACCESS_KEY_ID "" 1)
r2SecretAccessKey=$(read_env R2_SECRET_ACCESS_KEY "" 1)
r2BucketName=$(read_env R2_BUCKET_NAME "" 1)
r2Region=$(read_env R2_REGION "auto")
r2CustomDomain=$(read_env R2_CUSTOM_DOMAIN "")

aliyunAccessKeyId=$(read_env ALIYUN_OSS_ACCESS_KEY_ID "")
aliyunAccessKeySecret=$(read_env ALIYUN_OSS_ACCESS_KEY_SECRET "")
aliyunRegion=$(read_env ALIYUN_OSS_REGION "")
aliyunBucketName=$(read_env ALIYUN_OSS_BUCKET_NAME "")
aliyunEndpoint=$(read_env ALIYUN_OSS_ENDPOINT "")
aliyunCustomDomain=$(read_env ALIYUN_OSS_CUSTOM_DOMAIN "")

TMP_ENV=$(mktemp)
trap 'rm -f "$TMP_ENV"' EXIT INT TERM HUP

append_env IMAGE_TAG "$IMAGE_TAG"
append_env APP_ENV "prod"
append_env POSTGRES_PASSWORD "$postgresPassword"
append_env DATABASE_USER "poprako_s"
append_env DATABASE_NAME "db_poprako_s"
append_env DATABASE_PASSWORD "$databasePassword"
append_env DATABASE_HOST "prod-postgres"
append_env DATABASE_PORT "5432"
append_env JWT_SECRET "$jwtSecret"
append_env JWT_EXPIRATION_HOURS "$jwtExpirationHours"
append_env OSS_PLATFORM "$ossPlatform"
append_env R2_ACCOUNT_ID "$r2AccountId"
append_env R2_ACCESS_KEY_ID "$r2AccessKeyId"
append_env R2_SECRET_ACCESS_KEY "$r2SecretAccessKey"
append_env R2_BUCKET_NAME "$r2BucketName"
append_env R2_REGION "$r2Region"
append_env R2_CUSTOM_DOMAIN "$r2CustomDomain"
append_env ALIYUN_OSS_ACCESS_KEY_ID "$aliyunAccessKeyId"
append_env ALIYUN_OSS_ACCESS_KEY_SECRET "$aliyunAccessKeySecret"
append_env ALIYUN_OSS_REGION "$aliyunRegion"
append_env ALIYUN_OSS_BUCKET_NAME "$aliyunBucketName"
append_env ALIYUN_OSS_ENDPOINT "$aliyunEndpoint"
append_env ALIYUN_OSS_CUSTOM_DOMAIN "$aliyunCustomDomain"

ssh "${SERVER_USER}@${SERVER_HOST}" "mkdir -p '${DEPLOY_ROOT}/shared'"
scp "$TMP_ENV" "${SERVER_USER}@${SERVER_HOST}:${REMOTE_ENV_FILE}"
ssh "${SERVER_USER}@${SERVER_HOST}" "chmod 600 '${REMOTE_ENV_FILE}'"

printf '%s\n' "Uploaded runtime env to ${SERVER_USER}@${SERVER_HOST}:${REMOTE_ENV_FILE}"
