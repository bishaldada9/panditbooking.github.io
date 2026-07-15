# Bishal Puja Sewa - Secure Online Hindu Ritual Service & Verified Pandit Booking Platform

A production-grade secure web application for Kathmandu Valley that allows Hindu households to find, verify, book, and securely pay verified pandits for religious ceremonies. Built with modern cybersecurity principles and SSDLC practices.

> **Undergraduate Cybersecurity Thesis Project**
> Demonstrates: SSDLC, OWASP Top 10 mitigation, RBAC, MFA, secure payment workflow, audit logging, secure API design, JWT lifecycle, password security, input validation, output encoding, threat mitigation, trust and transparency.

## 🚀 Quick Start

### Prerequisites

- **Go 1.24+** - [Download](https://go.dev/dl/)
- **Node.js 20+** - [Download](https://nodejs.org/)
- **Docker & Docker Compose** - [Download](https://docs.docker.com/get-docker/) (recommended for PostgreSQL/Redis)

### One-Click Setup

#### Windows (PowerShell)
```powershell
.\setup.ps1
```

#### Windows (CMD)
```cmd
setup.bat
```

#### Linux/Mac
```bash
chmod +x setup.sh && ./setup.sh
```

### Manual Setup

```bash
# 1. Clone and enter the project
cd bishal-puja-sewa

# 2. Copy environment configuration
cp backend/.env.example backend/.env

# 3. Start database services (PostgreSQL + Redis)
docker compose up -d postgres redis

# 4. Start backend (Terminal 1)
cd backend
go mod download
go run ./cmd/main.go

# 5. Start frontend (Terminal 2)
cd frontend
npm install
npm run dev
```

### Using Docker (Full Stack)

```bash
docker compose up -d
```

Access:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

## 📋 Default Credentials

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@bishalpujasewa.com | AdminPass123! |
| Pandit | pandit.ram@example.com | PanditPass123! |
| Customer | customer@example.com | CustomerPass123! |

> Run `go run ./backend/scripts/seed.go` to seed the database with sample data.

## 🏗 Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Client (React + TypeScript)              │
└────────────────────┬────────────────────────────────────────┘
                     │ HTTPS
┌────────────────────▼────────────────────────────────────────┐
│            Nginx Reverse Proxy (SSL, Rate Limit, Security)   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│              Gin Router (CORS, CSRF, Auth)                   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Middleware Chain                         │   │
│  │  Logger → Validator → Rate Limiter → Auth → RBAC    │   │
│  └──────────────────────────────────────────────────────┘   │
│              │                                               │
│  ┌───────────┴───────────┐                                   │
│  │    Handler Layer      │                                   │
│  ├───────────────────────┤                                   │
│  │    Service Layer      │                                   │
│  ├───────────────────────┤                                   │
│  │   Repository Layer    │                                   │
│  └───────────┬───────────┘                                   │
└──────────────┼───────────────────────────────────────────────┘
               │
     ┌─────────┴─────────┐
     │                   │
┌────▼────┐       ┌─────▼─────┐
│PostgreSQL│       │   Redis   │
└─────────┘       └───────────┘
```

## 🧩 Tech Stack

### Backend
| Technology | Purpose |
|------------|---------|
| Go 1.24 | High-performance backend language |
| Gin Framework | HTTP web framework |
| GORM | ORM for PostgreSQL |
| JWT (golang-jwt) | Token-based authentication |
| bcrypt | Password hashing |
| TOTP (pquerna/otp) | Multi-factor authentication |
| Redis | Session blacklist & caching |
| Zap (uber) | Structured logging |
| PostgreSQL 16 | Primary database |
| Docker | Containerization |

### Frontend
| Technology | Purpose |
|------------|---------|
| React 18 | UI framework |
| TypeScript | Type safety |
| Vite | Build tool |
| TailwindCSS | Styling |
| React Router | Routing |
| React Query | Server state management |
| Zod | Schema validation |
| Axios | HTTP client |

### Infrastructure
| Technology | Purpose |
|------------|---------|
| Docker Compose | Orchestration |
| Nginx | Reverse proxy, SSL termination |
| GitHub Actions | CI/CD |

## 🔒 Security Features

### OWASP Top 10 Protection

| # | Category | Implementation |
|---|----------|----------------|
| 1 | Broken Access Control | RBAC middleware, JWT role validation |
| 2 | Cryptographic Failures | bcrypt (cost 12), TLS 1.3, AES-GCM |
| 3 | Injection | GORM parameterized queries, input sanitization |
| 4 | Insecure Design | Rate limiting, brute force protection |
| 5 | Security Misconfiguration | Environment validation, security headers |
| 6 | Vulnerable Components | Dependency scanning in CI/CD |
| 7 | Auth Failures | JWT rotation, MFA TOTP, account lockout |
| 8 | Data Integrity Failures | Webhook verification, audit logging |
| 9 | Logging Failures | Structured audit logs, request tracking |
| 10 | SSRF | Outbound request restrictions |

### Additional Security
- **JWT Rotation**: Access + Refresh token pair with rotation
- **Token Blacklisting**: Redis-based revocation
- **Multi-Factor Auth**: TOTP with recovery codes
- **Brute Force Protection**: Account lockout after 5 failed attempts
- **Rate Limiting**: Token bucket algorithm per IP/user
- **SQL Injection Detection**: Pattern matching on all inputs
- **XSS Prevention**: Output encoding, CSP headers
- **CSRF Protection**: Token-based prevention
- **Security Headers**: CSP, HSTS, X-Frame-Options, X-Content-Type-Options
- **CORS**: Restricted origins
- **Audit Logging**: All critical actions logged with context
- **Device Tracking**: Session management with device info

## 📁 Project Structure

```
├── backend/
│   ├── cmd/main.go              # Application entry point
│   ├── internal/
│   │   ├── auth/                # Authentication module
│   │   ├── users/               # User management
│   │   ├── pandits/             # Pandit registration & verification
│   │   ├── bookings/            # Booking lifecycle
│   │   ├── rituals/             # Ritual categories & items
│   │   ├── payments/            # Payment gateway abstraction
│   │   ├── reviews/             # Review & rating system
│   │   ├── admin/               # Admin dashboard & management
│   │   ├── audit/               # Audit logging
│   │   ├── notification/        # Notification system
│   │   ├── middleware/          # Security & auth middleware
│   │   └── routes/              # API route definitions
│   ├── pkg/                     # Shared utilities
│   ├── migrations/              # Database migration
│   └── tests/                   # Unit & integration tests
├── frontend/
│   └── src/
│       ├── components/          # Reusable components
│       ├── pages/               # Page components
│       ├── services/            # API client services
│       ├── store/               # Auth context
│       └── types/               # TypeScript types
├── infrastructure/              # Nginx config, monitoring
├── docs/                        # Documentation
└── docker-compose.yml           # Docker orchestration
```

## 📡 API Overview

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/refresh` | Refresh tokens |
| POST | `/api/v1/auth/logout` | Logout |
| POST | `/api/v1/auth/forgot-password` | Request password reset |
| POST | `/api/v1/auth/reset-password` | Reset password |
| POST | `/api/v1/auth/change-password` | Change password |
| POST | `/api/v1/auth/mfa/setup` | Setup MFA |
| POST | `/api/v1/auth/mfa/verify` | Verify MFA |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/profile` | Get profile |
| PUT | `/api/v1/profile` | Update profile |

### Pandits
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/pandits` | List verified pandits |
| GET | `/api/v1/pandits/:id` | Get pandit details |
| POST | `/api/v1/pandit/register` | Register as pandit |
| POST | `/api/v1/pandit/availability` | Set availability |

### Bookings
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/bookings` | Create booking |
| GET | `/api/v1/bookings` | My bookings |
| PUT | `/api/v1/bookings/:id/confirm` | Confirm booking |
| PUT | `/api/v1/bookings/:id/cancel` | Cancel booking |

### Payments
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/payments/initiate` | Initiate payment |
| POST | `/api/v1/payments/verify` | Verify payment |

### Admin
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/dashboard` | Dashboard metrics |
| GET | `/api/v1/admin/users` | List users |
| PUT | `/api/v1/admin/users/:id/suspend` | Suspend user |
| GET | `/api/v1/admin/audit-logs` | View audit logs |

## 🧪 Testing

```bash
# Run backend tests
cd backend
go test -v -race -count=1 ./...

# Run frontend type check
cd frontend
npx tsc --noEmit

# Build frontend
cd frontend
npx vite build
```

## 🔧 Configuration

Configuration is managed via environment variables. See `backend/.env.example` for all options.

### Required Variables
```bash
JWT_SECRET=<64-char-random-string>
JWT_REFRESH_SECRET=<64-char-random-string>
DB_PASSWORD=<strong-database-password>
```

## 📦 Deployment

### Production Build
```bash
docker compose -f docker-compose.yml up -d
```

### Generate Secure Keys
```bash
openssl rand -base64 48   # For JWT secrets
openssl rand -base64 32   # For encryption keys
```

## 📚 Documentation

- [API Documentation](docs/api/api_documentation.md)
- [Threat Model](docs/threat_model.md)
- [Security Architecture](docs/security_architecture.md)
- [Deployment Guide](docs/deployment_guide.md)
- [Postman Collection](postman/collection.json)

## 🎓 Thesis Relevance

This project demonstrates:

1. **Secure Software Development Life Cycle (SSDLC)**
   - Requirements analysis → Design → Implementation → Testing → Deployment
   
2. **OWASP Top 10 Mitigation**
   - Comprehensive protection against all top 10 vulnerabilities

3. **Role-Based Access Control (RBAC)**
   - Customer, Pandit, Admin roles with middleware enforcement

4. **Multi-Factor Authentication (MFA)**
   - TOTP-based with recovery codes

5. **Secure Payment Workflow**
   - Gateway abstraction, webhook verification, refund logic

6. **Comprehensive Audit Logging**
   - Every critical action logged with context

7. **Secure API Design**
   - RESTful, validated, rate-limited, versioned

8. **JWT Lifecycle Management**
   - Rotation, blacklisting, refresh tokens

9. **Password Security**
   - bcrypt hashing, strength validation, rotation policy

10. **Input Validation & Output Encoding**
    - Server-side validation, XSS prevention

## 📄 License

MIT - This project is for educational purposes as part of an undergraduate cybersecurity thesis.
