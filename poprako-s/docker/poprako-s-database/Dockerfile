FROM postgres:18-alpine

# Copy all migration up-scripts and init entrypoint
COPY migrations/ /migrations/
COPY docker/poprako-s-database/init.sh /docker-entrypoint-initdb.d/01-init.sh

RUN chmod +x /docker-entrypoint-initdb.d/01-init.sh
