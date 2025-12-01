#!/bin/bash

# Stop on error
set -e

echo " Cleaning system environment variables to ensure .env is used..."
unset SUPERADMIN_USERNAME
unset SUPERADMIN_PASSWORD
unset REDIS_PASSWORD
unset PG_HOST
unset PG_PORT
unset PG_USER
unset PG_PASSWORD
unset PG_DATABASE
unset TWILIO_URL
unset TWILIO_ACCOUNT_SID
unset TWILIO_AUTH_TOKEN

echo " Fixing line endings (CRLF -> LF) in .env and docker-compose.yml..."
# Convert CRLF to LF in .env
if [ -f .env ]; then
    sed -i 's/\r$//' .env
fi

# Convert CRLF to LF in docker-compose.yml
if [ -f docker-compose.yml ]; then
    sed -i 's/\r$//' docker-compose.yml
fi

echo "  Cleaning up everything (containers, volumes, orphans)..."
docker-compose down --volumes --remove-orphans

if [ "$1" == "dev" ]; then
    echo "Starting in DEVELOPMENT mode..."
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
else
    echo "Starting in PRODUCTION mode..."
    docker-compose up --build -d
    echo "   Services started in background."
    echo "   Backend: http://localhost:9002"
    echo "   Frontend: http://localhost:80"
    echo "   Logs: docker-compose logs -f"
fi
