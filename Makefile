CIRCLE_TAG ?= latest
LDFLAGS := -X main.version=$(CIRCLE_TAG)

all: ghcp

.PHONY: check
check:
	golangci-lint run
	go test -v -race -cover -coverprofile=coverage.out ./...

ghcp: $(wildcard **/*.go)
	go build -o $@ -ldflags "$(LDFLAGS)"

dist:
	# make the zip files for GitHub Releases
	VERSION=$(CIRCLE_TAG) CGO_ENABLED=0 goxzst -d dist/ -i "LICENSE" -o ghcp -t "ghcp.rb" -- -ldflags "$(LDFLAGS)"

.PHONY: acceptance-test
acceptance-test: ghcp
	make -C acceptance_test

.PHONY: release
release: ghcp dist
	# publish to the GitHub Releases
	./ghcp release -u "$(CIRCLE_PROJECT_USERNAME)" -r "$(CIRCLE_PROJECT_REPONAME)" -t "$(CIRCLE_TAG)" dist/
	# publish to the Homebrew tap repository
	./ghcp commit -u "$(CIRCLE_PROJECT_USERNAME)" -r "homebrew-$(CIRCLE_PROJECT_REPONAME)" -b "bump-$(CIRCLE_TAG)" \
		-m "Bump the version to $(CIRCLE_TAG)" \
		-C dist/ ghcp.rb
