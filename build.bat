@echo off
echo Building SCHEDULER CLI for multiple platforms

if not exist bin mkdir bin

echo Building for macOS (Intel)
set GOOS=darwin
set GOARCH=amd64
go build -o bin/scheduler-cli-darwin-amd64 main.go

echo Building for macOS (Apple Silicon)
set GOOS=darwin
set GOARCH=arm64
go build -o bin/scheduler-cli-darwin-arm64 main.go

echo Building for Windows
set GOOS=windows
set GOARCH=amd4
go build -o bin/scheduler-cli-windows.exe main.go

echo Building for Linux
set GOOS=linux
set GOARCH=amd64
go build -o bin/scheduler-cli-linux-amd64 main.go

echo Build complete! Binaries are in the bin/ directory
dir bin
