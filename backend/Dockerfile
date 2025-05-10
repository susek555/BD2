FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/car-dealer-api ./cmd/car-dealer-api

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates libc6-compat
COPY --from=builder /app/car-dealer-api .
COPY --from=builder /app/.env .
COPY --from=builder /app/data ./data

EXPOSE 8080

CMD ["/app/car-dealer-api"]