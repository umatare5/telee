.PHONY: help build lint test-unit test-unit-coverage snapshot clean pre-commit-install pre-commit-test pre-commit-uninstall

BINARY_NAME := telee
BUILD_DIR := ./tmp
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)
COVERAGE_DIR := ./coverage

LDFLAGS := -X github.com/umatare5/telee/cli.version=$(shell cat VERSION)
BUILD_FLAGS := -trimpath -ldflags "$(LDFLAGS)"

.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  build                - Build the binary into $(BINARY_PATH)"
	@echo "  lint                 - Run golangci-lint and go mod tidy"
	@echo "  test-unit            - Run unit tests with coverage"
	@echo "  test-unit-coverage   - Generate the HTML coverage report"
	@echo "  snapshot             - Build a goreleaser snapshot"
	@echo "  clean                - Remove build artifacts and backup files"
	@echo "  pre-commit-install   - Install the pre-commit hooks"
	@echo "  pre-commit-test      - Run every hook across the whole tree"
	@echo "  pre-commit-uninstall - Remove the pre-commit hooks"
	@echo ""
	@echo "Requirements:"
	@echo "  - gotestsum:     go install gotest.tools/gotestsum@latest"
	@echo "  - golangci-lint: https://golangci-lint.run/docs/welcome/install/"
	@echo "  - pre-commit:    https://pre-commit.com/#install"
	@echo "  - gitleaks:      https://github.com/gitleaks/gitleaks#installing"

build:
	@mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BINARY_PATH) ./cmd

# config verify comes first because `run` accepts an unknown nested setting key
# silently and reverts that setting to its default, so a typo leaves a rule the
# author believes is on quietly off.
lint:
	golangci-lint config verify
	golangci-lint run
	go mod tidy

# The TELEE_* variables are cleared because they are read at flag-parse time, so a developer's own
# shell must not decide what the CLI tests see. A new variable the CLI reads has to be added here
# as well.
test-unit:
	@command -v gotestsum >/dev/null 2>&1 || { echo "Error: gotestsum is not installed. Run: go install gotest.tools/gotestsum@latest"; exit 1; }
	mkdir -p $(COVERAGE_DIR)
	env -u TELEE_HOSTNAME -u TELEE_COMMAND -u TELEE_USERNAME -u TELEE_PASSWORD -u TELEE_PRIVPASSWORD -u TELEE_HOSTKEYPATH \
		gotestsum --format testname -- -race -coverprofile=$(COVERAGE_DIR)/report.out ./...

test-unit-coverage: test-unit
	go tool cover -html=$(COVERAGE_DIR)/report.out -o $(COVERAGE_DIR)/report.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/report.html"

snapshot:
	goreleaser release --snapshot --clean

# BUILD_DIR is ./tmp, which also holds the worktrees `git wt` creates.
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
