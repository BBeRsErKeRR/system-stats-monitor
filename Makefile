### SSM (system-stats-monitor) app automation recieps
# ¯¯¯¯¯¯¯¯

.DEFAULT_GOAL := help

SHELL=/bin/bash
SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

BIN_DAEMON := "./bin/ssm"
BIN_CLIENT := "./bin/ssm_client"
BIN_TEST := "./bin/test"

help: ## Display this help
	@IFS=$$'\n'; for line in `grep -h -E '^[a-zA-Z_#-]+:?.*?## .*$$' $(MAKEFILE_LIST)`; do if [ "$${line:0:2}" = "##" ]; then \
	echo $$line | awk 'BEGIN {FS = "## "}; {printf "\n\033[33m%s\033[0m\n", $$2}'; else \
	echo $$line | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'; fi; \
	done; unset IFS;

generate: ## Generate proto files
	go $@ ./...

build-daemon:
	go build -v -o $(BIN_DAEMON) -ldflags "$(LDFLAGS)" ./cmd/daemon

build-client:
	go build -v -o $(BIN_CLIENT) -ldflags "$(LDFLAGS)" ./cmd/client

build: build-daemon build-client

version: build  ## Project version
	$(BIN) version

run-daemon: build-daemon ## Run monitor app
	$(BIN_DAEMON) --config ./configs/config.toml

run-client: build-client ## Run client app
	$(BIN_CLIENT) --config ./configs/config_client.toml

test: ## Execute tests
	go test  -covermode=atomic -coverprofile=coverage.out -race -count 100 ./internal/... ./api/... ./pkg/...

integration: ## Execute integration tests
	ginkgo -p -v --repeat=10 ./tests --

coverage: test ## Test coverage
	go tool cover -html coverage.out

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.52.2

lint: install-lint-deps ## Run liniter
	golangci-lint run --config=$$(pwd)/.golangci.yml \
		--timeout 3m0s \
		--skip-dirs='/opt/hostedtoolcache/go|/go/pkg/mod' \
		./...

.PHONY: build run version test lint help coverage
