NAME   := gcr.io/timatal1-198022/app-blacklist

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

##############
# CI targets #
##############

CI_IMG_VERSION 	 := $(shell ./.circleci/scripts/image_version.sh)
CI_CHART_VERSION := $(shell ./.circleci/scripts/chart_version.sh)

ci-build-image:
	./.circleci/scripts/image_build.sh ${NAME}

ci-push-image:
	./.circleci/scripts/image_push.sh ${NAME}

ci-build-chart:
	./.circleci/scripts/chart_build.sh ${CI_CHART_VERSION} ${CI_IMG_VERSION}

ci-push-chart:
	./.circleci/scripts/chart_push.sh ${CI_CHART_VERSION} ${CI_IMG_VERSION}

ci-noona-sync:
	./.circleci/scripts/trigger_noona_sync.sh