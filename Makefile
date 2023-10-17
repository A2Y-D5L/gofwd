# Variables
GO := go
GOFLAGS :=
VERBOSE ?= @
OUTPUT_DIR := _output
BINARY_NAME := gofwd
VERSION := $(shell git describe --always --long --dirty)
GOLANGCI_LINT := $(shell which golangci-lint)

# Default target executed when no arguments are given to make.
default: build

# Install necessary tools
tools:
ifndef GOLANGCI_LINT
	$(VERBOSE)curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
endif

# Build the binary.
build:
	$(VERBOSE)mkdir -p $(OUTPUT_DIR)
	$(VERBOSE)$(GO) build $(GOFLAGS) -ldflags "-X main.Version=$(VERSION)" -o $(OUTPUT_DIR)/$(BINARY_NAME) -v

# Run tests.
test:
	$(VERBOSE)$(GO) test -v ./...

# Clean the binary.
clean:
	$(VERBOSE)$(GO) clean
	$(VERBOSE)rm -f $(OUTPUT_DIR)/$(BINARY_NAME)

# Run the program.
run: build
	$(VERBOSE)./$(OUTPUT_DIR)/$(BINARY_NAME)

# Install missing dependencies.
deps:
	$(VERBOSE)$(GO) get ./...

# Update dependencies
update-deps:
	$(VERBOSE)$(GO) get -u ./...

# Tidy up and clean module dependencies
tidy:
	$(VERBOSE)$(GO) mod tidy

# Lint the code
lint: tools
	$(VERBOSE)golangci-lint run

# Run static analysis on the code
vet:
	$(VERBOSE)$(GO) vet ./...

# Cross-compilation
build-linux:
	$(VERBOSE)GOOS=linux GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux -v

build-windows:
	$(VERBOSE)GOOS=windows GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME).exe -v

# Build and run the program.
all: build run

# Help information.
help:
	@echo "Use: make <target> where <target> is one of"
	@echo "  build        to build the gofwd binary"
	@echo "  test         to run tests"
	@echo "  clean        to remove the binary and temporary files"
	@echo "  run          to build and run gofwd"
	@echo "  deps         to fetch missing dependencies"
	@echo "  update-deps  to update dependencies"
	@echo "  tidy         to tidy up and clean module dependencies"
	@echo "  lint         to lint the code using golangci-lint (will install if not present)"
	@echo "  vet          to run static analysis on the code"
	@echo "  build-linux  to cross-compile for Linux"
	@echo "  build-windows to cross-compile for Windows"
	@echo "  all          to build and run the program"
	@echo "  help         to show this help information"

.PHONY: default build test clean run deps update-deps tidy lint vet build-linux build-windows all help tools