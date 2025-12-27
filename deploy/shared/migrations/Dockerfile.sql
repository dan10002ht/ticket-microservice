# Migration runner for raw SQL files
FROM postgres:15-alpine

WORKDIR /migrations

# Copy migration scripts
COPY run-sql-migrations.sh /run-sql-migrations.sh
RUN chmod +x /run-sql-migrations.sh

ENTRYPOINT ["/run-sql-migrations.sh"]
