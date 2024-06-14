# Define variables
BINARY_NAME := terraform-provider-telnyx
BINARY_VERSION := dev

# Directories
SRC_DIR := .
BUILD_DIR := ./bin
DOCS_DIR := ./docs

# Go commands
GO := go
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean
GOTEST := $(GO) test
GOFMT := $(GO) fmt

# All command
.PHONY: all
all: format build test

# Format code
.PHONY: format
format:
	$(GOFMT) $(SRC_DIR)/...

# Build binary
.PHONY: build
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)_$(BINARY_VERSION)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Generate docs
.PHONY: docs
docs:
	$(GO) run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name telnyx --provider-dir $(SRC_DIR) --rendered-website-dir $(DOCS_DIR)

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Help
.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  format    - Format the Go code."
	@echo "  build     - Build the binary."
	@echo "  test      - Run tests."
	@echo "  docs      - Generate documentation."
	@echo "  clean     - Clean build artifacts."
	@echo "  all       - Format, build, and test."
