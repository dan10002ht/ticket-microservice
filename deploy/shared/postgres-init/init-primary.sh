#!/bin/bash
set -e

echo "=== Initializing PostgreSQL Primary ==="

# Wait for PostgreSQL to be ready
until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

echo "PostgreSQL is ready"

# Create replication user
echo "Creating replication user..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'replicator') THEN
            CREATE USER replicator WITH REPLICATION LOGIN PASSWORD 'replicator_pass';
        END IF;
    END
    \$\$;
EOSQL

# Grant necessary permissions
echo "Granting permissions to replicator..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    GRANT CONNECT ON DATABASE $POSTGRES_DB TO replicator;
    GRANT USAGE ON SCHEMA public TO replicator;
    GRANT SELECT ON ALL TABLES IN SCHEMA public TO replicator;
    GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO replicator;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO replicator;
EOSQL

# Create replication slots (for staging/production)
echo "Creating replication slots..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    SELECT CASE
        WHEN NOT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = 'replica_slot_1')
        THEN pg_create_physical_replication_slot('replica_slot_1', true)
    END;

    SELECT CASE
        WHEN NOT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = 'replica_slot_2')
        THEN pg_create_physical_replication_slot('replica_slot_2', true)
    END;
EOSQL

# Configure pg_hba.conf for replication
echo "Configuring pg_hba.conf for replication..."
if ! grep -q "host replication replicator" /var/lib/postgresql/data/pg_hba.conf; then
    cat >> /var/lib/postgresql/data/pg_hba.conf <<EOF

# Replication connections
host    replication     replicator      0.0.0.0/0             md5
host    replication     replicator      ::/0                  md5
host    all             all             0.0.0.0/0             md5
EOF
fi

# Reload configuration
echo "Reloading configuration..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -c "SELECT pg_reload_conf();"

# Verify configuration
echo "Verifying configuration..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -c "
SELECT name, setting FROM pg_settings
WHERE name IN ('wal_level', 'max_wal_senders', 'max_replication_slots', 'hot_standby');
"

echo "=== Primary initialization completed! ==="
