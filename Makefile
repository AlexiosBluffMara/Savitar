VERSION ?= dev
BINARY := bin/savitar

.PHONY: help bootstrap verify doctor fmt test build run

help:
	@echo "Savitar targets"
	@echo "  bootstrap  Install local Mac prerequisites"
	@echo "  verify     Check workstation prerequisites"
	@echo "  doctor     Run the Savitar runtime doctor"
	@echo "  fmt        Format Go source"
	@echo "  test       Run unit tests"
	@echo "  build      Build the Savitar CLI"
	@echo "  run        Run the Savitar CLI"

bootstrap:
	./scripts/bootstrap-macos.sh

verify:
	./scripts/verify-prereqs.sh

doctor:
	go run -ldflags "-X main.version=$(VERSION)" ./cmd/savitar doctor

fmt:
	gofmt -w $$(find ./cmd ./internal -name '*.go' -print)

test:
	go test ./...

build:
	mkdir -p bin
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) ./cmd/savitar

run:
	go run -ldflags "-X main.version=$(VERSION)" ./cmd/savitar