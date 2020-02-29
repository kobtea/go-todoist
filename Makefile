GOPATH     ?= $(shell go env GOPATH)
GORELEASER ?= $(GOPATH)/bin/goreleaser
VERSION    := v$(shell cat VERSION)

.PHONY: test build build-snapshot sync-tag release
all: test build build-snapshot sync-tag release

test:
	@echo '>> unit test'
	@go test ./...

build:
	@echo '>> build'
	go build -ldflags='-X github.com/kobtea/go-todoist/cmd/todoist/cmd.Version=$(shell cat VERSION)' \
	-o dist/todoist \
	./cmd/todoist

build-snapshot: $(GORELEASER)
	@echo '>> cross-build for testing'
	$(GORELEASER) release --snapshot --rm-dist --debug

sync-tag:
	@git config user.name  || git config --local user.name  "circleci-job"
	@git config user.email || git config --local user.email "kobtea9696@gmail.com"
	@git rev-parse $(VERSION) > /dev/null 2>&1 || \
	(git tag -a $(VERSION) -m "release $(VERSION)" && git push origin $(VERSION))

release: $(GORELEASER)
	@echo '>> release'
	$(GORELEASER) release --rm-dist --debug

$(GORELEASER):
	@wget -O - "https://github.com/goreleaser/goreleaser/releases/download/v0.95.0/goreleaser_$(shell uname -o | cut -d'/' -f2)_$(shell uname -m).tar.gz" | tar xvzf - -C /tmp
	@mv /tmp/goreleaser $(GOPATH)/bin
