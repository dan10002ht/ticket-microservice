#!/bin/bash
# Check email-worker PgPool-II updates

echo "🔍 Checking email-worker PgPool-II updates..."

# Check database connection file
echo "📊 Checking email-worker/database/connection.go..."
if grep -q "PgPool-II" email-worker/database/connection.go; then
    echo "  ✅ Database connection updated for PgPool-II"
else
    echo "  ❌ Database connection not updated"
fi

if grep -q "Connection \*sqlx.DB" email-worker/database/connection.go; then
    echo "  ✅ Simplified to single connection"
else
    echo "  ❌ Still using master/slave pattern"
fi

# Check config file
echo "📊 Checking email-worker/config/config.go..."
if grep -q "PgPool-II" email-worker/config/config.go; then
    echo "  ✅ Config updated for PgPool-II"
else
    echo "  ❌ Config not updated"
fi

if ! grep -q "MasterHost\|SlaveHost" email-worker/config/config.go; then
    echo "  ✅ Removed master/slave config fields"
else
    echo "  ❌ Still has master/slave config fields"
fi

# Check environment file
echo "📊 Checking email-worker/env.example..."
if grep -q "pgpool-ticket" email-worker/env.example; then
    echo "  ✅ Environment updated for PgPool-II"
else
    echo "  ❌ Environment not updated"
fi

# Check dev script
echo "📊 Checking email-worker/scripts/dev-local.sh..."
if grep -q "pgpool" email-worker/scripts/dev-local.sh; then
    echo "  ✅ Dev script updated for PgPool-II"
else
    echo "  ❌ Dev script not updated"
fi

echo ""
echo "📋 Email Worker PgPool-II Update Summary:"
echo "  - Database connection: Single connection to PgPool-II"
echo "  - Config: Simplified to PgPool-II endpoints"
echo "  - Environment: Updated to use PgPool-II"
echo "  - Dev script: Updated to start PgPool-II infrastructure"
echo ""
echo "💡 Benefits:"
echo "  - Automatic master/slave routing by PgPool-II"
echo "  - Simplified configuration"
echo "  - Better performance with connection pooling"
echo "  - Automatic failover"
echo ""
echo "✅ Email Worker PgPool-II update check completed!"
