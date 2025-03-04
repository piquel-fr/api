run: build
	@./bin/main

build: sql $(wildcard *.go)
	@echo Starting main build...
	@go build -o bin/main main.go

sql: $(wildcard *.sql)
	@sqlc generate
	@echo Generated SQLc files...
