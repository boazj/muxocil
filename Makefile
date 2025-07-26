EXTERNAL_VARS := $(.VARIABLES)

# Variables
BINARY_NAME = muxocil
MAIN_FILE = main.go
BUILD_DIR = bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildTime=${BUILD_TIME} -w -s"

# Go related variables
GOCMD ?= go
GOGEN = $(GOCMD) generate
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
# BINARY_UNIX=$(BINARY_NAME)_unix

# Dependencies 
TOOL_AIR ?= github.com/air-verse/air@latest
TOOL_GODOC ?= golang.org/x/tools/cmd/godoc@latest
TOOL_STRINGER ?= golang.org/x/tools/cmd/stringer@latest
TOOL_GOSEC ?= github.com/securego/gosec/v2/cmd/gosec@latest
TOOL_GOVULNCHECK ?= golang.org/x/vuln/cmd/govulncheck@latest
TOOL_GOLANGCI_LINT ?= github.com/golangci/golangci-lint/cmd/golangci-lint@latest
TOOL_GOPLS ?= golang.org/x/tools/gopls@latest
TOOL_GOPLS_MODERNIZE ?= golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest

# Get all the simple variables defined in the make run so far
# Allows to compute values dynamically
MAKE_VARS := $(filter-out $(EXTERNAL_VARS), $(.VARIABLES))

TOOLS := $(foreach tool,$(filter TOOL_%,$(MAKE_VARS)), $(shell echo $($(tool))))
TOOL_CMDS := $(foreach tool, $(TOOLS), $(shell echo $(tool) | sed 's/^.*\/\(.*\)@.*/\1/g'))

define newline


endef

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)



###########################
# Build & Release targets #
###########################

.PHONY: codegen
codegen: check-tools ## Auto-generates code where relevant (like stringer) 
	@echo "Auto-generating code..."
	go generate ./...

.PHONY: build-prepare
build-prepare: check ## Prepares for build execution

.PHONY: build
build: codegen ## Local build, minimal slowdowns
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

.PHONY: build-linux
build-linux: build-prepare ## Build for Linux
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_FILE)

.PHONY: build-mac
build-mac: build-prepare ## Build for macOS
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-mac $(MAIN_FILE)

.PHONY: build-windows
build-windows: build-prepare ## Build for Windows
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe $(MAIN_FILE)

.PHONY: build-all
build-all: build-prepare build-linux build-mac build-windows ## Build the application for multiple platforms
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)

.PHONY: release
release: clean build-all ## Create release builds
	@echo "Creating release builds..."
	@mkdir -p release
	@cd $(BUILD_DIR) && for file in *; do \
		tar -czf ../release/$$file.tar.gz $$file; \
	done
	@echo "Release artifacts created in release/ directory"

#################
# Clean targets #
#################

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	# rm -f $(BINARY_UNIX)

################
# Test targets #
################

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

#########################
# Dependency management #
#########################

.PHONY: deps-go
deps-go: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download

.PHONY: deps-tools
deps-tools: ## Install development tools
	@echo "Installing development tools..."
	$(foreach t, $(TOOLS), \
		$(GOCMD) install $(t) $(newline) \
	)


.PHONY: deps
deps: deps-go deps-tools

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy


####################
# Checks & Linting #
####################

.PHONY: check-tools
check-tools: ## Verify all dependent tools are available
	@echo "Verifying tool dependencies availability..."
	$(foreach cmd,$(TOOL_CMDS),\
		$(if $(shell command -v $(cmd) 2> /dev/null),,$(error Missing dependency `$(cmd)`. Install with: make deps )))


.PHONY: check-fmt
check-fmt: codegen ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

.PHONY: check-vet
check-vet: codegen ## Vet code
	@echo "Vetting code..."
	$(GOCMD) vet ./...

.PHONY: check-lint
check-lint: check-tools codegen ## Run linter
	@echo "Running linter..."
	golangci-lint run

.PHONY: check-sec
check-sec: check-tools codegen ## Run security checks
	@echo "Running security checks..."
	gosec ./...

.PHONY: check-vuln
check-vuln: check-tools codegen ## Check for vulnerabilities
	@echo "Checking for vulnerabilities..."
	govulncheck ./...

.PHONY: check-code
check-code: check-fmt check-vet check-lint  ## Run code relevant checks
	@echo "Code checks completed"
	
.PHONY: check-sec
check-security: check-sec check-vuln ## Run security relevant checks
	@echo "Security checks completed"

.PHONY: check
check: check-code test check-security ## Run all checks
	@echo "All checks completed"

#######################
# Development targets #
#######################

.PHONY: dev
dev: codegen ## Run in development mode
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_FILE)

.PHONY: dev-watch
dev-watch: check-tools codegen ## Run with file watching (requires air)
	@echo "Running with file watching..."
	air

.PHONY: install
install: build ## Install the application
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

.PHONY: uninstall
uninstall: ## Uninstall the application
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f /usr/local/bin/$(BINARY_NAME)

#################################
# Additional artificats targets #
#################################

.PHONY: docs
docs: check-tools ## Generate documentation
	@echo "Generating documentation..."
	godoc -http=:6060

###################
# Utility targets #
###################

.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit Hash: $(COMMIT_HASH)"
	@echo "Build Time: $(BUILD_TIME)"


.PHONY: ci
ci: deps check build-all ## Run CI pipeline
	@echo "CI pipeline completed"

