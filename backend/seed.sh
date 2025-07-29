#!/bin/bash

set -e

echo "🌱 Seeding Kyooar database..."

# Check if --force flag is provided
FORCE_FLAG=""
if [ "$1" == "--force" ]; then
    FORCE_FLAG="--force"
    echo "⚠️  Force mode: Will recreate existing accounts and data"
fi

# Run subscription plans seeding first
echo "📋 Seeding subscription plans..."
go run cmd/seed-plans/main.go

# Run main data seeding
echo "🏢 Seeding organizations, products, and feedback..."
go run cmd/seed/main.go $FORCE_FLAG

echo "✅ Database seeding completed!"
echo ""
echo "🔑 Login credentials (password: Pass123!):"
echo "  Starter: admin_starter@kyooar.com / viewer_starter@kyooar.com"
echo "  Professional: admin_professional@kyooar.com / viewer_professional@kyooar.com"
echo "  Premium: admin_premium@kyooar.com / viewer_premium@kyooar.com"
echo ""
echo "📊 Each plan has 2 organizations with multiple products, QR codes, and feedback"