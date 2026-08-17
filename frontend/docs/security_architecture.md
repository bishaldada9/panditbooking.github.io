# Security Architecture

## Overview
The platform implements defense-in-depth security architecture with multiple layers of protection.

## Architecture Layers

### 1. Network Security
- TLS 1.3 termination at Nginx
- Rate limiting at reverse proxy level
- DDoS protection
- Internal network isolation via Docker networks

### 2. Application Security

#### Authentication Layer
- Password hashing with bcrypt (cost 12)
- JWT with RS256 signing
- Token rotation and refresh
- Redis-based blacklist for revoked tokens
- TOTP-based MFA
- Account lockout after 5 failed attempts

#### Authorization Layer
- RBAC with three roles: customer, pandit, admin
- Middleware-enforced permission checks
- JWT claim-based role validation

#### Input Validation Layer
- Request validation using go-playground/validator
- Custom Nepali phone number validator
- SQL injection detection
- XSS sanitization
- Request size limits

### 3. Data Security
- PostgreSQL with encrypted connections
- Redis with authentication
- Environment variable configuration
- No hardcoded credentials

### 4. Audit & Monitoring
- Structured logging via zap
- Comprehensive audit trails
- Request logging with correlation IDs
- Panic recovery middleware

## Security Flow Diagram

```
Client → HTTPS → Nginx (TLS, Rate Limit, Security Headers)
                → Gin Router (CORS, CSRF, Auth)
                → Middleware Chain (Logger, Validator, Rate Limiter, Auth)
                → Handler → Service → Repository → Database
                                                   → Redis
                → Audit Logger → PostgreSQL Audit Table
```

## Compliance
- OWASP Top 10 (2021)
- PCI DSS considerations for payment data
- GDPR/Data protection principles
