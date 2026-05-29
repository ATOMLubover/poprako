#!/bin/sh
set -eu

MIGRATIONS_DIR="${MIGRATIONS_DIR:-/migrations}"
TRACK_TABLE="${MIGRATION_TRACK_TABLE:-schema_migration_table}"

# Verify migration files exist.
set -- "$MIGRATIONS_DIR"/*.up.sql
if [ ! -e "$1" ]; then
    echo "No migration files found in: $MIGRATIONS_DIR" >&2
    exit 1
fi

if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "Migrations directory not found: $MIGRATIONS_DIR" >&2
    exit 1
fi

# Ensure the migration tracking table exists.
echo "Ensuring migration tracking table exists..."
psql -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE IF NOT EXISTS public.${TRACK_TABLE} (
    version  TEXT PRIMARY KEY,
    checksum TEXT NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
SQL

# Baseline: if tracking table is empty but the database already contains
# tables (e.g. a pre-existing dev volume), mark all current migrations as
# applied so subsequent runs only pick up new files.
track_count="$(psql -tAqc "SELECT count(*) FROM public.${TRACK_TABLE}")"
public_table_count="$(psql -tAqc "SELECT count(*) FROM pg_tables WHERE schemaname='public'")"

if [ "$track_count" = "0" ] && [ "$public_table_count" != "0" ]; then
    echo "Detected pre-existing schema (${public_table_count} tables). Baseline current migrations as applied."
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

# Incremental apply: only run migrations that have not been recorded.
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
