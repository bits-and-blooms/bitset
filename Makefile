# MAKEFILE

SHELL=/bin/bash
.SHELLFLAGS=-o pipefail -c

# Project name
PROJECT=bitset

# Project version
VERSION=$(shell cat VERSION)

# Current directory
CURRENTDIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

# Target directory
TARGETDIR=$(CURRENTDIR)target

# Directory where to store binary utility tools
BINUTIL=$(TARGETDIR)/binutil

# GO lang path
ifeq ($(GOPATH),)
	# extract the GOPATH
	GOPATH=$(firstword $(subst /src/, ,$(CURRENTDIR)))
endif

# Add the GO binary dir in the PATH
export PATH := $(GOPATH)/bin:$(PATH)

# sed argument for in-place substitutions
SEDINPLACE=-i
ifeq ($(shell uname -s),Darwin)
	SEDINPLACE=-i ''
endif

# Common commands
GO=GOPATH=$(GOPATH) $(shell which go)
GOVERSION=${shell go version | grep -Eo '(go[0-9]+.[0-9]+)'}
GOFMT=$(shell which gofmt)
GOTEST=$(GO) test
GODOC=GOPATH=$(GOPATH) $(shell which godoc)
GOLANGCILINT=$(BINUTIL)/golangci-lint
GOLANGCILINTVERSION=v2.12.2

# Directory containing the source code
SRCDIR=./

# List of packages
GOPKGS=$(shell $(GO) list $(SRCDIR)/...)

# Enable junit report when not in LOCAL mode
ifeq ($(strip $(DEVMODE)),LOCAL)
	TESTEXTRACMD=&& $(GO) tool cover -func=$(TARGETDIR)/report/coverage.out
else
	TESTEXTRACMD=2>&1 | tee >(PATH=$(GOPATH)/bin:$(PATH) go-junit-report > $(TARGETDIR)/test/report.xml); test $${PIPESTATUS[0]} -eq 0
endif

# --- MAKE TARGETS ---

.PHONY: help
help:
	@echo ""
	@echo "$(PROJECT) Makefile."
	@echo "GOPATH=$(GOPATH)"
	@echo "The following commands are available:"
	@echo ""
	@awk '/^## /{desc=substr($$0,4)} /^\.PHONY:/{if(NF>1) {target=$$2; if(desc) printf "  make %-15s: %s\n",target,desc; desc=""}}' Makefile
	@echo ""
	@echo "To test and build everything from scratch, use the shortcut:"
	@echo "    make x"
	@echo ""

# Alias for help target
all: help

## Test and build everything from scratch
.PHONY: x
x:
	DEVMODE=LOCAL $(MAKE) format clean mod deps generate qa

## Remove any build artifact
.PHONY: clean
clean:
	rm -rf $(TARGETDIR)

## Generate the coverage report
.PHONY: coverage
coverage: ensuretarget
	$(GO) tool cover -html=$(TARGETDIR)/report/coverage.out -o $(TARGETDIR)/report/coverage.html

## Get dependencies
.PHONY: deps
deps: ensuretarget
	curl --silent --show-error --fail --location "https://golangci-lint.run/install.sh" | sh -s -- -b $(BINUTIL) $(GOLANGCILINTVERSION)

## Create the target directories if missing
.PHONY: ensuretarget
ensuretarget:
	@mkdir -p $(TARGETDIR)/test
	@mkdir -p $(TARGETDIR)/report
	@mkdir -p $(TARGETDIR)/binutil

## Format the source code
.PHONY: format
format:
	@find $(SRCDIR) -type f -name "*.go" -exec $(GOFMT) -s -w {} \;

## Run fuzz tests (set FUZZTIME to override the duration or iterations, e.g. FUZZTIME=10s or FUZZTIME=1000x)
.PHONY: fuzz
fuzz:
	./run_fuzz_tests.sh $(FUZZTIME)

