.PHONY: check release

all: ghcp

check:
	go vet
	go test -v ./...

ghcp: check
	go build

dist:
	gox --osarch 'darwin/amd64 linux/amd64 windows/amd64' -output 'dist/bin/{{.Dir}}_{{.OS}}_{{.Arch}}'
	./.circleci/homebrew.sh > dist/ghcp.rb

release: ghcp dist
	ghr -u $(CIRCLE_PROJECT_USERNAME) -r $(CIRCLE_PROJECT_REPONAME) -b "$$(ghch -F markdown --latest)" $(CIRCLE_TAG) dist/bin
	./ghcp -u $(CIRCLE_PROJECT_USERNAME) -r homebrew-$(CIRCLE_PROJECT_REPONAME) -m $(CIRCLE_TAG) -C dist/ ghcp.rb
