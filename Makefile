GOLANGCI_LINT_PACKAGE=github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2
MOCKERY_PACKAGE=github.com/vektra/mockery/v3@v3.5.1

.PHONY: deps
deps:
	go install $(GOLANGCI_LINT_PACKAGE)
	go install $(MOCKERY_PACKAGE)

.PHONY: mocks
mocks:
	rm -rf ./internal/generated/mocks
	mkdir -p ./internal/generated/mocks

	mockery

.PHONY: lint
lint:
	./scripts/lint.sh

.PHONY: test
test:
	./scripts/test.sh
