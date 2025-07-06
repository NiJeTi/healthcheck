# cannot be migrated to go tool
# reason: https://golangci-lint.run/welcome/install/#install-from-sources
GOLANGCI_LINT_IMAGE=golangci/golangci-lint:v2.2-alpine

.PHONY: deps
deps:
	docker pull $(GOLANGCI_LINT_IMAGE)

.PHONY: mocks
mocks:
	rm -rf ./internal/generated/mocks
	mkdir -p ./internal/generated/mocks

	go tool mockery

.PHONY: fmt
fmt:
	$(MAKE) deps

	docker run -t --rm -v $(PWD):/src -w /src $(GOLANGCI_LINT_IMAGE) \
		golangci-lint fmt

.PHONY: lint
lint:
	$(MAKE) deps

	docker run -t --rm -v $(PWD):/src -w /src $(GOLANGCI_LINT_IMAGE) golangci-lint run

.PHONY: test
test:
	./scripts/test.sh
