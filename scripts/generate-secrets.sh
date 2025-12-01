#!/bin/bash
# Generate secure random secrets for production deployment

set -e

echo "🔐 Generating production secrets..."

# Function to generate random password
generate_password() {
    openssl rand -base64 32 | tr -d "=+/" | cut -c1-32
}

# Check if .env.production exists
if [ -f .env.production ]; then
    echo "⚠️  .env.production already exists!"
    read -p "Overwrite existing file? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
fi

# Copy template
cp .env.production.example .env.production

# Generate secrets
SUPERADMIN_PASSWORD=$(generate_password)
REDIS_PASSWORD=$(generate_password)
PG_PASSWORD=$(generate_password)

# Update the file with generated secrets
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/SUPERADMIN_PASSWORD=CHANGE_ME_STRONG_PASSWORD_HERE/SUPERADMIN_PASSWORD=$SUPERADMIN_PASSWORD/" .env.production
    sed -i '' "s/REDIS_PASSWORD=CHANGE_ME_REDIS_PASSWORD_HERE/REDIS_PASSWORD=$REDIS_PASSWORD/" .env.production
    sed -i '' "s/PG_PASSWORD=CHANGE_ME_DB_PASSWORD_HERE/PG_PASSWORD=$PG_PASSWORD/" .env.production
else
    # Linux
    sed -i "s/SUPERADMIN_PASSWORD=CHANGE_ME_STRONG_PASSWORD_HERE/SUPERADMIN_PASSWORD=$SUPERADMIN_PASSWORD/" .env.production
    sed -i "s/REDIS_PASSWORD=CHANGE_ME_REDIS_PASSWORD_HERE/REDIS_PASSWORD=$REDIS_PASSWORD/" .env.production
    sed -i "s/PG_PASSWORD=CHANGE_ME_DB_PASSWORD_HERE/PG_PASSWORD=$PG_PASSWORD/" .env.production
fi

echo ""
echo "✅ Secrets generated successfully!"
echo ""
echo "================================================"
echo "⚠️  IMPORTANT: Save these credentials securely!"
echo "================================================"
echo ""
echo "SuperAdmin Password: $SUPERADMIN_PASSWORD"
echo "Redis Password:      $REDIS_PASSWORD"
echo "PostgreSQL Password: $PG_PASSWORD"
echo ""
echo "These passwords have been saved to .env.production"
echo ""
echo "Next steps:"
echo "1. Update TRAEFIK_DOMAIN in .env.production with your domain"
echo "2. Update database host (PG_HOST) if using external database"
echo "3. Add Twilio credentials if using WhatsApp messaging"
echo "4. Review and update all other settings"
echo "5. NEVER commit .env.production to git!"
echo ""
