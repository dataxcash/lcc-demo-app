.PHONY: all generate build run run-dev test clean \
	build-basic build-pro build-ent generate-basic generate-pro generate-ent

FEATURES_BASIC := configs/lcc-features.basic.yaml
FEATURES_PRO   := configs/lcc-features.pro.yaml
FEATURES_ENT   := configs/lcc-features.ent.yaml

all: build

# Default generate/build use root lcc-features.yaml (mapped to Pro profile)
generate:
	@echo "Generating license wrappers (default: Pro profile)..."
	@lcc-codegen --config lcc-features.yaml --output ./

build: generate
	@echo "Building demo-app (default: Pro profile)..."
	@go build -o bin/demo-app ./cmd/demo

# --- Multi-product builds ---

generate-basic:
	@echo "Generating license wrappers for Basic product..."
	@lcc-codegen --config $(FEATURES_BASIC) --output ./

build-basic: generate-basic
	@echo "Building demo-basic (demo-analytics-basic)..."
	@go build -o bin/demo-basic ./cmd/demo


generate-pro:
	@echo "Generating license wrappers for Pro product..."
	@lcc-codegen --config $(FEATURES_PRO) --output ./

build-pro: generate-pro
	@echo "Building demo-pro (demo-analytics-pro)..."
	@go build -o bin/demo-pro ./cmd/demo


generate-ent:
	@echo "Generating license wrappers for Enterprise product..."
	@lcc-codegen --config $(FEATURES_ENT) --output ./

build-ent: generate-ent
	@echo "Building demo-ent (demo-analytics-ent)..."
	@go build -o bin/demo-ent ./cmd/demo

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
