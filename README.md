# Weather API: Docker + CI/CD + split into services

Project is split into two Go services:

- `api-service` — REST API, business logic, PostgreSQL, repository/service layers, DI, logging middleware, tests.
- `gateway-service` — external API gateway. It calls Open-Meteo and has no direct DB access.

Runtime flow:

```text
Client -> api-service:8080 -> gateway-service:8081 -> Open-Meteo public API
                         -> PostgreSQL
```

## Run

```bash
docker compose up --build
```

Health checks:

```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
```

## Demo requests

Register user:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alnur","email":"alnur@example.com","password":"password123"}'
```

Login:

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alnur@example.com","password":"password123"}' | sed -E 's/.*"token":"([^"]+)".*/\1/')
```

Add city:

```bash
curl -X POST http://localhost:8080/cities \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"city":"Almaty"}'
```

Get weather. API service calls gateway-service, gateway-service calls Open-Meteo, API service saves result into PostgreSQL:

```bash
curl http://localhost:8080/weather \
  -H "Authorization: Bearer $TOKEN"
```

Get history:

```bash
curl "http://localhost:8080/weather/history?city=Almaty" \
  -H "Authorization: Bearer $TOKEN"
```

Direct gateway check:

```bash
curl "http://localhost:8081/weather?city=Almaty"
```

## Tests

Unit tests:

```bash
cd api-service
go test ./internal/handler ./internal/service
```

Repository integration test requires PostgreSQL. Example:

```bash
docker run --name weather-test-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=weather_test \
  -p 5437:5432 \
  -d postgres:16-alpine

cd api-service
TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5437/weather_test?sslmode=disable" go test ./internal/repository
```

## CI/CD

GitHub Actions workflow is placed in `.github/workflows/ci.yml`.
It runs:

- API service tests;
- gateway service tests;
- Docker Compose image build.
