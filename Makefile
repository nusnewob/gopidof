# Makefile for gopidof

# Binary name
BINARY_NAME=pidof

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-darwin

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 .

.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 .

# Run go vet
.PHONY: vet
vet:
	$(GOVET)

# Format code
.PHONY: fmt
fmt:
	$(GOFMT)

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*

# Install the binary
.PHONY: install
install: build
	install -m 755 $(BINARY_NAME) /usr/local/bin/

# Uninstall the binary
.PHONY: uninstall
uninstall:
	rm -f /usr/local/bin/$(BINARY_NAME)

# Run all checks
.PHONY: check
check: fmt vet
