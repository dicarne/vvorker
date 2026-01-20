# Makefile for Vorker Project
# This Makefile replaces the Windows batch scripts for cross-platform compatibility

# Variables
# Detect Windows and use appropriate command to read VERSION.txt
ifeq ($(OS),Windows_NT)
	VERSION := $(shell type VERSION.txt 2>nul || echo "unknown")
else
	VERSION := $(shell cat VERSION.txt 2>/dev/null || echo "unknown")
endif
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GO := go
PNPM := pnpm
DOCKER := docker
DOCKER_COMPOSE := docker-compose


# Default target
.PHONY: all
all: build

# Build the main Go binary
.PHONY: build
build:
	@echo Building main binary...
	$(GO) build -o vvorker

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo Building for Linux...
	GOOS=linux GOARCH=amd64 $(GO) build -o vvorker

# Build all frontend
.PHONY: build-all
build-all: build-admin build-cli build-sdk build-go

# Build admin frontend
.PHONY: build-admin
build-admin:
	@echo Building admin frontend...
	cd admin && $(PNPM) i && $(PNPM) run build

# Build CLI
.PHONY: build-cli
build-cli:
	@echo Building CLI...
	cd cli && $(PNPM) i && $(PNPM) run build


# Build SDK
.PHONY: build-sdk
build-sdk:
	@echo Building SDK...
	cd sdk/js && $(PNPM) i && $(PNPM) update && $(PNPM) run build

# Build all extensions
.PHONY: build-ext
build-ext:
	@echo Building all extensions...
	@for dir in control ai pgsql mysql oss kv assets task; do \
		echo Building $$dir...; \
		cd ext/$$dir && $(PNPM) i && $(PNPM) run build && cd ../..; \
	done

# Build Go binary (after building extensions)
.PHONY: build-go
build-go:
	@echo Building Go binary...
	$(GO) build -o vvorker


# Docker build
.PHONY: docker-build
docker-build:
	@echo Building Docker image version $(VERSION)...
	$(DOCKER) buildx build --platform linux/amd64 -t git.cloud.zhishudali.ink/dicarne/vvorker:latest --push .

# Docker build with version tag
.PHONY: docker-build-version
docker-build-version:
	@echo Building Docker image version $(VERSION)...
	$(DOCKER) buildx build --platform linux/amd64 -t git.cloud.zhishudali.ink/dicarne/vvorker:$(VERSION) --push .

# Test with Docker Compose
.PHONY: test
test:
	@echo Running tests with Docker Compose...
	$(DOCKER) rmi vorker-vorker-agent vorker-vorker-master || true
	$(DOCKER_COMPOSE) up
	$(DOCKER_COMPOSE) down

# Test with Docker Compose (dev)
.PHONY: test-dev
test-dev:
	@echo Running tests with Docker Compose (dev)...
	$(DOCKER) rmi vorker-vorker-agent vorker-vorker-master || true
	$(DOCKER_COMPOSE) -f docker-compose-dev.yml up
	$(DOCKER_COMPOSE) -f docker-compose-dev.yml down

# Clean build artifacts
.PHONY: clean
clean:
	@echo Cleaning build artifacts...
	rm -f vvorker
	rm -rf admin/dist
	rm -rf cli/dist
	rm -rf sdk/js/dist

# Clean all (including node_modules)
.PHONY: clean-all
clean-all: clean
	@echo Cleaning all dependencies...
	rm -rf admin/node_modules
	rm -rf cli/node_modules
	rm -rf sdk/js/node_modules

# Install dependencies
.PHONY: install
install:
	@echo Installing Go dependencies...
	$(GO) mod tidy
	@echo Installing admin dependencies...
	cd admin && $(PNPM) i
	@echo Installing CLI dependencies...
	cd cli && $(PNPM) i
	@echo Installing SDK dependencies...
	cd sdk/js && $(PNPM) i

# Development server
.PHONY: dev
dev:
	@echo Starting development server...
	$(GO) build
	vvorker

# Format code
.PHONY: fmt
fmt:
	@echo Formatting Go code...
	$(GO) fmt ./...
	@echo Formatting TypeScript code...
	cd admin && $(PNPM) run format || true
	cd cli && $(PNPM) run format || true

# Lint code
.PHONY: lint
lint:
	@echo Linting Go code...
	$(GO) vet ./...
	@echo Linting TypeScript code...
	cd admin && $(PNPM) run lint || true
	cd cli && $(PNPM) run lint || true

# Run tests
.PHONY: test-go
test-go:
	@echo Running Go tests...
	$(GO) test ./...

# Generate code
.PHONY: generate
generate:
	@echo Generating code...
	$(GO) generate ./...

# Show help
.PHONY: help
help:
	@echo Available targets:
	@echo   build          - Build the main Go binary
	@echo   build-linux    - Build for Linux
	@echo   build-all      - Build all frontend
	@echo   build-ext      - Build all extensions
	@echo   build-admin    - Build admin frontend
	@echo   build-cli      - Build CLI
	@echo   build-sdk      - Build SDK
	@echo   docker-build   - Build Docker image
	@echo   test           - Run tests with Docker Compose
	@echo   clean          - Clean build artifacts
	@echo   clean-all      - Clean all (including node_modules)
	@echo   install        - Install dependencies
	@echo   dev            - Start development server
	@echo   fmt            - Format code
	@echo   lint           - Lint code
	@echo   test-go        - Run Go tests
	@echo   generate       - Generate code
	@echo   help           - Show this help message
