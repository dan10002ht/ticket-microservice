#!/bin/bash
# Check simplified database configuration

echo "🔍 Checking simplified database configuration..."

# Check if databaseConfig.js is simplified
echo "📊 Checking auth-service/src/config/databaseConfig.js..."
if grep -q "const db = knex" auth-service/src/config/databaseConfig.js; then
    echo "  ✅ Database config simplified to single connection"
else
    echo "  ❌ Database config not simplified"
fi

# Check if old files are removed
echo "📊 Checking for old database files..."
if [ ! -f "auth-service/src/config/database.js" ]; then
    echo "  ✅ Old database.js removed"
else
    echo "  ❌ Old database.js still exists"
fi

if [ ! -f "auth-service/src/config/databaseConfig.pgpool.js" ]; then
    echo "  ✅ Old databaseConfig.pgpool.js removed"
else
    echo "  ❌ Old databaseConfig.pgpool.js still exists"
fi

# Check if baseRepository uses simplified import
echo "📊 Checking baseRepository.js..."
if grep -q "import { db }" auth-service/src/repositories/baseRepository.js; then
    echo "  ✅ BaseRepository uses simplified import"
else
    echo "  ❌ BaseRepository not updated"
fi

# Check if repositories use this.db
echo "📊 Checking repository usage..."
REPO_FILES=$(find auth-service/src/repositories -name "*.js" -not -name "baseRepository.js" | head -3)
for file in $REPO_FILES; do
    if grep -q "this\.db" "$file"; then
        echo "  ✅ $file uses this.db"
    else
        echo "  ❌ $file not updated"
    fi
done

echo ""
echo "📋 Simplified Database Configuration Summary:"
echo "  - Single database connection (db)"
echo "  - PgPool-II handles master/slave routing"
echo "  - Removed complex master/slave logic"
echo "  - Cleaner, simpler code"
echo ""
echo "💡 Benefits:"
echo "  - Easier to maintain"
echo "  - Better performance with PgPool-II"
echo "  - Automatic load balancing"
echo "  - Automatic failover"
echo ""
echo "✅ Database configuration check completed!"
