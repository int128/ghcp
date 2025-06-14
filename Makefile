.PHONY: all
all:

.PHONY: generate
generate:
	go tool wire ./pkg/di
	rm -fr mocks/
	go tool mockery

.PHONY: lint
lint:
	go tool golangci-lint run
