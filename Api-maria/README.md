# Co-working Spaces — Microservices Platform

A production-ready microservices backend for a co-working space management system, orchestrated with Docker Compose and routed through an Nginx API gateway.

---

## Architecture

```
                        ┌──────────────────────────────────────────────┐
                        │              CLIENT (Browser / App)           │
                        └──────────────────┬───────────────────────────┘
                                           │ HTTP :80
                        ┌──────────────────▼───────────────────────────┐
                        │            NGINX API GATEWAY                  │
                        │                                               │
                        │  /api/auth/         → auth-service:8001       │
                        │  /api/spaces/       → space-service:8002      │
                        │  /api/billing/      → space-service:8002      │
                        │  /api/reports/      → space-service:8002      │
                        │  /api/reservations/ → reservation-service:8003│
                        └────────┬──────────────┬──────────────┬────────┘
                                 │              │              │
               ┌─────────────────▼──┐   ┌───────▼──────┐  ┌──▼────────────────┐
               │   auth-service     │   │ space-service │  │reservation-service│
               │   :8001            │   │ :8002         │  │:8003              │
               └─────────┬──────────┘   └──────┬────────┘  └────────┬──────────┘
                         │                     │                     │
               ┌─────────▼──────────┐ ┌────────▼────────┐ ┌─────────▼──────────┐
               │   auth-db          │ │   spaces-db      │ │  reservations-db   │
               │   postgres:15      │ │   postgres:15    │ │  postgres:15       │
               │   :5432            │ │   :5433          │ │  :5434             │
               └────────────────────┘ └──────────────────┘ └────────────────────┘
```

### Inter-service Communication

```
reservation-service ──[validates JWT]──► auth-service:8001/api/auth/validate
reservation-service ──[checks space ]──► space-service:8002/api/spaces/{id}
space-service       ──[validates JWT]──► auth-service:8001/api/auth/validate
```

---

## Prerequisites

| Tool | Minimum Version |
|------|----------------|
| Docker | 24.x |
| Docker Compose | v2.x (plugin) |
| `curl` | any (for testing) |

---

## Quick Start

```bash
# 1. Clone / enter the project directory
cd Api-maria

# 2. Build images and start all services
docker-compose up --build

# 3. To run in detached mode
docker-compose up --build -d

# 4. Tail logs for a specific service
docker-compose logs -f auth-service

# 5. Stop everything
docker-compose down

# 6. Stop and remove volumes (DESTRUCTIVE)
docker-compose down -v
```

---

## Services & Ports

| Service | Internal Port | Host Port | Description |
|---------|-------------|-----------|-------------|
| nginx-gateway | 80 | **80** | API Gateway — single entry point |
| auth-service | 8001 | 8001 | Authentication & JWT issuance |
| space-service | 8002 | 8002 | Spaces, billing & reports |
| reservation-service | 8003 | 8003 | Booking & reservation management |
| auth-db | 5432 | 5432 | PostgreSQL for auth-service |
| spaces-db | 5432 | 5433 | PostgreSQL for space-service |
| reservations-db | 5432 | 5434 | PostgreSQL for reservation-service |

---

## API Endpoints

All requests go through the gateway at `http://localhost:80`.

### 🔐 Auth Service — `/api/auth/`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/auth/register` | ❌ | Register a new user |
| POST | `/api/auth/login` | ❌ | Login and receive JWT |
| POST | `/api/auth/logout` | ✅ | Invalidate token |
| GET | `/api/auth/me` | ✅ | Get current user profile |
| POST | `/api/auth/refresh` | ✅ | Refresh JWT token |
| POST | `/api/auth/validate` | ✅ | Validate a JWT (inter-service) |

### 🏢 Space Service — `/api/spaces/`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/spaces/` | ❌ | List all available spaces |
| GET | `/api/spaces/{id}` | ❌ | Get space details |
| POST | `/api/spaces/` | ✅ Admin | Create a space |
| PUT | `/api/spaces/{id}` | ✅ Admin | Update a space |
| DELETE | `/api/spaces/{id}` | ✅ Admin | Delete a space |
| GET | `/api/spaces/{id}/availability` | ❌ | Check space availability |

