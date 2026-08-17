# Production Deployment Guide

## Prerequisites
- Docker & Docker Compose
- Domain name with SSL certificate
- PostgreSQL 16
- Redis 7
- SMTP service (SendGrid)

## Environment Variables

### Required Variables
```bash
JWT_SECRET=<64-char-random-string>
JWT_REFRESH_SECRET=<64-char-random-string>
DB_PASSWORD=<strong-database-password>
```

### Generate Secure Keys
```bash
openssl rand -base64 48  # For JWT secrets
openssl rand -base64 32  # For encryption keys
```

## Deployment Steps

### 1. Clone and Configure
```bash
git clone <repository>
cd bishal-puja-sewa
cp backend/.env.example backend/.env
# Edit .env with production values
```

### 2. SSL Certificates
```bash
# Using Let's Encrypt
docker compose run --rm certbot certonly --webroot -w /var/www/certbot -d yourdomain.com
```

### 3. Database Setup
```bash
# Create database
docker compose exec postgres createdb -U ritual_user hindu_ritual_db

# Run migrations (auto-migrate on startup)
```

### 4. Deploy
```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 5. Monitoring
```bash
# Check logs
docker compose logs -f backend
docker compose logs -f frontend

# Health check
curl https://yourdomain.com/health
```

## Security Checklist
- [ ] JWT secrets changed from defaults
- [ ] Database password changed
- [ ] Redis password configured
- [ ] SSL certificates installed
- [ ] Rate limiting configured
- [ ] CORS origins restricted
- [ ] Email service configured
- [ ] Payment gateway credentials set
- [ ] Audit logging enabled
- [ ] Backup strategy implemented

## Backup Strategy
- Daily database backups
- Weekly full system backups
- Backup retention: 30 days
- Encrypted backup storage
