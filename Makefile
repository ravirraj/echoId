BINARY=echoid
GOFLAGS=-ldflags="-s -w"

.PHONY: build run test bench vet lint clean docker help

build:
	go build $(GOFLAGS) -o ./bin/$(BINARY) ./cmd/fingerprint

run:
	go run ./cmd/fingerprint

test:
	go test -v -race ./...

bench:
	go test -bench=. -benchmem ./...

vet:
	go vet ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		go vet ./...; \
	fi

clean:
	rm -rf bin/ temp/

docker:
	docker build -t $(BINARY) .

help:
	@echo "Targets:"
	@echo "  build   - Build binary to ./bin/$(BINARY)"
	@echo "  run     - Run with 'go run'"
	@echo "  test    - Run tests with -v -race"
	@echo "  bench   - Run benchmarks"
	@echo "  vet     - Run go vet"
	@echo "  lint    - Run golangci-lint (falls back to vet)"
	@echo "  clean   - Remove bin/ and temp/"
	@echo "  docker  - Build Docker image"
	@echo "  help    - Show this help"
