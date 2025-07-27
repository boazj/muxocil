EXTERNAL_VARS := $(.VARIABLES)

# Variables
BINARY_NAME = muxocil
MAIN_FILE = main.go
BUILD_DIR = dist
RELEASE_DIR = release
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

# Dependencies 
TOOL_AIR ?= github.com/air-verse/air@latest
TOOL_GODOC ?= golang.org/x/tools/cmd/godoc@latest
TOOL_STRINGER ?= golang.org/x/tools/cmd/stringer@latest
TOOL_GOSEC ?= github.com/securego/gosec/v2/cmd/gosec@latest
TOOL_GOVULNCHECK ?= golang.org/x/vuln/cmd/govulncheck@latest
TOOL_GOLANGCI_LINT ?= github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.3.0
TOOL_GOPLS ?= golang.org/x/tools/gopls@latest
TOOL_GOPLS_MODERNIZE ?= golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest
TOOL_SCC ?= github.com/boyter/scc/v3@latest

CMD_TOOL_SCC = scc

# Get all the simple variables defined in the make run so far
# Allows to compute values dynamically
MAKE_VARS := $(filter-out $(EXTERNAL_VARS), $(.VARIABLES))

# List of all the tools urls
TOOLS_VARS := $(filter TOOL_%,$(MAKE_VARS))
# List all the tools install urls
TOOLS := $(foreach tool, $(TOOLS_VARS), $(shell echo $($(tool))))
# List all the commands from the tools (supports override via CMD_$(TOOL_<name>)
TOOL_CMDS := $(foreach tool, $(TOOLS_VARS), $(if $(value CMD_$(tool)), $(CMD_$(tool)), $(shell echo $($(tool)) | sed 's/^.*\/\(.*\)@.*/\1/g')))

define newline


endef

# Default target
.DEFAULT_GOAL := help

###  
### Build Tasks:
#
.PHONY: clean
## Clean build artifacts
clean: 
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	# rm -f $(BINARY_UNIX)

.PHONY: build
## Local default build
build: build-prepare
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

.PHONY: build-all
## Build for supported platforms
build-all: build-prepare build-linux build-mac build-windows 
	@echo "Building for multiple platforms... done"
	@mkdir -p $(BUILD_DIR)

.PHONY: release
## Create release builds
release: clean build-all 
	@echo "Creating release builds..."
	@mkdir -p $(RELEASE_DIR)
	@cd $(BUILD_DIR) && for file in *; do \
		tar -czf ../$(RELEASE_DIR)/$$file.tar.gz $$file; \
	done
	@echo "Release artifacts created in release/ directory"

.PHONY: ci
## Run CI pipeline
ci: deps check build-all 
	@echo "CI pipeline completed"

.PHONY: build-prepare
## Prepares for build execution
build-prepare: check 
	@mkdir -p $(BUILD_DIR)

$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64: build-prepare
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)

$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64: build-prepare
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)

$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64: build-prepare
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)

$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64: build-prepare
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)

$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64: build-prepare
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe $(MAIN_FILE)

.PHONY: build-linux
## Build for Linux
build-linux: build-prepare $(BUILD_DIR)/$(BINARY_NAME)-linux-%
	@echo "Building for Linux... \033[32;1;4mDone\033[0m"

.PHONY: build-mac
## Build for macOS
build-mac:  build-prepare $(BUILD_DIR)/$(BINARY_NAME)-darwin-%
	@echo "Building for macOS... \033[32;1;4mDone\033[0m"

.PHONY: build-windows
## Build for Windows
build-windows: build-prepare $(BUILD_DIR)/$(BINARY_NAME)-windows-%
	@echo "Building for Windows... \033[32;1;4mDone\033[0m"


###  
### Dependency Tasks:
#

.PHONY: deps
## Perform all dependencies tasks
deps: deps-go deps-tools

.PHONY: deps-go
## Download go dependencies
deps-go: 
	@echo "Downloading go dependencies..."
	$(GOMOD) download
	@echo "Downloading go dependencies... \033[32;1;4mDone\033[0m"


.PHONY: deps-tools
## Install dev tools
deps-tools: 
	@echo "Installing dev tools..."
	$(foreach t, $(TOOLS), \
		$(GOCMD) install $(t) $(newline) \
	)
	@echo "Installing dev tools... \033[32;1;4mDone\033[0m"

.PHONY: deps-update
## Update go dependencies
deps-update: 
	@echo "Updating go dependencies..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy
	@echo "Updating go dependencies... \033[32;1;4mDone\033[0m"

###  
### Checks & Tests Tasks:
#

.PHONY: check
## Run all checks
check: check-code test check-security 
	@echo "All checks completed"

.PHONY: test
## Run tests
test: 
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "Running tests... \033[32;1;4mDone\033[0m"

