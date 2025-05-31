# Makefile for Go DWG Extractor
# Cross-platform build support for Windows, Linux, and macOS

# Configuration
PROJECT_NAME := go-dwg-extractor
DIST_DIR := dist
VERSION ?= v0.1.0-dev
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS := -X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)

# Platform configurations
PLATFORMS := windows/amd64 linux/amd64 darwin/amd64 darwin/arm64

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := all

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Go DWG Extractor Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Variables:"
	@echo "  $(BLUE)VERSION$(NC)     Set build version (default: $(VERSION))"
	@echo "  $(BLUE)GIT_COMMIT$(NC)  Set git commit hash (auto-detected: $(GIT_COMMIT))"
	@echo ""
	@echo "Examples:"
	@echo "  make all VERSION=v1.0.0"
	@echo "  make build-windows"
	@echo "  make test"
	@echo "  make clean"

# Check prerequisites
.PHONY: check
check: ## Check build prerequisites
	@echo "$(BLUE)[INFO]$(NC) Checking prerequisites..."
	@command -v go >/dev/null 2>&1 || { echo "$(RED)[ERROR]$(NC) Go is not installed or not in PATH"; exit 1; }
	@command -v git >/dev/null 2>&1 || echo "$(YELLOW)[WARNING]$(NC) Git is not installed - version info may be incomplete"
	@echo "$(GREEN)[SUCCESS]$(NC) Prerequisites check completed"
	@echo "$(BLUE)[INFO]$(NC) Go version: $$(go version)"
	@echo "$(BLUE)[INFO]$(NC) Build version: $(VERSION)"
	@echo "$(BLUE)[INFO]$(NC) Git commit: $(GIT_COMMIT)"

# Clean build artifacts
.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)[INFO]$(NC) Cleaning build artifacts..."
	@rm -rf $(DIST_DIR)
	@echo "$(GREEN)[SUCCESS]$(NC) Clean completed"

# Run tests
.PHONY: test
test: ## Run all tests
	@echo "$(BLUE)[INFO]$(NC) Running tests..."
	@go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)[INFO]$(NC) Running tests with coverage..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)[SUCCESS]$(NC) Coverage report generated: coverage.html"

# Format code
.PHONY: fmt
fmt: ## Format Go code
	@echo "$(BLUE)[INFO]$(NC) Formatting code..."
	@go fmt ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Code formatting completed"

# Lint code
.PHONY: lint
lint: ## Lint Go code
	@echo "$(BLUE)[INFO]$(NC) Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)[WARNING]$(NC) golangci-lint not found, running go vet instead"; \
		go vet ./...; \
	fi
	@echo "$(GREEN)[SUCCESS]$(NC) Linting completed"

