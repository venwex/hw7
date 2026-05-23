FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api-service ./cmd/app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /src/api-service /app/api-service
COPY db/migrations /app/db/migrations

EXPOSE 8080

CMD ["/app/api-service"]
