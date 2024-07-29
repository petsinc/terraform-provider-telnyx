# Directories
MAIN_DIR := .
BIN_DIR := ./bin
DOCS_DIR := ../docs
PROVIDER_DIR := ./provider

# Binary name
BINARY := terraform-provider-telnyx

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOTEST := $(GOCMD) test
GOFMT := $(GOCMD) fmt

# Targets
.PHONY: all build run test format clean docs help

# All command
all: format build test

# Install dependencies
deps:
	cd $(PROVIDER_DIR) && $(GOCMD) mod tidy && $(GOCMD) mod download
	
# Build binary with version
build:
	cd $(PROVIDER_DIR) && $(GOBUILD) -ldflags "-X main.version=$(VERSION)" -o $(BIN_DIR)/$(BINARY) ./main.go

# Run the program
run:
	cd $(PROVIDER_DIR) && $(GORUN) main.go

# Run tests
test:
	cd $(PROVIDER_DIR) && . .env && TELNYX_API_KEY=$$TELNYX_API_KEY TF_ACC=1 $(GOTEST) ./internal/provider -v --count=1

# Format code
format:
	cd $(PROVIDER_DIR) && $(GOFMT) ./...

# Generate docs
docs:
	cd $(PROVIDER_DIR) && $(GOCMD) run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name telnyx --provider-dir . --rendered-website-dir $(DOCS_DIR)

# Clean build artifacts
clean:
	rm -rf $(PROVIDER_DIR)/$(BIN_DIR)

# Help command
help:
	@echo "Makefile commands:"
	@echo "  all       - Format, build, and test."
	@echo "  format    - Format the Go code."
	@echo "  build     - Build the binary."
	@echo "  run       - Run the program."
	@echo "  test      - Run tests."
	@echo "  docs      - Generate documentation."
	@echo "  clean     - Clean build artifacts."
