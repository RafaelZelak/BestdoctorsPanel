# BestDoctors Production Environment

## Production Files

- `.env.production` - Production environment variables (NOT in git)
- `docker-compose.prod.yml` - Production Docker configuration
- `start.prod.sh` - Production deployment script

## Development Files (Local Only)

- `.env` - Development environment variables
- `docker-compose.yml` - Development Docker configuration
- `start.sh` - Development startup script

## Security

**.env.production should NEVER be committed to git!**

All production secrets should be:

- Generated with `scripts/generate-secrets.sh`
- Stored securely (password manager, AWS Secrets Manager, etc.)
- Rotated regularly (every 90 days minimum)

## Deployment

### First Time Setup

1. **Generate production secrets:**

   ```bash
   chmod +x scripts/generate-secrets.sh
   ./scripts/generate-secrets.sh
   ```

2. **Update `.env.production`:**

   - Set `TRAEFIK_DOMAIN` to your domain
   - Update database credentials if using external DB
   - Add Twilio credentials if needed

3. **Prepare Traefik for SSL:**

   ```bash
   mkdir -p traefik
   touch traefik/acme.json
   chmod 600 traefik/acme.json
   ```

4. **Deploy to production:**
   ```bash
   chmod +x start.prod.sh
   ./start.prod.sh
   ```

### SSL Certificate Notes

- **Let's Encrypt** will automatically obtain certificates
- First run may take a few minutes for certificate issuance
- Certificates auto-renew every 60 days
- Email notifications sent to `admin@example.com` (update in `traefik/traefik.yml`)

### Manual Certificate (if needed)

If you want to use existing certificates instead of Let's Encrypt:

1. Place certificates in `traefik/certs/`:

   - `cert.pem` (certificate)
   - `key.pem` (private key)

2. Update `traefik/traefik.yml` to use file provider instead of ACME

## Monitoring & Maintenance

### View Logs

```bash
docker-compose -f docker-compose.prod.yml logs -f [service]
```

### Check Service Health

```bash
docker-compose -f docker-compose.prod.yml ps
```

### Restart Services

```bash
docker-compose -f docker-compose.prod.yml restart [service]
```

### Backup Database

```bash
chmod +x scripts/backup.sh
./scripts/backup.sh
```

Setup automated daily backups with cron:

```bash
0 2 * * * /path/to/bestdoctors_panel/scripts/backup.sh >> /var/log/bestdoctors_backup.log 2>&1
```

## Security Checklist

- [ ] All default passwords changed
- [ ] Firewall configured (only ports 80, 443 open)
- [ ] SSH key-only authentication
- [ ] Fail2ban installed and configured
- [ ] Automated backups enabled
- [ ] SSL certificates valid
- [ ] Resource limits configured
- [ ] Monitoring/alerting setup

## Troubleshooting

### SSL Certificate Issues

```bash
# Check certificate status
docker-compose -f docker-compose.prod.yml exec traefik cat /acme.json

# Force certificate regeneration
rm traefik/acme.json
touch traefik/acme.json
chmod 600 traefik/acme.json
docker-compose -f docker-compose.prod.yml restart traefik
```

### Service Won't Start

```bash
# Check logs
docker-compose -f docker-compose.prod.yml logs [service]

# Verify environment variables
docker-compose -f docker-compose.prod.yml config
```

### Database Connection Issues

```bash
# Test database connectivity
docker run -it --rm postgres:alpine psql -h $PG_HOST -U $PG_USER -d $PG_DATABASE
```
