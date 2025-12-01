#!/bin/bash
# Production deployment script for BestDoctors Panel
# Handles SSL, secrets, health checks, and safe deployment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "BestDoctors Production Deployment"

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
   echo -e "${RED} Do not run this script as root${NC}"
   exit 1
fi

# Check for production environment file
if [ ! -f .env.production ]; then
    echo -e "${RED} .env.production not found!${NC}"
    echo "Please create .env.production from .env.production.example"
    exit 1
fi

# Load production environment
export $(cat .env.production | grep -v '^#' | xargs)

# Validate required variables
REQUIRED_VARS=(
    "TRAEFIK_DOMAIN"
    "SUPERADMIN_USERNAME"
    "SUPERADMIN_PASSWORD"
    "REDIS_PASSWORD"
    "PG_HOST"
    "PG_PASSWORD"
)

echo "Validating environment variables..."
for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        echo -e "${RED} Missing required variable: $var${NC}"
        exit 1
    fi
done
echo -e "${GREEN} All required variables present${NC}"

# Check for placeholder passwords
if [[ "$SUPERADMIN_PASSWORD" == *"CHANGE_ME"* ]]; then
    echo -e "${RED} Please change default passwords in .env.production${NC}"
    exit 1
fi

# Prepare Traefik SSL directory
echo " Preparing SSL certificates directory..."
mkdir -p traefik
touch traefik/acme.json
chmod 600 traefik/acme.json

# Pre-flight checks
echo "Running pre-flight checks..."

# Check if domain is reachable
if ! ping -c 1 "$TRAEFIK_DOMAIN" &> /dev/null; then
    echo -e "${YELLOW} Warning: Domain $TRAEFIK_DOMAIN is not reachable${NC}"
    echo "Make sure DNS is configured correctly"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check if ports 80 and 443 are available
if ss -tuln | grep -q ':80\s'; then
    echo -e "${RED} Port 80 is already in use${NC}"
    exit 1
fi

if ss -tuln | grep -q ':443\s'; then
    echo -e "${RED} Port 443 is already in use${NC}"
    exit 1
fi

echo -e "${GREEN} Pre-flight checks passed${NC}"

# Build and deploy
echo "  Building Docker images..."
docker-compose -f docker-compose.prod.yml build --no-cache

echo " Starting production services..."
docker-compose -f docker-compose.prod.yml up -d

echo " Waiting for services to be healthy..."
sleep 10

# Health check
MAX_RETRIES=30
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -f -k "https://$TRAEFIK_DOMAIN/health" &> /dev/null; then
        echo -e "${GREEN} All services are healthy!${NC}"
        break
    fi
    
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "Waiting for services... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 5
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo -e "${RED} Services failed to become healthy${NC}"
    echo "Check logs with: docker-compose -f docker-compose.prod.yml logs"
    exit 1
fi

# Display summary
echo ""
echo -e "${GREEN} Deployment Complete!${NC}"
echo ""
echo " URLs:"
echo "   Frontend: https://$TRAEFIK_DOMAIN"
echo "   Admin:    https://$TRAEFIK_DOMAIN/admin"
echo "   Traefik:  https://traefik.$TRAEFIK_DOMAIN (if configured)"
echo ""
echo " Useful commands:"
echo "   Logs:     docker-compose -f docker-compose.prod.yml logs -f"
echo "   Status:   docker-compose -f docker-compose.prod.yml ps"
echo "   Stop:     docker-compose -f docker-compose.prod.yml down"
echo "   Restart:  docker-compose -f docker-compose.prod.yml restart <service>"
echo ""
echo " SSL Certificates:"
echo "   Let's Encrypt certificates will be automatically obtained"
echo "   Certificate storage: ./traefik/acme.json"
echo ""
echo -e "${YELLOW}  Important:${NC}"
echo "   - Monitor logs for the first 24 hours"
echo "   - Set up automated backups"
echo "   - Configure firewall rules"
echo "   - Enable monitoring/alerting"
echo ""
