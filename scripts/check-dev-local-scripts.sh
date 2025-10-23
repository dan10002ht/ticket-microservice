#!/bin/bash

# Script to check all dev-local scripts for PgPool-II compatibility
# This script verifies that all dev-local scripts are updated for PgPool-II

echo "🔍 Checking dev-local scripts for PgPool-II compatibility..."
echo ""

# Function to check if a file contains old master/slave references
check_old_references() {
    local file=$1
    local service_name=$2
    
    echo "📋 Checking $service_name dev-local script..."
    
    # Check for old master/slave references
    if grep -q "Master-Slave\|master-slave\|DB_MASTER\|DB_SLAVE" "$file"; then
        echo "⚠️  Found old master/slave references in $file"
        echo "   Lines containing old references:"
        grep -n "Master-Slave\|master-slave\|DB_MASTER\|DB_SLAVE" "$file" | sed 's/^/     /'
        echo ""
        return 1
    fi
    
    # Check for PgPool-II references
    if grep -q "PgPool-II\|pgpool" "$file"; then
        echo "✅ Found PgPool-II references in $file"
    else
        echo "⚠️  No PgPool-II references found in $file"
        return 1
    fi
    
    echo "✅ $service_name dev-local script looks good"
    echo ""
    return 0
}

# Check auth-service dev-local script
if [ -f "auth-service/scripts/dev-local.sh" ]; then
    check_old_references "auth-service/scripts/dev-local.sh" "Auth Service"
    auth_status=$?
else
    echo "❌ Auth Service dev-local script not found"
    auth_status=1
fi

# Check email-worker dev-local script
if [ -f "email-worker/scripts/dev-local.sh" ]; then
    check_old_references "email-worker/scripts/dev-local.sh" "Email Worker"
    email_status=$?
else
    echo "❌ Email Worker dev-local script not found"
    email_status=1
fi

# Check if there are other dev-local scripts
echo "🔍 Looking for other dev-local scripts..."
find . -name "dev-local.sh" -not -path "./auth-service/*" -not -path "./email-worker/*" | while read -r script; do
    echo "📋 Found additional dev-local script: $script"
    check_old_references "$script" "$(basename $(dirname $script))"
done

echo ""
echo "📊 Summary:"
echo "==========="

if [ $auth_status -eq 0 ]; then
    echo "✅ Auth Service dev-local script: PgPool-II compatible"
else
    echo "❌ Auth Service dev-local script: Needs updates"
fi

if [ $email_status -eq 0 ]; then
    echo "✅ Email Worker dev-local script: PgPool-II compatible"
else
    echo "❌ Email Worker dev-local script: Needs updates"
fi

echo ""
echo "🎯 Next Steps:"
echo "=============="
echo "1. All dev-local scripts should now be PgPool-II compatible"
echo "2. Test each service individually:"
echo "   - cd auth-service && ./scripts/dev-local.sh"
echo "   - cd email-worker && ./scripts/dev-local.sh"
echo "3. Use the main development script:"
echo "   - ./scripts/start-dev-with-pgpool.sh"
echo "4. Monitor logs to ensure PgPool-II connections work properly"
echo ""
echo "✅ Dev-local scripts check completed!"
