run: build
	@./bin/main

build:
	@sqlc generate
	@go build -o bin/main main.go
