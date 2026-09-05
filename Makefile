.PHONY: help build lint test-unit test-unit-coverage clean pre-commit-install pre-commit-test pre-commit-uninstall

# Binary name and paths
BINARY_NAME := telee
BUILD_DIR := ./tmp
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)
COVERAGE_DIR := ./coverage

# Go build flags
LDFLAGS := -X github.com/umatare5/telee/cli.version=$(shell cat VERSION)
BUILD_FLAGS := -ldflags "$(LDFLAGS)"

# Default target
.DEFAULT_GOAL := help

# Show available targets
help:
	@echo "Available targets:"
	@echo "  build                - Build the binary"
	@echo "  lint                 - Run linters (golangci-lint)"
	@echo "  test-unit            - Run unit tests with colored output"
	@echo "  test-unit-coverage   - Generate HTML coverage report"
	@echo "  clean                - Remove build artifacts and backup files"
	@echo "  pre-commit-install   - Install the pre-commit hooks"
	@echo "  pre-commit-test      - Run every hook across the whole tree"
	@echo "  pre-commit-uninstall - Remove the pre-commit hooks"
	@echo ""
	@echo "Requirements:"
	@echo "  - gotestsum: go install gotest.tools/gotestsum@latest"
	@echo "  - golangci-lint: https://golangci-lint.run/docs/welcome/install/"
	@echo "  - pre-commit: https://pre-commit.com/#install"
	@echo "  - gitleaks: https://github.com/gitleaks/gitleaks#installing"

# Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BINARY_PATH) ./cmd

# Lint the code
lint:
	golangci-lint run
	go mod tidy

# Run unit tests with gotestsum (shows individual test results with color)
test-unit:
	@command -v gotestsum >/dev/null 2>&1 || { echo "Error: gotestsum is not installed. Run: go install gotest.tools/gotestsum@latest"; exit 1; }
	mkdir -p $(COVERAGE_DIR)
	gotestsum --format testname -- -coverprofile=$(COVERAGE_DIR)/report.out ./...

# Generate coverage report (HTML)
test-unit-coverage: test-unit
	go tool cover -html=$(COVERAGE_DIR)/report.out -o $(COVERAGE_DIR)/report.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/report.html"

# BUILD_DIR is ./tmp, which also holds the worktrees `git wt` creates, so this removes
# what the build produced rather than the directory itself.
clean:
	rm -f $(BINARY_PATH)
	rm -rf $(BUILD_DIR)/dist
	rm -f $(COVERAGE_DIR)/report.html
	find . -name "*.bak*" -type f -delete 2>/dev/null || true

# --allow-missing-config is load-bearing: the hook path is the shared git common
# dir, so the hook installed from one worktree also runs on every other one and
# on main, where .pre-commit-config.yaml may not exist yet.
pre-commit-install:
	@command -v pre-commit >/dev/null 2>&1 || { echo "Error: pre-commit is not installed. See: https://pre-commit.com/#install"; exit 1; }
	@pre-commit install --allow-missing-config

pre-commit-test:
	@pre-commit run --all-files

pre-commit-uninstall:
	@pre-commit uninstall
