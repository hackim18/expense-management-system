# Expense Management System - Backend

Go backend service built with Gin + Gorm using clean architecture style.

## Tech Stack
- Go (see `go.mod`)
- Gin (HTTP)
- Gorm (PostgreSQL)
- Viper (config)
- Logrus (logging)
- Swaggo (Swagger UI)

## Local Setup
1) Copy env file:
```bash
cp .env.example .env
```

2) Run migrations + seed + server:
```bash
go run ./cmd/web --migrate --seed --run
```

Default seed users:
- Employee: `john@mail.com` / `12345678`
- Manager: `manager@mail.com` / `12345678`

## Docker
Use the root `docker-compose.yml` to run full stack:
```bash
docker compose up --build
```

## Environment Variables
- `APP_NAME`, `PORT`, `LOG_LEVEL`
- `DB_USERNAME`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME`
- `JWT_SECRET`, `JWT_ISSUER`, `JWT_AUDIENCE`, `JWT_EXPIRES_MINUTES`
- `PAYMENT_BASE_URL`, `PAYMENT_TIMEOUT_SECONDS`, `PAYMENT_RETRY_COUNT`, `PAYMENT_RETRY_DELAY_SECONDS`, `PAYMENT_QUEUE_BUFFER`
- `DROP_TABLE_NAMES` (comma separated)
- `CORS_ALLOW_ORIGINS` (comma separated), `CORS_ALLOW_CREDENTIALS`
- `RATE_LIMIT` (format like `100-M`)
- `RATE_LIMIT_EXCLUDE_PATHS` (comma separated; supports `/*` suffix)
- `SMTP_ENABLED`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`
- `SMTP_FROM_EMAIL`, `SMTP_FROM_NAME`

## API Endpoints
- `POST /api/auth/login`
- `POST /api/auth/register` (helper for local usage)
- `POST /api/expenses` (auth)
- `GET /api/expenses` (auth, supports `status`, `page`, `size`)
- `GET /api/expenses/:id` (auth)
- `GET /api/expenses/:id/history` (auth)
- `PUT /api/expenses/:id/approve` (auth, manager only)
- `PUT /api/expenses/:id/reject` (auth, manager only)
- `GET /api/health`
- `GET /api/metrics`

## Business Rules
- Currency is IDR only; amount is stored as integer.
- Min amount: IDR 10,000; max amount: IDR 50,000,000.
- Expenses >= IDR 1,000,000 require manager approval.
- Expenses < IDR 1,000,000 are auto-approved.
- Approved expenses trigger background payment processing.
- Payment success updates status to `completed`.
- Status values: `awaiting_approval`, `auto_approved`, `approved`, `rejected`, `completed`.
- Every expense status change is recorded in `expense_status_histories`.
- Expenses that require approval trigger an email notification to manager accounts (SMTP configurable).

## Approval & Payment Flow
- Approve endpoint sets status to `approved` when the current status is `awaiting_approval`.
- After approval, a payment job is enqueued. The background worker can process immediately, so a follow-up GET may show `completed` quickly if the payment mock succeeds.
- The approve response itself is generated before payment finishes, so it should still return `approved` at the time of response.

## Payment Processor Mock
- Base URL: `https://1620e98f-7759-431c-a2aa-f449d591150b.mock.pstmn.io`
- Endpoint: `POST /v1/payments`
- Request body: `amount` (int), `external_id` (idempotency key)
- Handles idempotency response when external_id already exists.

## Database Seeding
Seed files:
- `internal/migrations/json/users.json`
- `internal/migrations/json/expenses.json`
- `internal/migrations/json/approvals.json`
- `internal/migrations/json/expense_status_histories.json`

## Testing
```bash
go test ./...
```

## Swagger/OpenAPI
- Spec file: `api/openapi.yaml`
- Endpoint: `GET /api/openapi.yaml`
- Swagger UI: `GET /swagger/index.html`

## Architecture Notes
- Clean architecture style with separation of delivery, usecase, repository, and entity layers.
- External services (payment, email) are injected via interfaces for testability.
- Uses structured logging and centralized error handling.

## Assumptions
- Payment processing is mocked and considered successful when mock returns HTTP 200 or idempotent 400.
- Email notifications are best-effort; failures are logged but do not block requests.

## Improvements
- Add more business logic unit tests for usecases.
- Add integration tests for API flows (approve and payment).
- Add background job visibility (queue size, processing metrics).
