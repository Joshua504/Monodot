.PHONY: help run test benchmark cover lint fmt vet clean ci

help:
	@echo "Available commands:"
	@echo "  make run         - Run the web server"
	@echo "  make test        - Run all tests"
	@echo "  make benchmark   - Run processor benchmarks"
	@echo "  make cover       - Generate coverage report"
	@echo "  make fmt         - Format Go code"
	@echo "  make vet         - Run go vet"
	@echo "  make lint        - Run golangci-lint"
	@echo "  make clean       - Remove generated files"
	@echo "  make ci          - Run the full CI pipeline"

run:
	go run ./cmd/web

test:
	go test ./...

benchmark:
	go test -bench=. -benchmem ./internal/processor

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	golangci-lint run

clean:
	rm -f coverage.out coverage.html

ci: fmt vet test