### 💳 Billing — `/api/billing/`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/billing/invoices` | ✅ | List user invoices |
| GET | `/api/billing/invoices/{id}` | ✅ | Get invoice details |
| POST | `/api/billing/invoices` | ✅ | Create an invoice |
| PUT | `/api/billing/invoices/{id}/pay` | ✅ | Mark invoice as paid |

### 📊 Reports — `/api/reports/`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/reports/occupancy` | ✅ Admin | Occupancy report |
| GET | `/api/reports/revenue` | ✅ Admin | Revenue report |
| GET | `/api/reports/usage` | ✅ Admin | Space usage statistics |

### 📅 Reservation Service — `/api/reservations/`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/reservations/` | ✅ | List user reservations |
| POST | `/api/reservations/` | ✅ | Create a reservation |
| GET | `/api/reservations/{id}` | ✅ | Get reservation details |
| PUT | `/api/reservations/{id}` | ✅ | Update a reservation |
| DELETE | `/api/reservations/{id}` | ✅ | Cancel a reservation |
| GET | `/api/reservations/upcoming` | ✅ | List upcoming reservations |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | `coworking_jwt_secret_2024_secure_key` | HMAC secret for JWT signing |
| `AUTH_DB_HOST` | `auth-db` | Auth DB hostname |
| `AUTH_DB_PORT` | `5432` | Auth DB port |
| `AUTH_DB_USER` | `auth_user` | Auth DB username |
| `AUTH_DB_PASS` | `auth_pass` | Auth DB password |
| `AUTH_DB_NAME` | `auth_db` | Auth DB database name |
| `SPACES_DB_HOST` | `spaces-db` | Spaces DB hostname |
| `SPACES_DB_PORT` | `5432` | Spaces DB port |
| `SPACES_DB_USER` | `spaces_user` | Spaces DB username |
| `SPACES_DB_PASS` | `spaces_pass` | Spaces DB password |
| `SPACES_DB_NAME` | `spaces_db` | Spaces DB database name |
| `RES_DB_HOST` | `reservations-db` | Reservations DB hostname |
| `RES_DB_PORT` | `5432` | Reservations DB port |
| `RES_DB_USER` | `res_user` | Reservations DB username |
| `RES_DB_PASS` | `res_pass` | Reservations DB password |
| `RES_DB_NAME` | `reservations_db` | Reservations DB database name |
| `SPACE_SERVICE_URL` | `http://space-service:8002` | Internal URL for space-service |
| `AUTH_SERVICE_URL` | `http://auth-service:8001` | Internal URL for auth-service |
| `RESERVATION_SERVICE_URL` | `http://reservation-service:8003` | Internal URL for reservation-service |

---

## Rate Limiting

| Zone | Endpoint Prefix | Limit | Burst |
|------|----------------|-------|-------|
| `auth_zone` | `/api/auth/` | 10 req/s per IP | 20 |
| `api_zone` | All other `/api/*` | 30 req/s per IP | 60 |

Requests that exceed the rate limit receive HTTP **429** with a JSON body:
```json
{
  "error": "too_many_requests",
  "message": "Rate limit exceeded. Please slow down."
}
```

---

## Testing with curl

> Replace `<TOKEN>` with the JWT returned from the login endpoint.

### Register & Login

```bash
# Register a new user
curl -s -X POST http://localhost/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"Secret123!"}' | jq .

# Login and capture token
TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"Secret123!"}' | jq -r '.access_token')

echo "Token: $TOKEN"
```

### Auth — Profile & Token Validation

```bash
# Get current user profile
curl -s http://localhost/api/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq .

# Validate token (inter-service use)
curl -s -X POST http://localhost/api/auth/validate \
  -H "Authorization: Bearer $TOKEN" | jq .

# Refresh token
curl -s -X POST http://localhost/api/auth/refresh \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### Spaces

```bash
# List all spaces (public)
curl -s http://localhost/api/spaces/ | jq .

# Get a specific space
curl -s http://localhost/api/spaces/1 | jq .

