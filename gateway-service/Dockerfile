FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway-service ./cmd/app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /src/gateway-service /app/gateway-service

EXPOSE 8081

CMD ["/app/gateway-service"]
