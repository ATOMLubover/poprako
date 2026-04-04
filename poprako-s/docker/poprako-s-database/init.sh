#!/bin/sh
set -e

echo "Running database migrations..."
for f in $(ls /migrations/*.up.sql | sort); do
    echo "  -> $f"
    psql -v ON_ERROR_STOP=1 \
        --username "$POSTGRES_USER" \
        --dbname "$POSTGRES_DB" \
        -f "$f"
done
echo "Migrations complete."
