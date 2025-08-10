# TigerGraph CLI Makefile

# Variables
BINARY_NAME=tg
MODULE_NAME=github.com/zrougamed/tgCli
VERSION=0.1.1
BUILD_DIR=build
MAIN_FILE=cmd/main.go
GO_FILES=$(shell find . -name "*.go" -type f)
TEST_FILES=$(shell find . -name "*_test.go" -type f)

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
PURPLE=\033[0;35m
CYAN=\033[0;36m
WHITE=\033[1;37m
NC=\033[0m # No Color

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "$(CYAN)TigerGraph CLI Build System$(NC)"
	@echo "$(WHITE)Usage: make [target]$(NC)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@awk -F ':[^#]*## ' \
		-v GREEN="$(GREEN)" -v NC="$(NC)" \
		'/^[a-zA-Z0-9_-]+:.*##/ {printf "  %s%-20s%s %s\n", GREEN, $$1, NC, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  $(WHITE)make build$(NC)           # Build for current platform"
	@echo "  $(WHITE)make test$(NC)            # Run all tests"
	@echo "  $(WHITE)make test-short$(NC)      # Run all tests in short mode (skip interactive tests)"
	@echo "  $(WHITE)make test-coverage$(NC)   # Run tests with coverage"
	@echo "  $(WHITE)make build-all$(NC)       # Build for all platforms"
	@echo "  $(WHITE)make release$(NC)         # Prepare release packages"
	@echo "  $(WHITE)make dev$(NC)             # Quick development build and test"
	@echo ""
	@echo "$(YELLOW)Project Info:$(NC)"
	@echo "  Binary Name: $(PURPLE)$(BINARY_NAME)$(NC)"
	@echo "  Version:     $(PURPLE)$(VERSION)$(NC)"
	@echo "  Go Files:    $(PURPLE)$(words $(GO_FILES))$(NC)"
	@echo "  Test Files:  $(PURPLE)$(words $(TEST_FILES))$(NC)"


# Clean build artifacts
.PHONY: clean
clean: ## Clean build artifacts and test cache
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME).exe
	@rm -f $(BINARY_NAME)-*
	@rm -f coverage.out coverage.html
	@rm -f test_output.log test_report.json
	@rm -f cpu.prof mem.prof
	@go clean -testcache
	@echo "✓ Clean completed"

# Initialize Go modules
.PHONY: mod
mod: ## Download and tidy Go modules
	@echo "Downloading Go modules..."
	@go mod download
	@go mod tidy
	@echo "✓ Modules updated"

# Format Go code
.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "✓ Code formatted"

# Run tests
.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...
	@echo "✓ Tests completed"

.PHONY: test-short
test-short: ## Run tests in short mode (skip interactive tests)
	@echo "Running tests in short mode..."
	@go test -v -short ./...
	@echo "✓ Short tests completed"

# Run tests with coverage
.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Run tests with race detection
.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@go test -v -race ./...
	@echo "✓ Race tests completed"

# Run benchmarks
.PHONY: bench
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Run specific package tests
.PHONY: test-pkg
test-pkg: ## Run tests for specific package (usage: make test-pkg PKG=./internal/cloud)
	@echo "Running tests for package: $(PKG)"
	@go test -v $(PKG)


# Lint code
.PHONY: lint
lint: ## Lint Go code (requires golangci-lint)
	@echo "Linting Go code..."
	@golangci-lint run || echo "golangci-lint not installed, skipping lint"

# Build for current platform
.PHONY: build
build: mod fmt ## Build binary for current platform
	@echo "$(YELLOW)Building $(BINARY_NAME) for current platform...$(NC)"
	@go build -ldflags "-X '$(MODULE_NAME)/pkg/constants.VERSION_CLI=$(VERSION)'" -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "$(GREEN)✓ Build completed: $(BINARY_NAME)$(NC)"

# Install binary
.PHONY: install
install: build ## Install binary to GOPATH/bin
	@echo "$(YELLOW)Installing $(BINARY_NAME)...$(NC)"
	@go install -ldflags "-X '$(MODULE_NAME)/pkg/constants.VERSION_CLI=$(VERSION)'" $(MAIN_FILE)
	@echo "$(GREEN)✓ Installation completed$(NC)"

# Build for all platforms
.PHONY: build-all
build-all: mod fmt ## Build binaries for all platforms
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	
	@echo "Building for Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X '$(MODULE_NAME)/pkg/constants.VERSION_CLI=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	
	@echo "Building for Linux ARM64..."
	@GOOS=linux GOARCH=arm64 go build -ldflags "-X '$(MODULE_NAME)/pkg/constants.VERSION_CLI=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)
	
	@echo "Building for macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X '$(MODULE_NAME)/pkg/constants.VERSION_CLI=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	
	@echo "Building for macOS ARM64 (M1)..."
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-X '$(MODULE_NAME)/pkg/constants.VERSION_CLI=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	
	@echo "Building for Windows AMD64..."
	@GOOS=windows GOARCH=amd64 go build -ldflags "-X '$(MODULE_NAME)/pkg/constants.VERSION_CLI=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	
	@echo "Building for Windows ARM64..."
	@GOOS=windows GOARCH=arm64 go build -ldflags "-X '$(MODULE_NAME)/pkg/constants.VERSION_CLI=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_FILE)

# Create release packages
.PHONY: package
package: build-all ## Create release packages
	@echo "Creating release packages..."
	@mkdir -p $(BUILD_DIR)/packages
	
	# Linux AMD64
	@tar -czf $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-amd64
	
	# Linux ARM64
	@tar -czf $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-arm64
	
	# macOS AMD64
	@tar -czf $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-amd64
	
	# macOS ARM64
	@tar -czf $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-arm64
	
	# Windows AMD64
	@zip -j $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	
	# Windows ARM64
	@zip -j $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-windows-arm64.zip $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe
	
	@echo "Release packages created in $(BUILD_DIR)/packages/"

# Quick development build and test
.PHONY: dev
dev: build test-short ## Quick development build and test
	@echo "$(YELLOW)Running version check...$(NC)"
	@./$(BINARY_NAME) version
	@echo "$(GREEN)✓ Development workflow completed$(NC)"

# Run the application
.PHONY: run
run: ## Run the application with optional ARGS
	@echo "$(YELLOW)Running application...$(NC)"
	@go run $(MAIN_FILE) $(ARGS)

# Setup development environment
.PHONY: setup
setup: ## Setup development environment
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@go mod download
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest || echo "$(RED)Failed to install golangci-lint$(NC)"
	@echo "$(GREEN)✓ Development environment ready!$(NC)"

# Check for Go and required tools
.PHONY: check
check: ## Check prerequisites
	@echo "Checking prerequisites..."
	@which go >/dev/null || (echo "Go is not installed" && exit 1)
	@echo "Go version: $$(go version)"
	@echo "Project structure check..."
	@test -f $(MAIN_FILE) || (echo "Main file $(MAIN_FILE) not found" && exit 1)
	@test -f go.mod || (echo "go.mod not found" && exit 1)
	@echo "Test files found: $(words $(TEST_FILES))"
	@echo "All prerequisites satisfied!"

# Show build information
.PHONY: info
info: ## Show build and project information
	@echo "$(CYAN)TigerGraph CLI Build Information$(NC)"
	@echo "$(WHITE)================================$(NC)"
	@echo "$(YELLOW)Project Details:$(NC)"
	@echo "  Binary Name: $(PURPLE)$(BINARY_NAME)$(NC)"
	@echo "  Module Name: $(PURPLE)$(MODULE_NAME)$(NC)"
	@echo "  Version:     $(PURPLE)$(VERSION)$(NC)"
	@echo "  Main File:   $(PURPLE)$(MAIN_FILE)$(NC)"
	@echo "  Build Dir:   $(PURPLE)$(BUILD_DIR)$(NC)"
	@echo ""
	@echo "$(YELLOW)Environment:$(NC)"
	@echo "  Go Version:  $(PURPLE)$(go version)$(NC)"
	@echo "  Platform:    $(PURPLE)$(go env GOOS)/$(go env GOARCH)$(NC)"
	@echo ""
	@echo "$(YELLOW)Code Statistics:$(NC)"
	@echo "  Go Files:    $(PURPLE)$(words $(GO_FILES)) files$(NC)"
	@echo "  Test Files:  $(PURPLE)$(words $(TEST_FILES)) files$(NC)"

# Generate documentation
.PHONY: docs
docs: build ## Generate CLI documentation
	@echo "Generating CLI documentation..."
	@mkdir -p docs
	@./$(BINARY_NAME) --help > docs/help.txt 2>&1 || echo "Help command failed"
	@./$(BINARY_NAME) cloud --help > docs/cloud-help.txt 2>&1 || echo "Cloud help failed"
	@./$(BINARY_NAME) server --help > docs/server-help.txt 2>&1 || echo "Server help failed"
	@./$(BINARY_NAME) conf --help > docs/conf-help.txt 2>&1 || echo "Conf help failed"
	@echo "Documentation generated in docs/"

# Security scan
.PHONY: security
security: ## Run security scan
	@echo "Running security scan..."
	@gosec ./... || echo "gosec not installed, skipping security scan"

# Vendor dependencies
.PHONY: vendor
vendor: ## Vendor Go modules
	@echo "Vendoring dependencies..."
	@go mod vendor

# Full CI pipeline
.PHONY: ci
ci: check mod fmt lint test build ## Run full CI pipeline
	@echo "$(GREEN)✓ CI pipeline completed successfully!$(NC)"

# Release preparation
.PHONY: release
release: clean ci test-coverage build-all package ## Prepare release with full testing
	@echo "$(GREEN)✓ Release $(VERSION) prepared!$(NC)"
	@echo "$(YELLOW)Packages available in $(BUILD_DIR)/packages/$(NC)"
	@ls -la $(BUILD_DIR)/packages/

# Debug build with additional flags
.PHONY: debug
debug: mod ## Build debug version
	@echo "Building debug version..."
	@go build -gcflags="all=-N -l" -ldflags "-X '$(MODULE_NAME)/pkg/constants.VERSION_CLI=$(VERSION)-debug'" -o $(BINARY_NAME)-debug $(MAIN_FILE)

# Cross-platform verification
.PHONY: verify
verify: build-all ## Verify all builds
	@echo "Verifying builds..."
	@for binary in $(BUILD_DIR)/$(BINARY_NAME)-*; do \
		echo "Checking $$binary..."; \
		file $$binary 2>/dev/null || echo "File command not available"; \
	done

# Test with verbose output and specific pattern
.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@go test -v -run $(PATTERN) ./... 2>&1 | tee test_output.log

# Clean test cache and run tests
.PHONY: test-clean
test-clean: ## Clean test cache and run tests
	@echo "Cleaning test cache and running tests..."
	@go clean -testcache
	@go test -v ./...

# Generate test report
.PHONY: test-report
test-report: ## Generate detailed test report
	@echo "Generating test report..."
	@go test -v -json ./... > test_report.json
	@echo "Test report generated: test_report.json"

# Run integration tests (if any)
.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test -v -tags=integration ./...

# Run unit tests only
.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v -short -tags=unit ./...

# Performance profiling
.PHONY: profile
profile: ## Run performance profiling
	@echo "Running performance profiling..."
	@go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...
	@echo "Profiles generated: cpu.prof, mem.prof"

# Watch for changes and run tests
.PHONY: watch
watch: ## Watch for changes and run tests (requires fswatch)
	@echo "Watching for changes..."
	@fswatch -o . | xargs -n1 -I{} make test-short || echo "fswatch not installed"

# Release preparation with version tagging
.PHONY: release-prepare
release-prepare: ## Prepare a new release (usage: make release-prepare VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make release-prepare VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Preparing release $(VERSION)..."
	@chmod +x scripts/release.sh
	@./scripts/release.sh $(VERSION)

# Quick release with current version
.PHONY: release-quick
release-quick: clean ci build-all package ## Quick release with current version
	@echo "Release $(VERSION) prepared!"
	@echo "Packages available in $(BUILD_DIR)/packages/"
	@echo ""
	@echo "To create a GitHub release:"
	@echo "1. git tag v$(VERSION)"
	@echo "2. git push origin v$(VERSION)"
	@echo "3. GitHub Actions will automatically create the release"

# Release dry run (test without creating tag)
.PHONY: release-dry-run
release-dry-run: clean ci test-coverage build-all package ## Test release process without creating tag
	@echo "✓ Release dry run completed successfully!"
	@echo "All checks passed. Ready for actual release."

# Tag and trigger release
.PHONY: tag-release
tag-release: ## Create and push release tag (usage: make tag-release VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make tag-release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release tag $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "✓ Tag $(VERSION) created and pushed!"
	@echo "GitHub Actions will now build and release automatically."

# Check release prerequisites
.PHONY: release-check
release-check: ## Check if repository is ready for release
	@echo "Checking release prerequisites..."
	@git diff-index --quiet HEAD -- || (echo "Working directory not clean" && exit 1)
	@git branch --show-current | grep -E '^(main|master)$$' > /dev/null || (echo "Not on main/master branch" && exit 1)
	@make test-short > /dev/null || (echo "Tests failing" && exit 1)
	@make build > /dev/null || (echo "Build failing" && exit 1)
	@echo "Repository is ready for release!"

# Show current version
.PHONY: version
version: ## Show current version
	@echo "Current version: $(VERSION)"
	@git tag --sort=-version:refname | head -1 | sed 's/^/Latest tag: /' || echo "No tags found"