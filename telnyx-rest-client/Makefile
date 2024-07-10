# Directories
MAIN_DIR := cmd
BIN_DIR := bin

# Binary name
BINARY := telnyx-client

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
	TELNYX_API_TOKEN="w/e" TF_ACC=1 $(GOTEST) ./internal/provider -v --count=1

format:
	$(GOFMT) ./...

clean:
	rm -f $(BIN_DIR)/$(BINARY)
