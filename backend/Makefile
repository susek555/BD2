build:
	go build -o bin/main cmd/car-dealer-api/main.go

run:
	go run cmd/car-dealer-api/main.go

live-server:
	 CompileDaemon  -build="make build"  -command="./bin/main" -color=true


docs:
	swag init --parseDependency --parseInternal -g cmd/car-dealer-api/main.go