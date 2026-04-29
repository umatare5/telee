.PHONY: help build lint clean

# Binary name and paths
BINARY_NAME := telee
BUILD_DIR := ./tmp
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)

# Go build flags
LDFLAGS := -X github.com/umatare5/telee/cli.version=$(shell cat VERSION)
BUILD_FLAGS := -ldflags "$(LDFLAGS)"

# Default target
.DEFAULT_GOAL := help

# Show available targets
help:
	@echo "Available targets:"
	@echo "  build              - Build the binary"
	@echo "  lint               - Run linters (golangci-lint)"
	@echo "  clean              - Remove build artifacts and backup files"
	@echo ""
	@echo "Requirements:"
	@echo "  - golangci-lint: https://golangci-lint.run/usage/install/"

build: $(BINARY_PATH)

# Build the binary
$(BINARY_PATH):
	mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BINARY_PATH) ./cmd

# Lint the code
lint:
	golangci-lint run
	go mod tidy

# Clean build artifacts and backup files
clean:
	rm -rf $(BUILD_DIR)
	find . -name "*.bak*" -type f -delete 2>/dev/null || true
