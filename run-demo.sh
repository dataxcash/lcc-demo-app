#!/bin/bash

echo "=== LCC Demo Application Setup ==="
echo ""

# Check if LCC server is running
echo "Checking LCC server..."
if ! curl -sk https://localhost:8088/api/lcc/info > /dev/null 2>&1; then
    echo "⚠ Warning: LCC server not responding at https://localhost:8088"
    echo "   Please start LCC server first:"
    echo "   cd /home/fila/jqdDev_2025/lcc && ./lcc_server"
    echo ""
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    echo "✓ LCC server is running"
fi

echo ""
echo "Building demo application..."
cd "$(dirname "$0")"

# Download dependencies
echo "  - Downloading dependencies..."
go mod download 2>/dev/null || true
go mod tidy 2>/dev/null || true

# Build
echo "  - Compiling..."
if go build -o demo cmd/demo/main.go; then
    echo "✓ Build successful"
    echo ""
    echo "Starting demo application..."
    echo "================================"
    echo ""
    ./demo
else
    echo "✗ Build failed"
    exit 1
fi
