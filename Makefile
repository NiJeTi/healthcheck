.PHONY: deps
deps:
	docker pull vektra/mockery

.PHONY: mocks
mocks:
	rm -rf ./internal/generated/mocks
	docker run --rm -v ${PWD}:/src -w /src vektra/mockery

.PHONY: tests
tests:
	go test -v . ./servers/fasthttp ./servers/http
