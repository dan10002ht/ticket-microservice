#!/bin/bash
set -e

echo "=== Initializing PostgreSQL Replica ==="

# Environment variables (set by docker-compose)
PRIMARY_HOST=${PRIMARY_HOST:-postgres-primary}
PRIMARY_PORT=${PRIMARY_PORT:-5432}
REPLICATION_USER=${REPLICATION_USER:-replicator}
REPLICATION_PASSWORD=${REPLICATION_PASSWORD:-replicator_pass}
REPLICA_NAME=${REPLICA_NAME:-replica1}

# Check if data directory is empty (first run)
if [ -z "$(ls -A /var/lib/postgresql/data 2>/dev/null)" ]; then
    echo "Data directory is empty. Setting up replica from primary..."

    # Wait for primary to be ready
    echo "Waiting for primary at $PRIMARY_HOST:$PRIMARY_PORT..."
    until PGPASSWORD=$REPLICATION_PASSWORD pg_isready -h $PRIMARY_HOST -p $PRIMARY_PORT -U $REPLICATION_USER; do
        echo "Primary not ready, waiting..."
        sleep 5
    done
    echo "Primary is ready!"

    # Create base backup from primary
    echo "Creating base backup from primary..."
    PGPASSWORD=$REPLICATION_PASSWORD pg_basebackup \
        -h $PRIMARY_HOST \
        -p $PRIMARY_PORT \
        -U $REPLICATION_USER \
        -D /var/lib/postgresql/data \
        -Fp \
        -Xs \
        -P \
        -R

    # Configure standby.signal (PostgreSQL 12+)
    echo "Configuring standby mode..."
    touch /var/lib/postgresql/data/standby.signal

    # Configure primary connection info
    echo "Configuring primary connection info..."
    cat >> /var/lib/postgresql/data/postgresql.auto.conf <<EOF

# Replica configuration
primary_conninfo = 'host=$PRIMARY_HOST port=$PRIMARY_PORT user=$REPLICATION_USER password=$REPLICATION_PASSWORD application_name=$REPLICA_NAME'
primary_slot_name = 'replica_slot_1'
EOF

    # Set permissions
    chown -R postgres:postgres /var/lib/postgresql/data
    chmod 700 /var/lib/postgresql/data

    echo "Replica setup completed!"
else
    echo "Data directory is not empty. Using existing data."

    # Ensure standby.signal exists
    if [ ! -f /var/lib/postgresql/data/standby.signal ]; then
        echo "Creating standby.signal..."
        touch /var/lib/postgresql/data/standby.signal
    fi
fi

echo "=== Replica initialization completed! ==="
