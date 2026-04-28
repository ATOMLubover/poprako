#!/bin/sh
set -eu

MIGRATIONS_DIR="${MIGRATIONS_DIR:-/migrations}"
TRACK_TABLE="${MIGRATION_TRACK_TABLE:-schema_migration_table}"
BASELINE_TABLE="${MIGRATION_BASELINE_TABLE:-public.chapter_invitation}"

set -- "$MIGRATIONS_DIR"/*.up.sql
if [ ! -e "$1" ]; then
    echo "No migration files found in: $MIGRATIONS_DIR" >&2
    exit 1
fi

if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "Migrations directory not found: $MIGRATIONS_DIR" >&2
    exit 1
fi

echo "Ensuring migration tracking table exists..."
psql -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE IF NOT EXISTS public.${TRACK_TABLE} (
    version TEXT PRIMARY KEY,
    checksum TEXT NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
SQL

table_count="$(psql -tAqc "SELECT count(*) FROM public.${TRACK_TABLE}")"
baseline_exists="$(psql -tAqc "SELECT to_regclass('${BASELINE_TABLE}')")"

if [ "$table_count" = "0" ] && [ "$baseline_exists" = "$BASELINE_TABLE" ]; then
    echo "Detected pre-existing schema. Baseline current migrations as applied."
    for f in $(printf '%s\n' "$MIGRATIONS_DIR"/*.up.sql | sort); do
        version="$(basename "$f")"
        checksum="$(sha256sum "$f" | awk '{print $1}')"

        psql -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.${TRACK_TABLE} (version, checksum)
VALUES ('${version}', '${checksum}')
ON CONFLICT (version) DO NOTHING;
SQL

        echo "  -> baselined $version"
    done

    echo "Baseline complete."
fi

echo "Running incremental migrations..."
for f in $(printf '%s\n' "$MIGRATIONS_DIR"/*.up.sql | sort); do
    version="$(basename "$f")"
    checksum="$(sha256sum "$f" | awk '{print $1}')"
    db_checksum="$(psql -tAqc "SELECT checksum FROM public.${TRACK_TABLE} WHERE version='${version}' LIMIT 1")"

    if [ -n "$db_checksum" ]; then
        if [ "$db_checksum" != "$checksum" ]; then
            echo "Checksum mismatch for applied migration: $version" >&2
            echo "  recorded: $db_checksum" >&2
            echo "  current : $checksum" >&2
            exit 1
        fi

        echo "  -> skip $version (already applied)"
        continue
    fi

    echo "  -> apply $version"
    psql -v ON_ERROR_STOP=1 -f "$f"

    psql -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.${TRACK_TABLE} (version, checksum)
VALUES ('${version}', '${checksum}');
SQL
done

echo "Incremental migration complete."
