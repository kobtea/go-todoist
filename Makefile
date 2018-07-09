GOPATH     ?= $(shell go env GOPATH)
DEP        ?= $(GOPATH)/bin/dep
GORELEASER ?= $(GOPATH)/bin/goreleaser
VERSION    := v$(shell cat VERSION)

.PHONY: setup test build build-snapshot sync-tag release
all: setup test build build-snapshot sync-tag release

setup: $(DEP)
	@echo '>> setup'
	@$(DEP) ensure -v

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

$(DEP):
	go get -u github.com/golang/dep/cmd/dep

$(GORELEASER): $(DEP)
	go get golang.org/x/tools/cmd/stringer
	go get -d github.com/goreleaser/goreleaser
	cd $(GOPATH)/src/github.com/goreleaser/goreleaser && \
	$(DEP) ensure -vendor-only && \
	make build && \
	mv ./goreleaser $(GOPATH)/bin
