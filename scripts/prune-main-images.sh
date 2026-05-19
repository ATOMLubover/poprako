#!/bin/sh
set -eu

MAIN_IMAGE=${MAIN_IMAGE:-poprako-s-main}
KEEP_COUNT=${KEEP_COUNT:-4}

case "$KEEP_COUNT" in
    ''|*[!0-9]*)
        echo "KEEP_COUNT must be a non-negative integer" >&2
        exit 1
        ;;
esac

if [ "$KEEP_COUNT" -eq 0 ]; then
    old_tags=$(docker image ls "$MAIN_IMAGE" --format '{{.Tag}}')
else
    old_tags=$(docker image ls "$MAIN_IMAGE" --format '{{.Tag}}' | awk -v keep="$KEEP_COUNT" 'NR > keep')
fi

[ -n "$old_tags" ] || {
    printf '%s\n' "No old ${MAIN_IMAGE} tags to remove"
    exit 0
}

printf '%s\n' "$old_tags" | while IFS= read -r tag; do
    [ -n "$tag" ] || continue

    docker image rm "${MAIN_IMAGE}:${tag}"
done
