# List available recipes
default:
    @just --list

# Standard Vetting
vet:
    go vet ./...

# Comprehensive golangci-lint
lint:
    golangci-lint run ./...