.PHONY: test-coverage
## Run tests with coverage
test-coverage: 
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out.tmp ./...
	cat coverage.out.tmp | grep -v "_string.go" > coverage.out
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Running tests with coverage... \033[32;1;4mDone\033[0m"

.PHONY: codegen
## Auto-generates (like stringer)
codegen: check-tools  
	@echo "Auto-generating code..."
	go generate ./...
	@echo "Auto-generating code... \033[32;1;4mDone\033[0m"

.PHONY: check-tools
## Verify dev tools are installed
check-tools: 
	@echo "Verifying dev tools..."
	$(foreach cmd,$(TOOL_CMDS),\
		$(if $(shell command -v $(cmd) 2> /dev/null),,$(error Missing dependency `$(cmd)`. Install with: make deps )))
	@echo "Verifying dev tools... \033[32;1;4mDone\033[0m"


.PHONY: check-fmt
## Format code
check-fmt: codegen 
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "Formatting code... \033[32;1;4mDone\033[0m"

.PHONY: check-vet
## Vet code
check-vet: codegen 
	@echo "Vetting code..."
	$(GOCMD) vet ./...
	@echo "Vetting code... \033[32;1;4mDone\033[0m"

.PHONY: check-lint
## Run linter
check-lint: check-tools codegen 
	@echo "Running linter..."
	golangci-lint run
	@echo "Running linter... \033[32;1;4mDone\033[0m"

.PHONY: check-sec
## Run security checks
check-sec: check-tools codegen 
	@echo "Running security checks..."
	gosec ./...
	@echo "Running security checks... \033[32;1;4mDone\033[0m"

.PHONY: check-vuln
## Check for vulnerabilities
check-vuln: check-tools codegen 
	@echo "Running vulnerability checks..."
	govulncheck ./...
	@echo "Running vulnerability checks... \033[32;1;4mDone\033[0m"

.PHONY: check-code
## Run code relevant checks
check-code: check-fmt check-vet check-lint  
	@echo "Code checks... \033[32;1;4mDone\033[0m"
	
.PHONY: check-sec
## Run security relevant checks
check-security: check-sec check-vuln 
	@echo "Security checks... \033[32;1;4mDone\033[0m"


###  
### Development Tasks:
#

.PHONY: dev
## Run in development mode
dev: codegen 
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_FILE)

.PHONY: dev-watch
## Run with file watching
dev-watch: check-tools codegen 
	@echo "Running with file watching..."
	air

.PHONY: dev-metrics
## Print development metrics
dev-metrics: 
	@echo "Generating development metrics... \033[32;1;4mDone\033[0m"
	@git ls-files -co --exclude-standard --full-name -- . ':!:*.md' ':!:go.*' ':!:LICENSE' | scc

###  
### Other Artifacts Tasks:
#

.PHONY: docs
## Generate documentation
docs: check-tools 
	@echo "Generating documentation..."
	godoc -http=:6060
	@echo "Generating documentation... \033[32;1;4mDone\033[0m"

.PHONY: install
## Install to /usr/local/bin
install: build 
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Installing $(BINARY_NAME)... \033[32;1;4mDone\033[0m"

.PHONY: uninstall
## Uninstall from /usr/local/bin
uninstall: 
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstalling $(BINARY_NAME)... \033[32;1;4mDone\033[0m"

###  
### Utility Tasks:
#

.PHONY: version
## Show version information
version : 
	@echo "Version: $(VERSION)"
	@echo "Commit Hash: $(COMMIT_HASH)"
	@echo "Build Time: $(BUILD_TIME)"

.PHONY: help
## Show this help message
help :
	@echo 
	@echo Makefile help for $(BINARY_NAME)
	@echo 
	@echo To exclude a task from execution use the -o flag: make build -o check
	@echo 
	@printf "    %-17s     %-26s     [%s]\n" "Task" "Description" "Direct Dependencies"
	@printf " ==================== ============================== ==============================\n"
	@awk -v dir="$(BUILD_DIR)" -v app="$(BINARY_NAME)" '/^### /, /^[:alpha:][[:alnum:]_-]+\s*:/ { \
		if ($$0 ~ /^### /){ title = substr($$0, 5); print title; prev = $$0; next; } \
		if (prev !~ /^## /){ prev = $$0; next; } \
		desc = substr(prev, 4); \
		if (! match($$0, /^([[:alpha:]][[:alnum:]_-]+)\s*:(.*)/, m)) { prev = $$0; next; } \
		task = m[1]; deps = m[2]; \
		gsub(/\$$\(BUILD_DIR\)/, dir , task); \
		gsub(/\$$\(BINARY_NAME\)/, app , task); \
		gsub(/^[ \t]+/,"", deps); \
		gsub(/[ \t]+$$/,"", deps); \
		gsub(/\$$\(BUILD_DIR\)/, dir , deps); \
		gsub(/\$$\(BINARY_NAME\)/, app , deps); \
		printf "    %-17s %-30s [%s]\n", task, desc, deps; \
		prev = $$0; \
		}' $(MAKEFILE_LIST)

