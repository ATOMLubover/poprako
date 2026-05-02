#!/bin/sh
set -eu

final_table="public.t_assignment_invitation"

if [ "$(psql -tAqc "SELECT to_regclass('$final_table')")" = "$final_table" ]; then
    echo "Schema already initialized, skipping migrations."
    exit 0
fi

table_count="$(psql -tAqc "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public'")"
if [ "$table_count" != "0" ]; then
    echo "Database contains partial schema but is not fully initialized." >&2
    echo "Reset the dev volume or complete the migration manually before continuing." >&2
    exit 1
fi

echo "Running database setup..."

for f in $(ls /migrations/*.up.sql | sort); do
    echo "Executing $f"
    psql -v ON_ERROR_STOP=1 -f "$f"
done

echo "Setup completed!"
