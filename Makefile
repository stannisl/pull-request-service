.PHONY: test-e2e
test-e2e:
	@echo "Running E2E tests..."
	go test -v ./tests/e2e/... -timeout=5m

.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	go test -v ./internal/... -timeout=1m

.PHONY: test-load
test-load:
	@echo "Running load tests..."
	k6 run ./tests/load_tests/load_test.js


.PHONY: test-all
test-all: test-load test-unit test-e2e