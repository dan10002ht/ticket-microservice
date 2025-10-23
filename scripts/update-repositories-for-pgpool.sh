#!/bin/bash
# Update all repositories to use this.db instead of this.getSlaveDb() and this.getMasterDb()

echo "🔄 Updating repositories for PgPool-II compatibility..."

# Find all repository files
REPO_FILES=$(find auth-service/src/repositories -name "*.js" -not -name "baseRepository.js")

for file in $REPO_FILES; do
    echo "📝 Updating $file..."
    
    # Replace this.getSlaveDb() with this.db
    sed -i '' 's/this\.getSlaveDb()/this.db/g' "$file"
    
    # Replace this.getMasterDb() with this.db
    sed -i '' 's/this\.getMasterDb()/this.db/g' "$file"
    
    echo "✅ Updated $file"
done

echo ""
echo "🎉 All repositories updated for PgPool-II compatibility!"
echo ""
echo "📊 Changes made:"
echo "  - this.getSlaveDb() → this.db"
echo "  - this.getMasterDb() → this.db"
echo ""
echo "💡 Benefits:"
echo "  - Cleaner code with single database connection"
echo "  - PgPool-II handles master/slave routing automatically"
echo "  - Better performance with connection pooling"
echo ""
