build:
	go build -o bin/main cmd/car-dealer-api/main.go

run:
	go run cmd/car-dealer-api/main.go

live-server:
	CompileDaemon  -build="make build"  -command="./bin/main" -color=true

create-docs:
	swag init -g ./cmd/car-dealer-api/main.go

reload-database:
	rm data/test.db
	touch data/test.db

test-unit:
	go test ./internal/test/... -tags unit

test-integration:
	go test ./internal/test/... -tags integration