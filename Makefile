GO_MODULE := $(shell git config --get remote.origin.url | grep -o 'github\.com[:/][^.]*' | tr ':' '/')
CMD_NAME := kanopy-codegen
VERSION := $(shell grep '^const Version =' internal/version/version.go | cut -d\" -f2)

RUN ?= .*
PKG ?= ./...

.PHONY: test
test: ## Run tests in local environment
	golangci-lint run --timeout=5m $(PKG)
	go test -short -cover -run=$(RUN) $(PKG)

.PHONY: build
build:
	go build -o ./bin/ ./cmd/kanopy-codegen

.PHONY: dist
dist: ## Cross compile binaries into ./dist/
	mkdir -p bin dist
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(CMD_NAME) ./cmd/$(CMD_NAME)/
	tar -zcvf dist/$(CMD_NAME)-linux-$(VERSION).tgz ./bin/$(CMD_NAME) README.md
	GOOS=darwin GOARCH=amd64 go build -o ./bin/$(CMD_NAME) ./cmd/$(CMD_NAME)/
	tar -zcvf dist/$(CMD_NAME)-macos-$(VERSION).tgz ./bin/$(CMD_NAME) README.md

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
