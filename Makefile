.PHONY: clean all init generate generate_mocks

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<


init: generate
	go mod tidy
	go mod vendor

test:
	go test -short -coverprofile coverage.out -v ./...