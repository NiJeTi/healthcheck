.PHONY: deps
deps:
	go install github.com/vektra/mockery/v2@v2.43.2

.PHONY: mocks
mocks:
	rm -rf ./internal/generated/mocks
	mockery
