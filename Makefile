GOBIN ?= $(shell go env GOBIN)

GO_ARGS=-trimpath

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Run locally
	go run main.go

.PHONY: build
build: ## Build project
	go build ${GO_ARGS} ./...

.PHONY: install
install: ## Install on local system
	go install ${GO_ARGS} ./...
