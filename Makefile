.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags '-X main.version=$(RELEASE_NAME)'

.PHONY: generate
generate:
	go tool wire ./pkg/di
	rm -fr mocks/
	go tool mockery

.PHONY: lint
lint:
	go tool golangci-lint run
