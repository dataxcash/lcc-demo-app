.PHONY: all generate build run run-dev test clean

all: build

generate:
	@echo "Generating license wrappers..."
	@lcc-codegen --config lcc-features.yaml --output ./

build: generate
	@echo "Building demo-app..."
	@go build -o bin/demo-app ./cmd/demo

run: build
	@echo "Running demo-app..."
	@./bin/demo-app

run-dev:
	@echo "Running in development mode (all features unlocked)..."
	@LCC_DEV_MODE=true go run ./cmd/demo

test:
	@echo "Running tests..."
	@go test -v ./...

test-basic:
	@echo "Testing with basic tier..."
	@LCC_LICENSE_TIER=basic go test -v ./...

test-professional:
	@echo "Testing with professional tier..."
	@LCC_LICENSE_TIER=professional go test -v ./...

test-enterprise:
	@echo "Testing with enterprise tier..."
	@LCC_LICENSE_TIER=enterprise go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@find . -name "lcc_gen.go" -delete

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
