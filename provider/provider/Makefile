# Directories
MAIN_DIR := .
BIN_DIR := bin

# Binary name
BINARY := provider

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOTEST := $(GOCMD) test
GOFMT := $(GOCMD) fmt

# Targets
.PHONY: build run test format clean

build:
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY) $(MAIN_DIR)/main.go

run:
	$(GORUN) $(MAIN_DIR)/main.go

test:
	@. .env && TELNYX_API_TOKEN=$$TELNYX_API_TOKEN TF_ACC=1 $(GOTEST) ./internal/provider -v --count=1

format:
	$(GOFMT) ./...
	goimports -w .

clean:
	rm -f $(BIN_DIR)/$(BINARY)
