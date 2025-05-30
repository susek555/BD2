FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git build-base
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN rm -rf ./docs

# Generate Swagger documentation
RUN swag init --parseDependency --parseInternal -g ./cmd/car-dealer-api/main.go -o ./docs
# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/car-dealer-api

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates libc6-compat curl

COPY --from=builder /app/main .

RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup
RUN chown -R appuser:appgroup /root

USER appuser
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["./main"]