# Check availability for a space
curl -s "http://localhost/api/spaces/1/availability?date=2024-07-15" | jq .

# Create a space (admin)
curl -s -X POST http://localhost/api/spaces/ \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Open Desk A1",
    "type": "desk",
    "capacity": 1,
    "price_per_hour": 5.00,
    "amenities": ["wifi","monitor","locker"]
  }' | jq .
```

### Reservations

```bash
# Create a reservation
curl -s -X POST http://localhost/api/reservations/ \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "space_id": 1,
    "start_time": "2024-07-15T09:00:00Z",
    "end_time":   "2024-07-15T17:00:00Z"
  }' | jq .

# List your reservations
curl -s http://localhost/api/reservations/ \
  -H "Authorization: Bearer $TOKEN" | jq .

# Upcoming reservations
curl -s http://localhost/api/reservations/upcoming \
  -H "Authorization: Bearer $TOKEN" | jq .

# Cancel a reservation
curl -s -X DELETE http://localhost/api/reservations/1 \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### Billing

```bash
# List invoices
curl -s http://localhost/api/billing/invoices \
  -H "Authorization: Bearer $TOKEN" | jq .

# Get invoice details
curl -s http://localhost/api/billing/invoices/1 \
  -H "Authorization: Bearer $TOKEN" | jq .

# Pay an invoice
curl -s -X PUT http://localhost/api/billing/invoices/1/pay \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"payment_method":"credit_card"}' | jq .
```

### Reports (Admin)

```bash
# Occupancy report
curl -s "http://localhost/api/reports/occupancy?from=2024-07-01&to=2024-07-31" \
  -H "Authorization: Bearer $TOKEN" | jq .

# Revenue report
curl -s "http://localhost/api/reports/revenue?from=2024-07-01&to=2024-07-31" \
  -H "Authorization: Bearer $TOKEN" | jq .

# Usage statistics
curl -s http://localhost/api/reports/usage \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### Gateway Health

```bash
# Nginx gateway health
curl -s http://localhost/health
```

---

## Project Structure

```
Api-maria/
├── .env                          # Shared environment variables
├── docker-compose.yml            # Service orchestration
├── README.md                     # This file
│
├── nginx/
│   ├── Dockerfile                # Builds nginx:alpine image
│   └── nginx.conf                # API gateway config (routing, CORS, rate limiting)
│
├── auth-service/
│   ├── Dockerfile
│   └── ...                       # Auth service source code
│
├── space-service/
│   ├── Dockerfile
│   └── ...                       # Space/billing/reports service source code
│
└── reservation-service/
    ├── Dockerfile
    └── ...                       # Reservation service source code
```

---

## Networking

All services communicate on the internal Docker network `coworking-network`. Services are addressable by their container name:

| From | To | Via |
|------|----|-----|
| nginx | auth-service | `http://auth-service:8001` |
| nginx | space-service | `http://space-service:8002` |
| nginx | reservation-service | `http://reservation-service:8003` |
| reservation-service | auth-service | `http://auth-service:8001` (env: `AUTH_SERVICE_URL`) |
| reservation-service | space-service | `http://space-service:8002` (env: `SPACE_SERVICE_URL`) |
| space-service | auth-service | `http://auth-service:8001` (env: `AUTH_SERVICE_URL`) |

---

## Health Checks & Start Order

Docker Compose uses `condition: service_healthy` to guarantee the boot sequence:

```
auth-db ──healthy──► auth-service ──healthy──► space-service ──healthy──► reservation-service
spaces-db ──healthy──►↑                                                           │
reservations-db ──healthy──────────────────────────────────────────────────────►↑
All three services healthy ──────────────────────────────────────────────────► nginx
```

---

## Security Notes

- Change `JWT_SECRET` to a cryptographically random string in production (e.g., `openssl rand -hex 32`).
- Do not commit `.env` to version control. Add it to `.gitignore`.
- Database ports `5432–5434` are exposed to localhost for development. Remove `ports:` from DB services in production.
- Consider adding TLS termination at the Nginx level (Let's Encrypt / cert-manager) for production deployments.