# Build for current platform
.PHONY: build
build: check ## Build for current platform
	@echo "$(BLUE)[INFO]$(NC) Building for current platform..."
	@mkdir -p $(DIST_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(PROJECT_NAME) .
	@echo "$(GREEN)[SUCCESS]$(NC) Built $(DIST_DIR)/$(PROJECT_NAME)"

# Platform-specific build targets
.PHONY: build-windows
build-windows: check ## Build for Windows (amd64)
	@echo "$(BLUE)[INFO]$(NC) Building for windows/amd64..."
	@mkdir -p $(DIST_DIR)/windows-amd64
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/windows-amd64/$(PROJECT_NAME).exe .
	@if [ -f README.md ]; then cp README.md $(DIST_DIR)/windows-amd64/; fi
	@if [ -f LICENSE ]; then cp LICENSE $(DIST_DIR)/windows-amd64/; fi
	@echo "$(GREEN)[SUCCESS]$(NC) Built $(DIST_DIR)/windows-amd64/$(PROJECT_NAME).exe"

.PHONY: build-linux
build-linux: check ## Build for Linux (amd64)
	@echo "$(BLUE)[INFO]$(NC) Building for linux/amd64..."
	@mkdir -p $(DIST_DIR)/linux-amd64
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/linux-amd64/$(PROJECT_NAME) .
	@if [ -f README.md ]; then cp README.md $(DIST_DIR)/linux-amd64/; fi
	@if [ -f LICENSE ]; then cp LICENSE $(DIST_DIR)/linux-amd64/; fi
	@echo "$(GREEN)[SUCCESS]$(NC) Built $(DIST_DIR)/linux-amd64/$(PROJECT_NAME)"

.PHONY: build-darwin
build-darwin: check ## Build for macOS (amd64)
	@echo "$(BLUE)[INFO]$(NC) Building for darwin/amd64..."
	@mkdir -p $(DIST_DIR)/darwin-amd64
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/darwin-amd64/$(PROJECT_NAME) .
	@if [ -f README.md ]; then cp README.md $(DIST_DIR)/darwin-amd64/; fi
	@if [ -f LICENSE ]; then cp LICENSE $(DIST_DIR)/darwin-amd64/; fi
	@echo "$(GREEN)[SUCCESS]$(NC) Built $(DIST_DIR)/darwin-amd64/$(PROJECT_NAME)"

.PHONY: build-darwin-arm64
build-darwin-arm64: check ## Build for macOS (arm64)
	@echo "$(BLUE)[INFO]$(NC) Building for darwin/arm64..."
	@mkdir -p $(DIST_DIR)/darwin-arm64
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/darwin-arm64/$(PROJECT_NAME) .
	@if [ -f README.md ]; then cp README.md $(DIST_DIR)/darwin-arm64/; fi
	@if [ -f LICENSE ]; then cp LICENSE $(DIST_DIR)/darwin-arm64/; fi
	@echo "$(GREEN)[SUCCESS]$(NC) Built $(DIST_DIR)/darwin-arm64/$(PROJECT_NAME)"

# Build all platforms
.PHONY: build-all
build-all: check build-windows build-linux build-darwin build-darwin-arm64 ## Build for all platforms
	@echo "$(GREEN)[SUCCESS]$(NC) All platform builds completed!"

# Create distribution packages
.PHONY: package
package: ## Create distribution packages
	@echo "$(BLUE)[INFO]$(NC) Creating distribution packages..."
	@if [ -d "$(DIST_DIR)/windows-amd64" ]; then \
		cd $(DIST_DIR)/windows-amd64 && \
		if command -v zip >/dev/null 2>&1; then \
			zip -r ../$(PROJECT_NAME)-windows-amd64.zip .; \
			echo "$(GREEN)[SUCCESS]$(NC) Created $(PROJECT_NAME)-windows-amd64.zip"; \
		else \
			echo "$(YELLOW)[WARNING]$(NC) zip command not found, skipping Windows package"; \
		fi; \
	fi
	@if [ -d "$(DIST_DIR)/linux-amd64" ]; then \
		cd $(DIST_DIR)/linux-amd64 && \
		tar -czf ../$(PROJECT_NAME)-linux-amd64.tar.gz . && \
		echo "$(GREEN)[SUCCESS]$(NC) Created $(PROJECT_NAME)-linux-amd64.tar.gz"; \
	fi
	@if [ -d "$(DIST_DIR)/darwin-amd64" ]; then \
		cd $(DIST_DIR)/darwin-amd64 && \
		if command -v zip >/dev/null 2>&1; then \
			zip -r ../$(PROJECT_NAME)-darwin-amd64.zip .; \
			echo "$(GREEN)[SUCCESS]$(NC) Created $(PROJECT_NAME)-darwin-amd64.zip"; \
		else \
			tar -czf ../$(PROJECT_NAME)-darwin-amd64.tar.gz . && \
			echo "$(GREEN)[SUCCESS]$(NC) Created $(PROJECT_NAME)-darwin-amd64.tar.gz"; \
		fi; \
	fi
	@if [ -d "$(DIST_DIR)/darwin-arm64" ]; then \
		cd $(DIST_DIR)/darwin-arm64 && \
		if command -v zip >/dev/null 2>&1; then \
			zip -r ../$(PROJECT_NAME)-darwin-arm64.zip .; \
			echo "$(GREEN)[SUCCESS]$(NC) Created $(PROJECT_NAME)-darwin-arm64.zip"; \
		else \
			tar -czf ../$(PROJECT_NAME)-darwin-arm64.tar.gz . && \
			echo "$(GREEN)[SUCCESS]$(NC) Created $(PROJECT_NAME)-darwin-arm64.tar.gz"; \
		fi; \
	fi

# Development targets
.PHONY: dev
dev: fmt lint test build ## Development workflow (format, lint, test, build)

.PHONY: ci
ci: check fmt lint test-coverage build-all ## CI workflow (check, format, lint, test with coverage, build all)

# Install to local bin
.PHONY: install
install: build ## Install to local bin directory
	@echo "$(BLUE)[INFO]$(NC) Installing to local bin..."
	@mkdir -p ~/bin
	@cp $(DIST_DIR)/$(PROJECT_NAME) ~/bin/
	@echo "$(GREEN)[SUCCESS]$(NC) Installed to ~/bin/$(PROJECT_NAME)"
	@echo "$(YELLOW)[NOTE]$(NC) Make sure ~/bin is in your PATH"

# Show build summary
.PHONY: summary
summary: ## Show build summary
	@echo "$(BLUE)[INFO]$(NC) Build Summary:"
	@echo ""
	@if [ -d "$(DIST_DIR)" ]; then \
		echo "Generated files:"; \
		find $(DIST_DIR) -type f -name "$(PROJECT_NAME)*" -exec sh -c 'printf "  %s (%s)\n" "$$1" "$$(du -h "$$1" | cut -f1)"' _ {} \; 2>/dev/null || true; \
		echo ""; \
		echo "Distribution packages:"; \
		find $(DIST_DIR) -type f \( -name "*.zip" -o -name "*.tar.gz" \) -exec sh -c 'printf "  %s (%s)\n" "$$1" "$$(du -h "$$1" | cut -f1)"' _ {} \; 2>/dev/null || true; \
	else \
		echo "$(YELLOW)[WARNING]$(NC) No build output found"; \
	fi

# Complete build pipeline
.PHONY: all
all: clean build-all package summary ## Complete build pipeline (clean, build all platforms, package, summary)

# Version information
.PHONY: version
version: ## Show version information
	@echo "Project: $(PROJECT_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)" 