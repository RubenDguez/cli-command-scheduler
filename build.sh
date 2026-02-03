#!/bin/bash

echo "Building SCHEDULER CLI for multiple platforms"

mkdir -p bin

echo "- Building for macOS (Intel)"
GOOS=darwin GOARCH=amd64 go build -o bin/scheduler-darwin-amd64 main.go

echo "- Building for macOS (Apple Silicon)"
GOOS=darwin GOARCH=arm64 go build -o bin/scheduler-darwin-arm64 main.go

echo "- Building for Windows"
GOOS=windows GOARCH=amd64 go build -o bin/scheduler-windows.exe main.go

echo "- Building for Linux"
GOOS=linux GOARCH=amd64 go build -o bin/scheduler-linux-amd64 main.go

echo "✓ Build complete! Binaries are in the bin/ directory"
ls -lh bin/