## Generate go code automatically
.PHONY: generate
generate:
	@find $(SRCDIR) -type f -name "*mock_test.go" -exec rm {} \;
	$(GO) generate $(GOPKGS)

## Check code against multiple linters
.PHONY: linter
linter:
	@echo -e "\n\n>>> START: Static code analysis <<<\n\n"
	$(GOLANGCILINT) run --max-issues-per-linter 0 --max-same-issues 0 $(SRCDIR)/...
	@echo -e "\n\n>>> END: Static code analysis <<<\n\n"

## Download dependencies
.PHONY: mod
mod: gotools
	$(GO) mod download all

## Run all tests and static analysis tools
.PHONY: qa
qa: linter test coverage

## Tag the Git repository
.PHONY: tag
tag:
	git tag -a "v$(VERSION)" -m "Version $(VERSION)" && \
	git push origin --tags

## Run unit tests
.PHONY: test
test: ensuretarget
	@echo -e "\n\n>>> START: Unit Tests <<<\n\n"
	$(GOTEST) \
	-shuffle=on \
	-tags=unit,benchmark \
	-covermode=atomic \
	-bench=. \
	-benchtime=1x \
	-race \
	-failfast \
	-coverprofile=$(TARGETDIR)/report/coverage.out \
	-v $(GOPKGS) $(TESTEXTRACMD)
	@echo -e "\n\n>>> END: Unit Tests <<<\n\n"

## Run benchmarks (real measurements, without -race or coverage)
.PHONY: bench
bench: ensuretarget
	@echo -e "\n\n>>> START: Benchmarks <<<\n\n"
	$(GOTEST) \
	-tags=unit,benchmark \
	-run=^$$ \
	-bench=. \
	-benchmem \
	-v $(GOPKGS)
	@echo -e "\n\n>>> END: Benchmarks <<<\n\n"

## Get the go tools
.PHONY: gotools
gotools:
	$(GO) install github.com/jstemmer/go-junit-report/v2@latest

## Update everything
.PHONY: updateall
updateall: updatelint updatemod

## Update Go version
.PHONY: updatego
updatego:
	$(eval LAST_GO_TOOLCHAIN=$(shell curl -s https://go.dev/dl/ | grep -oE 'go[0-9]+\.[0-9]+\.[0-9]+\.linux-amd64\.tar\.gz' | head -n 1 | grep -oE 'go[0-9]+\.[0-9]+\.[0-9]+'))
	$(eval LAST_GO_VERSION=$(shell echo ${LAST_GO_TOOLCHAIN} | grep -oE '[0-9]+\.[0-9]+'))
	sed $(SEDINPLACE) "s|^go [0-9]*\.[0-9]*.*$$|go ${LAST_GO_VERSION}|g" go.mod
	sed $(SEDINPLACE) "s|^toolchain go[0-9]*\.[0-9]*\.[0-9]*$$|toolchain ${LAST_GO_TOOLCHAIN}|g" go.mod

## Update golangci-lint version
.PHONY: updatelint
updatelint:
	$(eval LAST_GOLANGCILINT_VERSION=$(shell curl -sL https://github.com/golangci/golangci-lint/releases/latest | sed -n 's/.*<title>Release \(v[0-9]*\.[0-9]*\.[0-9]*\).*/\1/p'))
	sed $(SEDINPLACE) "s|^GOLANGCILINTVERSION=v[0-9]*\.[0-9]*\.[0-9]*$$|GOLANGCILINTVERSION=${LAST_GOLANGCILINT_VERSION}|g" Makefile

## Update dependencies
.PHONY: updatemod
updatemod: mod
	$(GO) get -t -u ./... && \
	$(GO) mod tidy -compat=$(shell sed -n -E 's/^go ([0-9]+\.[0-9]+).*/\1/p' go.mod)

## Increase the patch number in the VERSION file
.PHONY: versionup
versionup:
	echo ${VERSION} | awk -F. '{printf("%d.%d.%d\n",$$1,$$2,(($$3+1)));}' > VERSION
