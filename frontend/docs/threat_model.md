# Threat Model - Bishal Puja Sewa

## Overview
This document outlines the threat model for the Secure Online Hindu Ritual Service platform following STRIDE methodology.

## Assets
1. User credentials and personal data
2. Payment information and transaction data
3. Pandit verification documents
4. Booking and scheduling data
5. Audit logs
6. JWT tokens and session data

## Threats (STRIDE)

### Spoofing
- **Threat**: Attacker impersonates a verified pandit
- **Mitigation**: Multi-factor authentication, document verification, JWT with strong signing

### Tampering
- **Threat**: Modification of booking or payment data in transit
- **Mitigation**: TLS 1.3, request signing, CSRF tokens, input validation

### Repudiation
- **Threat**: User denies performing an action
- **Mitigation**: Comprehensive audit logging with timestamps, IP, user agent

### Information Disclosure
- **Threat**: Exposure of sensitive user data
- **Mitigation**: Encryption at rest, HTTPS, secure headers, data minimization

### Denial of Service
- **Threat**: Overwhelming the API with requests
- **Mitigation**: Rate limiting, request size limits, connection pooling

### Elevation of Privilege
- **Threat**: User accessing admin functionality
- **Mitigation**: RBAC, middleware enforcement, JWT role validation

## Risk Assessment

| Threat | Impact | Likelihood | Risk Level |
|--------|--------|------------|------------|
| SQL Injection | Critical | Low | High |
| XSS | High | Medium | High |
| Broken Authentication | Critical | Medium | Critical |
| Sensitive Data Exposure | High | Medium | High |
| Brute Force | Medium | High | Medium |

## Mitigation Summary

1. **SQL Injection**: Parameterized queries via GORM
2. **XSS**: Output encoding, CSP headers
3. **CSRF**: Token implementation
4. **Authentication**: bcrypt, MFA, account lockout
5. **Authorization**: JWT with role claims
6. **Session Management**: Redis blacklist, token rotation
7. **Data Protection**: TLS, encrypted storage
8. **Monitoring**: Audit logs, structured logging
