TARGET := ghcp

ifndef $(CIRCLE_TAG)
GHCP_FLAGS := --dry-run
endif

CIRCLE_TAG ?= latest
LDFLAGS := -X main.version=$(CIRCLE_TAG)

all: $(TARGET)

.PHONY: check
check:
	golangci-lint run
	go test -v -race -cover -coverprofile=coverage.out ./...

$(TARGET): $(wildcard *.go)
	go build -o $@ -ldflags "$(LDFLAGS)"

dist:
	# make the zip files for GitHub Releases
	VERSION=$(CIRCLE_TAG) CGO_ENABLED=0 goxzst -d dist/ -i "LICENSE" -o "$(TARGET)" -t "ghcp.rb" -- -ldflags "$(LDFLAGS)"

.PHONY: release
release: dist $(TARGET)
	# publish to the GitHub Releases
	./ghcp release $(GHCP_FLAGS) -u "$(CIRCLE_PROJECT_USERNAME)" -r "$(CIRCLE_PROJECT_REPONAME)" -t "$(CIRCLE_TAG)" dist/
	# publish to the Homebrew tap repository
	./ghcp commit $(GHCP_FLAGS) -u "$(CIRCLE_PROJECT_USERNAME)" -r "homebrew-$(CIRCLE_PROJECT_REPONAME)" -b "bump-$(CIRCLE_TAG)" \
		-m "Bump the version to $(CIRCLE_TAG)" \
		-C dist/ ghcp.rb
