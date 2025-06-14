ghcp:
	CGO_ENABLED=0 go build -o $@ -ldflags '-X main.version=$(GITHUB_REF_NAME)'

.PHONY: generate
generate:
	go tool wire ./pkg/di
	rm -fr mocks/
	go tool mockery

.PHONY: lint
lint:
	go tool golangci-lint run
