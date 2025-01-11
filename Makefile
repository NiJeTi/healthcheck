MOCKERY_VERSION=github.com/vektra/mockery/v2@v2.50
GOLANGCI_LINT_IMAGE=golangci/golangci-lint:v1.63-alpine

.PHONY: deps
deps:
	go install $(MOCKERY_VERSION)
	docker pull $(GOLANGCI_LINT_IMAGE)

.PHONY: mocks
mocks:
	$(MAKE) deps

	rm -rf ./internal/generated/mocks
	mockery

.PHONY: lint
lint:
	$(MAKE) deps

	./scripts/lint.sh $(GOLANGCI_LINT_IMAGE)

.PHONY: test
test:
	./scripts/test.sh
