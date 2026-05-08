# List available recipes
default:
    @just --list

# Standard Vetting
vet:
    go vet ./...

tidy:
    go mod tidy

# Comprehensive golangci-lint
lint:
    golangci-lint run ./...
