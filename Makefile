SPEC_DIR=./api/specification
SPEC_OUT_DIR=./api
SPECS=$(shell find $(SPEC_DIR) -name '*.json')
OUT_SPECS=$(SPECS:$(SPEC_DIR)/%.json=$(SPEC_OUT_DIR)/%.gen.go)
SPEC_GEN_CONFIG=$(SPEC_DIR)/oapi-codegen.yml

.PHONY: run
run: build
	@./bin/main

.PHONY: build
build: sql $(wildcard *.go) $(OUT_SPECS)
	@echo Starting main build...
	@go build -o bin/main main.go

.PHONY: sql
sql: $(wildcard *.sql)
	@sqlc generate
	@echo Generated SQLc files...

$(SPEC_OUT_DIR)/%.gen.go: $(SPEC_DIR)/%.json
	@echo Generating spec for $*...
	@oapi-codegen --config $(SPEC_GEN_CONFIG) -o $@ $<
