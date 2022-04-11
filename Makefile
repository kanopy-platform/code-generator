GO_MODULE := $(shell git config --get remote.origin.url | grep -o 'github\.com[:/][^.]*' | tr ':' '/')
CMD_NAME := $(shell basename ${GO_MODULE})

RUN ?= .*
PKG ?= ./...

.PHONY: test
test: ## Run tests in local environment
	golangci-lint run --timeout=5m $(PKG)
	go test -short -cover -run=$(RUN) $(PKG)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
