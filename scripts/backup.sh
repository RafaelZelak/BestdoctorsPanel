#!/bin/bash
# Automated backup script for BestDoctors production database

set -e

# Configuration
BACKUP_DIR="/var/backups/bestdoctors"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Load production environment
if [ -f .env.production ]; then
    export $(cat .env.production | grep -v '^#' | xargs)
else
    echo "❌ .env.production not found!"
    exit 1
fi

mkdir -p "$BACKUP_DIR"

echo "🗄️  Starting backup: $DATE"

# PostgreSQL backup
echo "📦 Backing up PostgreSQL..."
PGPASSWORD=$PG_PASSWORD pg_dump \
    -h $PG_HOST \
    -p $PG_PORT \
    -U $PG_USER \
    -d $PG_DATABASE \
    --format=custom \
    --compress=9 \
    --file="$BACKUP_DIR/postgres_$DATE.dump"

# Verify backup
if [ -f "$BACKUP_DIR/postgres_$DATE.dump" ]; then
    SIZE=$(du -h "$BACKUP_DIR/postgres_$DATE.dump" | cut -f1)
    echo "✅ PostgreSQL backup completed: $SIZE"
else
    echo "❌ PostgreSQL backup failed!"
    exit 1
fi

# Redis backup
echo "📦 Backing up Redis..."
docker exec bestdoctors_redis redis-cli --no-auth-warning -a "$REDIS_PASSWORD" BGSAVE
sleep 5
docker cp bestdoctors_redis:/data/dump.rdb "$BACKUP_DIR/redis_$DATE.rdb"

if [ -f "$BACKUP_DIR/redis_$DATE.rdb" ]; then
    SIZE=$(du -h "$BACKUP_DIR/redis_$DATE.rdb" | cut -f1)
    echo "✅ Redis backup completed: $SIZE"
else
    echo "❌ Redis backup failed!"
    exit 1
fi

# Optional: Upload to S3
if [ ! -z "$AWS_S3_BUCKET" ]; then
    echo "☁️  Uploading to S3..."
    aws s3 cp "$BACKUP_DIR/postgres_$DATE.dump" "s3://$AWS_S3_BUCKET/backups/postgres_$DATE.dump"
    aws s3 cp "$BACKUP_DIR/redis_$DATE.rdb" "s3://$AWS_S3_BUCKET/backups/redis_$DATE.rdb"
    echo "✅ Uploaded to S3"
fi

# Cleanup old backups
echo "🗑️  Cleaning up old backups (older than $RETENTION_DAYS days)..."
find "$BACKUP_DIR" -name "postgres_*.dump" -mtime +$RETENTION_DAYS -delete
find "$BACKUP_DIR" -name "redis_*.rdb" -mtime +$RETENTION_DAYS -delete

echo "✅ Backup completed successfully!"
echo "Location: $BACKUP_DIR"
