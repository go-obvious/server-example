SHELL := /bin/bash

.PHONY: build
build:
	@echo "Building lifecycle example..."
	go build -o bin/lifecycle .

.PHONY: run
run: build
	@echo "Running lifecycle example with graceful shutdown..."
	@echo "Press Ctrl+C to trigger graceful shutdown"
	./bin/lifecycle

.PHONY: run-with-config
run-with-config: build
	@echo "Running with custom configuration..."
	DATABASE_MAX_CONNECTIONS=5 \
	WORKER_INTERVAL=10s \
	WORKER_MAX_JOBS=50 \
	./bin/lifecycle

.PHONY: test-apis
test-apis:
	@echo "Testing database API..."
	curl -s http://localhost:8080/api/database/users | jq .
	@echo
	@echo "Testing worker API..."
	curl -s http://localhost:8080/api/worker/status | jq .
	@echo
	@echo "Creating a job..."
	curl -s -X POST http://localhost:8080/api/worker/jobs \
		-H "Content-Type: application/json" \
		-d '{"type":"example","data":"test data"}' | jq .

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build           - Build the lifecycle example"
	@echo "  run             - Run the example server"
	@echo "  run-with-config - Run with custom environment variables"
	@echo "  test-apis       - Test the API endpoints"
	@echo "  clean           - Remove build artifacts"