# build.ps1
# Build Go project for Windows and macOS (x64)

$ErrorActionPreference = "Stop"

Write-Host "Building Go project for Windows (x64)..."
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o bin/stateflow-windows-x64.exe .

Write-Host "Building Go project for macOS (x64)..."
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o bin/stateflow-macos-x64 .

Write-Host "Building Go project for macOS (arm64)..."
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o bin/stateflow-macos-arm64 .

Write-Host "Builds completed successfully!"
