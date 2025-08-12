run: build
	@./bin/main

.PHONY: build
build: sql $(wildcard *.go)
	@echo Starting main build...
	@go build -o bin/main main.go

.PHONY: sql
sql: $(wildcard *.sql)
	@sqlc generate
	@echo Generated SQLc files...
