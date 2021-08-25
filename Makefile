all: build

COMMIT := $(shell git log -1 --format='%h')
BRANCH_PRETTY := $(subst /,-,$(BRANCH))
BUILT  := $(shell date -u +%F-%T-%Z)

# don't override user values
VERSION ?= $(shell git describe --exact-match 2>/dev/null)
VERSION ?= $(BRANCH_PRETTY)-$(COMMIT)

ldflags = -w -s \
	-X github.com/mtps/tq/version.Name=tq \
	-X github.com/mtps/tq/version.Version=$(VERSION) \
	-X github.com/mtps/tq/version.Commit=$(COMMIT) \
	-X github.com/mtps/tq/version.Built=$(BUILT)

.PHONY: build
build:
	go build -ldflags '$(ldflags)' -mod=vendor ./cmd/tq

.PHONY: test
test:
	go test -mod=vendor -race ./...

.PHONY: install
install:
	go install -ldflags '$(ldflags)' -mod=vendor ./cmd/tq
